package handler

import (
	"MediaWarp/internal/config"
	"MediaWarp/internal/logging"
	"MediaWarp/internal/service"
	"MediaWarp/internal/service/alist"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/allegro/bigcache/v3"
)

type StrmHandlerFunc func(content string, ua string) string

func getHTTPStrmHandler() (StrmHandlerFunc, error) {
	var cache *bigcache.BigCache
	if config.Cache.Enable && config.Cache.HTTPStrmTTL > 0 && config.HTTPStrm.FinalURL {
		var err error
		cache, err = bigcache.New(context.Background(), bigcache.DefaultConfig(config.Cache.HTTPStrmTTL))
		if err != nil {
			return nil, fmt.Errorf("创建 HTTPStrm 缓存失败: %w", err)
		}
		logging.Info("启用 HTTPStrm 缓存，TTL: ", config.Cache.HTTPStrmTTL)
	}

	client := &http.Client{ // 创建自定义HTTP客户端配置
		Timeout: RedirectTimeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// 禁止自动重定向，以便手动处理
			return http.ErrUseLastResponse
		},
	}
	return func(content string, ua string) string {
		if config.HTTPStrm.FinalURL {
			if cache != nil {
				if cachedURL, err := cache.Get(content); err == nil {
					logging.Infof("HTTPStrm 重定向至: %s (缓存)", string(cachedURL))
					return string(cachedURL)
				}
			}

			logging.Debug("HTTPStrm 启用获取最终 URL，开始尝试获取最终 URL")
			finalURL, err := getFinalURL(client, content, ua)
			if err != nil {
				logging.Warning("获取最终 URL 失败，使用原始 URL: ", err)
			} else {
				logging.Info("HTTPStrm 重定向至: ", finalURL)
			}
			if cache != nil {
				if err := cache.Set(content, []byte(finalURL)); err != nil {
					logging.Warning("缓存 HTTPStrm URL 失败: ", err)
				} else {
					logging.Debug("缓存 HTTPStrm URL 成功")
				}
			}
			return finalURL
		} else {
			logging.Debug("HTTPStrm 未启用获取最终 URL，直接使用原始 URL: ", content)
			return content
		}
	}, nil
}

type resolutionInfo struct {
	width  uint
	height uint
	name   string
}
type TranscodeResourceInfo struct {
	url        string
	isM3U8     bool
	expireAt   time.Time
	resolution resolutionInfo
}

type alistStrmResult struct {
	url                string                  // 重定向 URL
	fileSize           int64                   // 文件大小（字节）
	transcodeResources []TranscodeResourceInfo // 转码资源列表
}

func alistStrmHandler(content string, alistAddr string, needTranscodeResourceInfo bool) (*alistStrmResult, error) {
	startTime := time.Now()
	defer func() {
		logging.Debugf("获取 AlistStrm 重定向 URL 耗时：%s", time.Since(startTime))
	}()

	client, err := service.GetAlistClient(alistAddr)
	if err != nil {
		return nil, fmt.Errorf("获取 AlistClient 失败：%w", err)
	}

	fileData, err := client.FsGet(&alist.FsGetRequest{Path: content, Page: 1})
	if err != nil {
		return nil, fmt.Errorf("获取文件信息失败：%w", err)
	}

	res := alistStrmResult{
		transcodeResources: make([]TranscodeResourceInfo, 0),
	}

	if config.AlistStrm.RawURL {
		res.url = fileData.RawURL
	} else {
		var u strings.Builder
		u.WriteString(client.GetEndpoint())
		if fileData.Sign != "" {
			u.WriteString("?sign=" + fileData.Sign)
		}
		u.WriteString(path.Join("/d", client.GetUserInfo().BasePath, content))
		res.url = u.String()
	}
	logging.Infof("AlistStrm 重定向至：%s", res.url)

	res.fileSize = fileData.Size

	if needTranscodeResourceInfo {
		previewData, err := client.GetVideoPreviewData(content, "")
		if err != nil {
			logging.Warningf("%#v 获取视频预览信息失败：%+v", fileData, err)
			return &res, nil // 即使获取预览信息失败，也返回基本的重定向 URL 和文件大小
		}
		for _, task := range previewData.VideoPreviewPlayInfo.LiveTranscodingTaskList {
			if task.Url != "" {
				u, err := url.Parse(task.Url)
				if err != nil {
					logging.Warningf("解析转码资源 URL 失败: %s, URL: %s", err, task.Url)
					continue
				}
				expireStr := u.Query().Get("x-oss-expires")
				if expireStr == "" {
					logging.Warningf("转码资源 URL 中未找到 x-oss-expires 参数，URL: %s", task.Url)
					continue
				}
				tsInt, err := strconv.ParseInt(expireStr, 10, 64)
				if err != nil {
					logging.Warningf("解析转码资源 URL 中的 x-oss-expires 参数失败: %+v, URL: %s", err, task.Url)
					continue
				}
				info := TranscodeResourceInfo{
					url:      task.Url,
					isM3U8:   strings.HasSuffix(u.Path, ".m3u8"),
					expireAt: time.Unix(tsInt, 0),
					resolution: resolutionInfo{
						width:  uint(task.TemplateHeight),
						height: uint(task.TemplateHeight),
						name:   task.TemplateName,
					},
				}
				res.transcodeResources = append(res.transcodeResources, info)
			}
		}
	}

	return &res, nil
}
