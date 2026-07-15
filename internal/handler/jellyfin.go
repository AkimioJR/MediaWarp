package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/AkimioJR/MediaWarp/constants"
	"github.com/AkimioJR/MediaWarp/internal/config"
	"github.com/AkimioJR/MediaWarp/internal/logging"
	"github.com/AkimioJR/MediaWarp/internal/service/jellyfin"
	"github.com/AkimioJR/MediaWarp/utils"

	"github.com/gin-gonic/gin"
)

// Jellyfin 服务器处理器
type JellyfinHandler struct {
	client          *jellyfin.Client       // Jellyfin 客户端
	routerRules     []RegexpRouteRule      // 正则路由规则
	proxy           *httputil.ReverseProxy // 反向代理
	httpStrmHandler StrmHandlerFunc
	// playbackInfoMutex sync.Map // 视频流处理并发控制，确保同一个 item ID 的重定向请求串行化，避免重复获取缓存
}

func NewJellyfinHandler(addr string, apiKey string) (*JellyfinHandler, error) {
	handler := JellyfinHandler{}
	client, err := jellyfin.New(addr, apiKey)
	if err != nil {
		return nil, fmt.Errorf("创建 Jellyfin 客户端实例失败: %w", err)
	}
	handler.client = client
	target, err := url.Parse(handler.client.GetBaseURLString())
	if err != nil {
		return nil, err
	}
	handler.proxy = &httputil.ReverseProxy{
		Rewrite: newSingleHostRewrite(target),
	}

	{ // 初始化路由规则
		handler.routerRules = []RegexpRouteRule{
			{
				Regexp: constants.JellyfinRegexp.Router.ModifyPlaybackInfo,
				Handler: responseModifyCreater(
					&httputil.ReverseProxy{Rewrite: handler.proxy.Rewrite},
					handler.ModifyPlaybackInfo,
				),
			},
			{
				Regexp:  constants.JellyfinRegexp.Router.VideosHandler,
				Handler: handler.VideosHandler,
			},
		}
		if config.Web.Enable {
			if config.Web.Index || config.Web.Head != "" || config.Web.ExternalPlayerUrl || config.Web.VideoTogether {
				handler.routerRules = append(
					handler.routerRules,
					RegexpRouteRule{
						Regexp: constants.JellyfinRegexp.Router.ModifyIndex,
						Handler: responseModifyCreater(
							&httputil.ReverseProxy{Rewrite: handler.proxy.Rewrite},
							generateWebModifier(constants.JELLYFIN),
						),
					},
				)
			}
		}
	}

	handler.httpStrmHandler, err = getHTTPStrmHandler()
	if err != nil {
		return nil, fmt.Errorf("创建 HTTPStrm 处理器失败: %w", err)
	}
	return &handler, nil
}

// 转发请求至上游服务器
func (handler *JellyfinHandler) ReverseProxy(rw http.ResponseWriter, req *http.Request) {
	handler.proxy.ServeHTTP(rw, req)
}

// 正则路由表
func (handler *JellyfinHandler) GetRegexpRouteRules() []RegexpRouteRule {
	return handler.routerRules
}

func (handler *JellyfinHandler) GetImageCacheRegexp() *regexp.Regexp {
	return constants.JellyfinRegexp.Cache.Image
}

func (*JellyfinHandler) GetSubtitleCacheRegexp() *regexp.Regexp {
	return constants.JellyfinRegexp.Cache.Subtitle
}

