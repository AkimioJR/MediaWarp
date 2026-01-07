package handler

import (
	"MediaWarp/constants"
	"MediaWarp/internal/config"
	"MediaWarp/internal/logging"
	"MediaWarp/internal/service"
	"MediaWarp/internal/service/alist"
	"MediaWarp/internal/service/jellyfin"
	"MediaWarp/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/gin-gonic/gin"
)

// Jellyfin 服务器处理器
type JellyfinHandler struct {
	server            *jellyfin.Jellyfin     // Jellyfin 服务器
	routerRules       []RegexpRouteRule      // 正则路由规则
	proxy             *httputil.ReverseProxy // 反向代理
	httpStrmHandler   StrmHandlerFunc
	playbackInfoMutex sync.Map // 视频流处理并发控制，确保同一个 item ID 的重定向请求串行化，避免重复获取缓存
}

func NewJellyfinHander(addr string, apiKey string) (*JellyfinHandler, error) {
	jellyfinHandler := JellyfinHandler{}
	jellyfinHandler.server = jellyfin.New(addr, apiKey)
	target, err := url.Parse(jellyfinHandler.server.GetEndpoint())
	if err != nil {
		return nil, err
	}
	jellyfinHandler.proxy = httputil.NewSingleHostReverseProxy(target)

	{ // 初始化路由规则
		jellyfinHandler.routerRules = []RegexpRouteRule{
			{
				Regexp: constants.JellyfinRegexp.Router.ModifyPlaybackInfo,
				Handler: responseModifyCreater(
					&httputil.ReverseProxy{Director: jellyfinHandler.proxy.Director},
					jellyfinHandler.ModifyPlaybackInfo,
				),
			},
			{
				Regexp:  constants.JellyfinRegexp.Router.VideosHandler,
				Handler: jellyfinHandler.VideosHandler,
			},
		}
		if config.Web.Enable {
			if config.Web.Index || config.Web.Head != "" || config.Web.ExternalPlayerUrl || config.Web.VideoTogether {
				jellyfinHandler.routerRules = append(
					jellyfinHandler.routerRules,
					RegexpRouteRule{
						Regexp: constants.JellyfinRegexp.Router.ModifyIndex,
						Handler: responseModifyCreater(
							&httputil.ReverseProxy{Director: jellyfinHandler.proxy.Director},
							jellyfinHandler.ModifyIndex,
						),
					},
				)
			}
		}
	}

	jellyfinHandler.httpStrmHandler, err = getHTTPStrmHandler()
	if err != nil {
		return nil, fmt.Errorf("创建 HTTPStrm 处理器失败: %w", err)
	}
	return &jellyfinHandler, nil
}

// 转发请求至上游服务器
func (jellyfinHandler *JellyfinHandler) ReverseProxy(rw http.ResponseWriter, req *http.Request) {
	jellyfinHandler.proxy.ServeHTTP(rw, req)
}

// 正则路由表
func (jellyfinHandler *JellyfinHandler) GetRegexpRouteRules() []RegexpRouteRule {
	return jellyfinHandler.routerRules
}

func (jellyfinHandler *JellyfinHandler) GetImageCacheRegexp() *regexp.Regexp {
	return constants.JellyfinRegexp.Cache.Image
}

func (*JellyfinHandler) GetSubtitleCacheRegexp() *regexp.Regexp {
	return constants.JellyfinRegexp.Cache.Subtitle
}

