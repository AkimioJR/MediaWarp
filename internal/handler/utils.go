package handler

import (
	"MediaWarp/constants"
	"MediaWarp/internal/config"
	"MediaWarp/internal/logging"
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"reflect"
	"runtime"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// 响应修改创建器
//
// 将需要修改上游响应的处理器包装成一个 gin.HandlerFunc 处理器
func responseModifyCreater(proxy *httputil.ReverseProxy, modifyResponseFN func(rw *http.Response) error) gin.HandlerFunc {
	funcPtr := reflect.ValueOf(modifyResponseFN).Pointer()
	funcName := strings.ReplaceAll(runtime.FuncForPC(funcPtr).Name(), "-fm", "")
	logging.Debugf("创建响应修改处理器：%s", funcName)

	proxy.ModifyResponse = func(rw *http.Response) error {
		defer func() {
			if r := recover(); r != nil {
				logging.Errorf("%s 发生 panic：%s\n%s", funcName, r, string(debug.Stack()))
			}
		}()
		return modifyResponseFN(rw)
	}

	return func(ctx *gin.Context) {
		proxy.ServeHTTP(ctx.Writer, ctx.Request)
	}
}

// 不区分大小写地获取查询参数值
//
// 从 url.Values 中查找指定键名的值，忽略大小写
func getQueryValueCaseInsensitive(query url.Values, key string) string {
	keyLower := strings.ToLower(key)
	for k, v := range query {
		if strings.ToLower(k) == keyLower {
			if len(v) > 0 {
				return v[0]
			}
			return ""
		}
	}
	return ""
}

// 根据 Strm 文件路径识别 Strm 文件类型
//
// 返回 Strm 文件类型和一个可选配置
func recgonizeStrmFileType(strmFilePath string) (constants.StrmFileType, any) {
	if config.HTTPStrm.Enable {
		for _, prefix := range config.HTTPStrm.PrefixList {
			if strings.HasPrefix(strmFilePath, prefix) {
				logging.Debugf("%s 成功匹配路径：%s，Strm 类型：%s", strmFilePath, prefix, constants.HTTPStrm)
				return constants.HTTPStrm, nil
			}
		}
	}
	if config.AlistStrm.Enable {
		for _, alistStrmConfig := range config.AlistStrm.List {
			for _, prefix := range alistStrmConfig.PrefixList {
				if strings.HasPrefix(strmFilePath, prefix) {
					logging.Debugf("%s 成功匹配路径：%s，Strm 类型：%s，AlistServer 地址：%s", strmFilePath, prefix, constants.AlistStrm, alistStrmConfig.ADDR)
					return constants.AlistStrm, alistStrmConfig.ADDR
				}
			}
		}
	}
	logging.Debugf("%s 未匹配任何路径，Strm 类型：%s", strmFilePath, constants.UnknownStrm)
	return constants.UnknownStrm, nil
}

const (
	MaxRedirectAttempts = 10               // 最大重定向次数限制
	RedirectTimeout     = 10 * time.Second // 最大超时时间

)

var (
	ErrInvalidLocationHeader = errors.New("重定向 Location 头无效")
	ErrMaxRedirectsExceeded  = fmt.Errorf("超过最大重定向次数限制（%d）", MaxRedirectAttempts)
)

// 获取URL的最终目标地址（自动跟踪重定向）
func getFinalURL(client *http.Client, rawURL string, ua string) (string, error) {
	startTime := time.Now()
	defer func() {
		logging.Debugf("获取 %s 最终URL耗时：%s", rawURL, time.Since(startTime))
	}()

	parsedURL, err := url.Parse(rawURL) // 验证并解析输入URL
	if err != nil {
		return "", fmt.Errorf("非法 URL： %w", err)
	}
	if parsedURL.Scheme == "" {
		return "", fmt.Errorf("URL 缺少协议头： %s", parsedURL)
	}

	currentURL := parsedURL.String()
	visited := make(map[string]struct{}, MaxRedirectAttempts)
	redirectChain := make([]string, 0, MaxRedirectAttempts+1)

	// 跟踪重定向链
	for i := 0; i <= MaxRedirectAttempts; i++ {
		// 检测循环重定向
		if _, exists := visited[currentURL]; exists {
			return "", fmt.Errorf("检测到循环重定向，重定向链: %s", strings.Join(redirectChain, " -> "))
		}
		visited[currentURL] = struct{}{}
		redirectChain = append(redirectChain, currentURL)

		req, err := http.NewRequest(http.MethodHead, currentURL, nil) // 创建 HEAD 请求（更高效，只获取头部信息）
		if err != nil {
			return "", fmt.Errorf("创建请求失败: %w", err)
		}
		req.Header.Set("User-Agent", ua) // 设置 User-Agent 头部

		resp, err := client.Do(req)
		if err != nil {
			return "", fmt.Errorf("发送 HTTP 请求失败：%w", err)
		}
		defer resp.Body.Close()

		// 检查是否需要重定向 (3xx 状态码)
		if resp.StatusCode >= http.StatusMultipleChoices && resp.StatusCode < http.StatusBadRequest {
			location, err := resp.Location()
			if err != nil {
				return "", ErrInvalidLocationHeader
			}
			currentURL = location.String()
			continue
		}

		// 返回最终的非重定向URL
		logging.Debug("重定向链：", strings.Join(redirectChain, " -> "))
		return resp.Request.URL.String(), nil
	}

	return "", ErrMaxRedirectsExceeded
}
