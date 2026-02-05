package handler

import (
	"MediaWarp/constants"
	"MediaWarp/internal/logging"
	"MediaWarp/utils"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strconv"
	"time"

	"github.com/tidwall/gjson"
)

type FNTVHandler struct {
	routerRules     []RegexpRouteRule      // 正则路由规则
	proxy           *httputil.ReverseProxy // 反向代理
	httpStrmHandler StrmHandlerFunc
}

func NewFNTVHandler(addr string) (*FNTVHandler, error) {
	hanler := FNTVHandler{}
	target, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}
	hanler.proxy = httputil.NewSingleHostReverseProxy(target)

	hanler.routerRules = []RegexpRouteRule{
		{
			Regexp: constants.FNTVRegexp.StreamHandler,
			Handler: responseModifyCreater(
				&httputil.ReverseProxy{Director: hanler.proxy.Director},
				hanler.ModifyStream,
			),
		},
	}

	hanler.httpStrmHandler, err = getHTTPStrmHandler()
	if err != nil {
		return nil, fmt.Errorf("创建 HTTPStrm 处理器失败: %w", err)
	}

	return &hanler, nil
}

// 转发请求至上游服务器
func (hanler *FNTVHandler) ReverseProxy(writer http.ResponseWriter, request *http.Request) {
	hanler.proxy.ServeHTTP(writer, request)
}

// 获取正则路由表
func (hanler *FNTVHandler) GetRegexpRouteRules() []RegexpRouteRule {
	return hanler.routerRules
}

// 获取图片缓存正则表达式
func (hanler *FNTVHandler) GetImageCacheRegexp() *regexp.Regexp {
	return constants.FNTVRegexp.Cache.Image
}

// 获取字幕缓存正则表达式
func (hanler *FNTVHandler) GetSubtitleCacheRegexp() *regexp.Regexp {
	return constants.FNTVRegexp.Cache.Subtitle
}

func (hanler *FNTVHandler) ModifyStream(rw *http.Response) error {
	startTime := time.Now()
	defer func() {
		logging.Debugf("FNTV ModifyStream 处理耗时: %s", time.Since(startTime).String())
	}()

	data, err := io.ReadAll(rw.Body)
	if err != nil {
		logging.Warning("读取响应体失败：", err)
		return err
	}
	defer rw.Body.Close()
	logging.Debug(string(data))

	jsonChain := utils.NewFromBytesWithCopy(data, jsonChainOption)

	code := jsonChain.Get("code").Int()
	if code != 0 {
		logging.Warningf("stream 响应 code: %d, msg: %s", code, jsonChain.Get("msg").String())
		return nil
	}

	filePathRes := jsonChain.Get("data.file_stream.path")
	if filePathRes.Type != gjson.String {
		logging.Warningf("stream 响应 data.file_stream.path 字段不正确: %#v", filePathRes)
		return nil
	}

	filePath := filePathRes.String()

	strmFileType, opt := recgonizeStrmFileType(filePath)

	switch strmFileType {
	case constants.HTTPStrm: // HTTPStrm 设置支持直链播放并且支持转码
		urlRes := jsonChain.Get("data.direct_link_qualities.0.url")
		if urlRes.Type != gjson.String {
			logging.Warningf("stream 响应 data.direct_link_qualities.0.url 字段不正确: %#v", urlRes)
			return nil
		}

		redirectURL := hanler.httpStrmHandler(urlRes.String(), rw.Request.Header.Get("User-Agent"))
		jsonChain.Set(
			"data.direct_link_qualities.0.resolution",
			"HTTPStrm 直链",
		).Set(
			"data.direct_link_qualities.0.url",
			redirectURL,
		)

	case constants.AlistStrm: // AlistStm 设置支持直链播放并且禁止转码
		remoteFilepathRes := jsonChain.Get("data.direct_link_qualities.0.url")
		if remoteFilepathRes.Type != gjson.String {
			logging.Warningf("stream 响应 data.direct_link_qualities.0.url 字段不正确: %#v", remoteFilepathRes)
			return nil
		}

		redirectURL := alistStrmHandler(remoteFilepathRes.String(), opt.(string))
		jsonChain.Set(
			"data.direct_link_qualities.0.resolution",
			"AlistStrm 直链 - 原画",
		).Set(
			"data.direct_link_qualities.0.url",
			redirectURL,
		)

	default:
		logging.Debugf("%s 未匹配任何 Strm 类型，保持原有播放链接不变", filePath)
	}

	data, err = jsonChain.Result()
	if err != nil {
		logging.Warning("操作 FNTV Stream Json 错误：", err)
		return err
	}
	rw.Header.Set("Content-Type", "application/json") // 更新 Content-Type 头
	rw.Header.Set("Content-Length", strconv.Itoa(len(data)))
	rw.Body = io.NopCloser(bytes.NewReader(data))

	return nil
}

var _ MediaServerHandler = (*FNTVHandler)(nil)
