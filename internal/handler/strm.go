package handler

import (
	"MediaWarp/internal/config"
	"MediaWarp/internal/logging"
	"MediaWarp/internal/service"
	"context"
	"fmt"
	"net/http"
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

func alistStrmHandler(content string, alistAddr string) string {
	startTime := time.Now()
	defer func() {
		logging.Debugf("获取 AlistStrm 重定向 URL 耗时：%s", time.Since(startTime))
	}()

	alistClient, err := service.GetAlistClient(alistAddr)
	if err != nil {
		logging.Warning("获取 AlistClient 失败：", err)
		return ""
	}
	url, err := alistClient.GetFileURL(content, config.AlistStrm.RawURL)
	if err != nil {
		logging.Warning("获取文件 URL 失败：", err)
		return ""
	}
	logging.Infof("AlistStrm 重定向至：%s", url)
	return url
}