// 修改播放信息请求
//
// /Items/:itemId
// 强制将 HTTPStrm 设置为支持直链播放和转码、AlistStrm 设置为支持直链播放并且禁止转码
func (handler *JellyfinHandler) ModifyPlaybackInfo(rw *http.Response) error {
	startTime := time.Now()
	defer func() {
		logging.Debugf("处理 ModifyPlaybackInfo 耗时：%s", time.Since(startTime))
	}()

	defer rw.Body.Close()
	data, err := io.ReadAll(rw.Body)
	if err != nil {
		logging.Warning("读取响应体失败：", err)
		return err
	}

	jsonChain := utils.NewJsonChainFromBytesWithCopy(data, jsonChainOption)

	var playbackInfoResponse jellyfin.PlaybackInfoResponse
	if err = json.Unmarshal(data, &playbackInfoResponse); err != nil {
		logging.Warning("解析 jellyfin.PlaybackInfoResponse JSON 错误：", err)
		return err
	}

	for index, mediasource := range playbackInfoResponse.MediaSources {
		if mediasource.ID == nil {
			logging.Debugf("跳过 ID 为空的 MediaSource, index: %d", index)
			continue
		}
		startTime := time.Now()
		logging.Debug("请求 ItemsServiceQueryItem：" + *mediasource.ID)
		itemResponse, err := handler.client.ItemsServiceQueryItem(*mediasource.ID, 1, "Path,MediaSources") // 查询 item 需要去除前缀仅保留数字部分
		if err != nil {
			logging.Warning("请求 ItemsServiceQueryItem 失败：", err)
			continue
		}
		if len(itemResponse.Items) == 0 {
			logging.Warningf("ItemsServiceQueryItem 返回空结果，MediaSource ID: %s", *mediasource.ID)
			continue
		}
		item := itemResponse.Items[0]
		if item.Path == nil {
			logging.Warningf("Item 的 Path 为空，MediaSource ID: %s", *mediasource.ID)
			continue
		}
		strmFileType, opt := recgonizeStrmFileType(*item.Path)
		bsePath := "MediaSources." + strconv.Itoa(index) + "."
		switch strmFileType {
		case constants.HTTPStrm: // HTTPStrm 设置支持直链播放并且支持转码
			processHTTPStrmPlaybackInfo(
				jsonChain,
				bsePath,
				*mediasource.ID,
				*mediasource.ID,
				mediasource.DirectStreamURL,
			)

		case constants.AlistStrm: // AlistStm 设置支持直链播放并且禁止转码
			processAlistStrmPlaybackInfo(
				jsonChain,
				bsePath,
				*mediasource.ID,
				*mediasource.ID,
				opt.(string),
				mediasource.DirectStreamURL,
				*item.Path,
				mediasource.Size,
			)
		}

		logging.Debugf("处理 %s 的 MediaSource %s 耗时：%s", *item.Path, *mediasource.ID, time.Since(startTime))
	}

	data, err = jsonChain.Result()
	if err != nil {
		logging.Warning("操作 jellyfin.PlaybackInfoResponse Json 错误：", err)
		return err
	}

	rw.Header.Set("Content-Type", "application/json") // 更新 Content-Type 头
	rw.Header.Set("Content-Length", strconv.Itoa(len(data)))
	rw.Body = io.NopCloser(bytes.NewReader(data))
	return nil
}

// 视频流处理器
//
// 支持播放本地视频、重定向 HttpStrm、AlistStrm
func (handler *JellyfinHandler) VideosHandler(ctx *gin.Context) {
	if ctx.Request.Method == http.MethodHead { // 不额外处理 HEAD 请求
		handler.ReverseProxy(ctx.Writer, ctx.Request)
		logging.Debug("VideosHandler 不处理 HEAD 请求，转发至上游服务器")
		return
	}

	mediaSourceID := ctx.Query("mediasourceid")
	logging.Debugf("请求 ItemsServiceQueryItem：%s", mediaSourceID)
	itemResponse, err := handler.client.ItemsServiceQueryItem(mediaSourceID, 1, "Path,MediaSources") // 查询 item 需要去除前缀仅保留数字部分
	if err != nil {
		logging.Warning("请求 ItemsServiceQueryItem 失败：", err)
		handler.proxy.ServeHTTP(ctx.Writer, ctx.Request)
		return
	}

	if len(itemResponse.Items) == 0 {
		logging.Warningf("ItemsServiceQueryItem 返回空结果，mediaSourceID: %s", mediaSourceID)
		handler.proxy.ServeHTTP(ctx.Writer, ctx.Request)
		return
	}

	item := itemResponse.Items[0]

	if item.Path == nil || !strings.HasSuffix(strings.ToLower(*item.Path), ".strm") { // 不是 Strm 文件
		logging.Debugf("播放本地视频，不进行处理")
		handler.proxy.ServeHTTP(ctx.Writer, ctx.Request)
		return
	}

	strmFileType, opt := recgonizeStrmFileType(*item.Path)
	for _, mediasource := range item.MediaSources {
		if mediasource.ID == nil {
			continue
		}
		if *mediasource.ID == mediaSourceID { // EmbyServer >= 4.9 返回的ID带有前缀mediasource_
			if mediasource.Path == nil {
				logging.Warningf("MediaSource (ID: %s) 的 Path 为空", *mediasource.ID)
				handler.ReverseProxy(ctx.Writer, ctx.Request)
				return
			}
			switch strmFileType {
			case constants.HTTPStrm:
				if mediasource.Protocol != nil && *mediasource.Protocol == jellyfin.HTTP {
					ctx.Redirect(http.StatusFound, handler.httpStrmHandler(*mediasource.Path, ctx.Request.UserAgent()))
					return
				}

			case constants.AlistStrm: // 无需判断 *mediasource.Container 是否以Strm结尾，当 AlistStrm 存储的位置有对应的文件时，*mediasource.Container 会被设置为文件后缀
				res, err := alistStrmHandler(*mediasource.Path, opt.(string), false)
				if err != nil {
					logging.Warningf("获取 AlistStrm 重定向 URL 失败:%#v", err)
					handler.ReverseProxy(ctx.Writer, ctx.Request)
					return
				}
				ctx.Redirect(http.StatusFound, res.url)
				return

			case constants.UnknownStrm:
				handler.proxy.ServeHTTP(ctx.Writer, ctx.Request)
				return
			}
		}
	}
}

var _ MediaServerHandler = (*JellyfinHandler)(nil) // 确保 JellyfinHandler 实现 MediaServerHandler 接口