// 修改播放信息请求
//
// /Items/:itemId
// 强制将 HTTPStrm 设置为支持直链播放和转码、AlistStrm 设置为支持直链播放并且禁止转码
func (jellyfinHandler *JellyfinHandler) ModifyPlaybackInfo(rw *http.Response) error {
	// 检查 IsPlayback 参数，如果为 false 则不做修改直接返回
	// 从响应的请求中获取参数，因为响应对象包含原始请求
	// 使用不区分大小写的方式获取查询参数
	isPlayback := getQueryValueCaseInsensitive(rw.Request.URL.Query(), "IsPlayback")
	logging.Debugf("IsPlayback 参数值: '%s' (请求 URL: %s)", isPlayback, rw.Request.URL.String())
	if strings.ToLower(isPlayback) == "false" {
		logging.Debug("IsPlayback=false，跳过 PlaybackInfo 修改")
		return nil
	}

	defer rw.Body.Close()
	data, err := io.ReadAll(rw.Body)
	if err != nil {
		logging.Warning("读取响应体失败：", err)
		return err
	}

	var playbackInfoResponse jellyfin.PlaybackInfoResponse
	if err = json.Unmarshal(data, &playbackInfoResponse); err != nil {
		logging.Warning("解析 jellyfin.PlaybackInfoResponse JSON 错误：", err)
		return err
	}

	for index, mediasource := range playbackInfoResponse.MediaSources {
		logging.Debug("请求 ItemsServiceQueryItem：" + *mediasource.ID)
		itemResponse, err := jellyfinHandler.server.ItemsServiceQueryItem(*mediasource.ID, 1, "Path,MediaSources") // 查询 item 需要去除前缀仅保留数字部分
		if err != nil {
			logging.Warning("请求 ItemsServiceQueryItem 失败：", err)
			continue
		}
		item := itemResponse.Items[0]
		strmFileType, opt := recgonizeStrmFileType(*item.Path)
		switch strmFileType {
		case constants.HTTPStrm: // HTTPStrm 设置支持直链播放并且支持转码
			if !config.HTTPStrm.TransCode {
				*playbackInfoResponse.MediaSources[index].SupportsDirectPlay = true
				*playbackInfoResponse.MediaSources[index].SupportsDirectStream = true
				playbackInfoResponse.MediaSources[index].TranscodingURL = nil
				playbackInfoResponse.MediaSources[index].TranscodingSubProtocol = nil
				playbackInfoResponse.MediaSources[index].TranscodingContainer = nil
				if mediasource.DirectStreamURL != nil {
					apikeypair, err := utils.ResolveEmbyAPIKVPairs(mediasource.DirectStreamURL)
					if err != nil {
						logging.Warning("解析API键值对失败：", err)
						continue
					}
					directStreamURL := fmt.Sprintf("/Videos/%s/stream?MediaSourceId=%s&Static=true&%s", *mediasource.ID, *mediasource.ID, apikeypair)
					playbackInfoResponse.MediaSources[index].DirectStreamURL = &directStreamURL
					logging.Info(*mediasource.Name, " 强制禁止转码，直链播放链接为: ", directStreamURL)
				}
			}

		case constants.AlistStrm: // AlistStm 设置支持直链播放并且禁止转码
			if !config.AlistStrm.TransCode {
				*playbackInfoResponse.MediaSources[index].SupportsDirectPlay = true
				*playbackInfoResponse.MediaSources[index].SupportsDirectStream = true
				*playbackInfoResponse.MediaSources[index].SupportsTranscoding = false
				playbackInfoResponse.MediaSources[index].TranscodingURL = nil
				playbackInfoResponse.MediaSources[index].TranscodingSubProtocol = nil
				playbackInfoResponse.MediaSources[index].TranscodingContainer = nil
				directStreamURL := fmt.Sprintf("/Videos/%s/stream?MediaSourceId=%s&Static=true", *mediasource.ID, *mediasource.ID)
				if mediasource.DirectStreamURL != nil {
					logging.Debugf("%s 原直链播放链接： %s", *mediasource.Name, *mediasource.DirectStreamURL)
					apikeypair, err := utils.ResolveEmbyAPIKVPairs(mediasource.DirectStreamURL)
					if err != nil {
						logging.Warning("解析API键值对失败：", err)
						continue
					}
					directStreamURL += "&" + apikeypair
				}
				playbackInfoResponse.MediaSources[index].DirectStreamURL = &directStreamURL
				container := strings.TrimPrefix(path.Ext(*mediasource.Path), ".")
				playbackInfoResponse.MediaSources[index].Container = &container
				logging.Infof("%s 强制禁止转码，直链播放链接为：%s，容器为： %s", *mediasource.Name, directStreamURL, container)
			} else {
				logging.Infof("%s 保持原有转码设置", *mediasource.Name)
			}

			if playbackInfoResponse.MediaSources[index].Size == nil {
				alistClient, err := service.GetAlistClient(opt.(string))
				if err != nil {
					logging.Warning("获取 AlistClient 失败：", err)
					continue
				}
				fsGetData, err := alistClient.FsGet(&alist.FsGetRequest{Path: *mediasource.Path, Page: 1})
				if err != nil {
					logging.Warning("请求 FsGet 失败：", err)
					continue
				}
				playbackInfoResponse.MediaSources[index].Size = &fsGetData.Size
				logging.Infof("%s 设置文件大小为：%d", *mediasource.Name, fsGetData.Size)
			}
		}
	}

	if data, err = json.Marshal(playbackInfoResponse); err != nil {
		logging.Warning("序列化 jellyfin.PlaybackInfoResponse Json 错误：", err)
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
func (jellyfinHandler *JellyfinHandler) VideosHandler(ctx *gin.Context) {
	if ctx.Request.Method == http.MethodHead { // 不额外处理 HEAD 请求
		jellyfinHandler.ReverseProxy(ctx.Writer, ctx.Request)
		logging.Debug("VideosHandler 不处理 HEAD 请求，转发至上游服务器")
		return
	}

	// 从 URL 中提取 item ID（例如：/Videos/813a630bcf9c3f693a2ec8c498f868d2/stream 中的 813a630bcf9c3f693a2ec8c498f868d2）
	var itemID string
	path := ctx.Request.URL.Path
	if matches := constants.JellyfinRegexp.Router.VideosHandler.FindStringSubmatch(path); len(matches) > 0 {
		parts := strings.Split(path, "/")
		for i, part := range parts {
			if part == "Videos" && i+1 < len(parts) {
				itemID = parts[i+1]
				break
			}
		}
	}

	// 并发控制：确保同一个 item ID 只有一个任务在运行
	// 将整个处理流程放在锁内，避免重复查询和重复获取重定向 URL
	var muWrapper *mutexWithRefCount
	if itemID != "" {
		// 加载或创建 mutex wrapper
		value, _ := jellyfinHandler.playbackInfoMutex.LoadOrStore(itemID, &mutexWithRefCount{})
		muWrapper = value.(*mutexWithRefCount)

		// 增加引用计数
		atomic.AddInt32(&muWrapper.refCount, 1)

		// 锁定并处理
		muWrapper.mu.Lock()
		defer func() {
			muWrapper.mu.Unlock()
			// 减少引用计数
			refCount := atomic.AddInt32(&muWrapper.refCount, -1)
			// 如果没有其他 goroutine 在使用，删除这个 mutex 以避免内存泄漏
			if refCount == 0 {
				jellyfinHandler.playbackInfoMutex.Delete(itemID)
			}
		}()
		logging.Debugf("开始处理 item %s 的 VideosHandler 请求", itemID)
	}

	mediaSourceID := ctx.Query("mediasourceid")
	logging.Debugf("请求 ItemsServiceQueryItem：%s", mediaSourceID)
	itemResponse, err := jellyfinHandler.server.ItemsServiceQueryItem(mediaSourceID, 1, "Path,MediaSources") // 查询 item 需要去除前缀仅保留数字部分
	if err != nil {
		logging.Warning("请求 ItemsServiceQueryItem 失败：", err)
		jellyfinHandler.proxy.ServeHTTP(ctx.Writer, ctx.Request)
		return
	}

	item := itemResponse.Items[0]

	if !strings.HasSuffix(strings.ToLower(*item.Path), ".strm") { // 不是 Strm 文件
		logging.Debugf("播放本地视频：%s，不进行处理", *item.Path)
		jellyfinHandler.proxy.ServeHTTP(ctx.Writer, ctx.Request)
		return
	}

	strmFileType, opt := recgonizeStrmFileType(*item.Path)
	for _, mediasource := range item.MediaSources {
		if *mediasource.ID == mediaSourceID { // EmbyServer >= 4.9 返回的ID带有前缀mediasource_
			switch strmFileType {
			case constants.HTTPStrm:
				if *mediasource.Protocol == jellyfin.HTTP {
					// httpStrmHandler 内部有缓存机制，锁确保串行化访问
					ctx.Redirect(http.StatusFound, jellyfinHandler.httpStrmHandler(*mediasource.Path, ctx.Request.UserAgent()))
					return
				}

			case constants.AlistStrm: // 无需判断 *mediasource.Container 是否以Strm结尾，当 AlistStrm 存储的位置有对应的文件时，*mediasource.Container 会被设置为文件后缀
				redirectURL := alistStrmHandler(*mediasource.Path, opt.(string))
				if redirectURL != "" {
					ctx.Redirect(http.StatusFound, redirectURL)
				}
				return

			case constants.UnknownStrm:
				jellyfinHandler.proxy.ServeHTTP(ctx.Writer, ctx.Request)
				return
			}
		}
	}
}

// 修改首页函数
func (jellyfinHandler *JellyfinHandler) ModifyIndex(rw *http.Response) error {
	var (
		htmlFilePath string = path.Join(config.CostomDir(), "index.html")
		htmlContent  []byte
		addHEAD      bytes.Buffer
		err          error
	)

	defer rw.Body.Close() // 无论哪种情况，最终都要确保原 Body 被关闭，避免内存泄漏
	if config.Web.Index { // 从本地文件读取index.html
		if htmlContent, err = os.ReadFile(htmlFilePath); err != nil {
			logging.Warning("读取文件内容出错，错误信息：", err)
			return err
		}
	} else { // 从上游获取响应体
		if htmlContent, err = io.ReadAll(rw.Body); err != nil {
			return err
		}
	}

	if config.Web.Head != "" { // 用户自定义HEAD
		addHEAD.WriteString(config.Web.Head + "\n")
	}
	if config.Web.ExternalPlayerUrl { // 外部播放器
		addHEAD.WriteString(`<script src="/MediaWarp/static/embyExternalUrl/embyWebAddExternalUrl/embyLaunchPotplayer.js"></script>` + "\n")
	}
	if config.Web.Crx { // crx 美化
		addHEAD.WriteString(`<link rel="stylesheet" id="theme-css" href="/MediaWarp/static/jellyfin-crx/static/css/style.css" type="text/css" media="all" />
    <script src="/MediaWarp/static/jellyfin-crx/static/js/common-utils.js"></script>
    <script src="/MediaWarp/static/jellyfin-crx/static/js/jquery-3.6.0.min.js"></script>
    <script src="/MediaWarp/static/jellyfin-crx/static/js/md5.min.js"></script>
    <script src="/MediaWarp/static/jellyfin-crx/content/main.js"></script>` + "\n")
	}
	if config.Web.ActorPlus { // 过滤没有头像的演员和制作人员
		addHEAD.WriteString(`<script src="/MediaWarp/static/emby-web-mod/actorPlus/actorPlus.js"></script>` + "\n")
	}
	if config.Web.FanartShow { // 显示同人图（fanart图）
		addHEAD.WriteString(`<script src="/MediaWarp/static/emby-web-mod/fanart_show/fanart_show.js"></script>` + "\n")
	}
	if config.Web.Danmaku { // 弹幕
		addHEAD.WriteString(`<script src="/MediaWarp/static/jellyfin-danmaku/ede.js" defer></script>` + "\n")
	}
	if config.Web.VideoTogether { // VideoTogether
		addHEAD.WriteString(`<script src="https://2gether.video/release/extension.website.user.js"></script>` + "\n")
	}

	addHEAD.WriteString(`<!-- MediaWarp Web 页面修改功能 -->` + "\n" + "</head>")

	htmlContent = bytes.Replace(htmlContent, []byte("</head>"), addHEAD.Bytes(), 1) // 将添加HEAD

	rw.Header.Set("Content-Length", strconv.Itoa(len(htmlContent)))
	rw.Body = io.NopCloser(bytes.NewReader(htmlContent))
	return nil
}

var _ MediaServerHandler = (*JellyfinHandler)(nil) // 确保 JellyfinHandler 实现 MediaServerHandler 接口
