package handler

import (
	"MediaWarp/constants"
	"MediaWarp/internal/config"
	"MediaWarp/internal/logging"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"reflect"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/andybalholm/brotli"
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

// 读取响应体
//
// 读取响应体，解压缩 GZIP、Brotli 数据（若响应体被压缩）
func readBody(rw *http.Response) ([]byte, error) {
	encoding := rw.Header.Get("Content-Encoding")

	var reader io.Reader
	switch encoding {
	case "gzip":
		logging.Debug("解码 GZIP 数据")
		gr, err := gzip.NewReader(rw.Body)
		if err != nil {
			return nil, fmt.Errorf("gzip reader error: %w", err)
		}
		defer gr.Close()
		reader = gr

	case "br":
		logging.Debug("解码 Brotli 数据")
		reader = brotli.NewReader(rw.Body)

	case "": // 无压缩
		logging.Debug("无压缩数据")
		reader = rw.Body

	default:
		return nil, fmt.Errorf("unsupported Content-Encoding: %s", encoding)
	}
	return io.ReadAll(reader)
}

// 更新响应体
//
// 修改响应体、更新Content-Length
func updateBody(rw *http.Response, content []byte) error {
	encoding := rw.Header.Get("Content-Encoding")
	var (
		compressed bytes.Buffer
		writer     io.Writer
	)

	// 根据原始编码选择压缩方式
	switch encoding {
	case "gzip":
		logging.Debug("使用 GZIP 重新编码数据")
		writer = gzip.NewWriter(&compressed)

	case "br":
		logging.Debug("使用 Brotli 重新编码数据")
		writer = brotli.NewWriter(&compressed)

	case "": // 无压缩
		logging.Debug("无压缩数据")
		writer = &compressed

	default:
		logging.Warningf("不支持的重新编码：%s，将不对数据进行压缩编码", encoding)
		rw.Header.Del("Content-Encoding")
		writer = &compressed
	}

	if _, err := writer.Write(content); err != nil {
		return fmt.Errorf("compression write error: %w", err)
	}

	// 显式关闭以确保数据完整写入 buffer
	if gw, ok := writer.(*gzip.Writer); ok {
		if err := gw.Close(); err != nil {
			return fmt.Errorf("gzip close error: %w", err)
		}
	} else if bw, ok := writer.(*brotli.Writer); ok {
		if err := bw.Close(); err != nil {
			return fmt.Errorf("brotli close error: %w", err)
		}
	}

	// 设置新 Body
	rw.Body = io.NopCloser(bytes.NewReader(compressed.Bytes()))
	rw.ContentLength = int64(compressed.Len())
	rw.Header.Set("Content-Length", strconv.Itoa(compressed.Len())) // 更新响应头

	return nil
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
