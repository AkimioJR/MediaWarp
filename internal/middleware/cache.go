package middleware

import (
	"MediaWarp/internal/logging"
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"sort"

	"github.com/allegro/bigcache/v3"
	"github.com/gin-gonic/gin"
)

type CacheData struct {
	StatusCode int         // code 响应码
	Header     http.Header // header 响应头信息
	Body       []byte      // body 响应体
}

func (c *CacheData) Json() ([]byte, error) {
	return json.Marshal(c)
}

func (c *CacheData) WriteResponse(ctx *gin.Context) {
	ctx.Status(c.StatusCode)            // 设置响应码
	for key, values := range c.Header { // 设置响应头
		for _, value := range values {
			ctx.Writer.Header().Add(key, value)
		}
	}
	ctx.Writer.Write(c.Body) // 设置响应体
}

func ParseCacheData(data []byte) (*CacheData, error) {
	var cacheData CacheData
	if err := json.Unmarshal(data, &cacheData); err != nil {
		return nil, err
	}
	return &cacheData, nil
}

// 自定义的请求响应器
//
// 用于记录缓存数据
type WriterWarp struct {
	gin.ResponseWriter
	Body bytes.Buffer
}

var _ gin.ResponseWriter = (*WriterWarp)(nil)

func (w *WriterWarp) Write(data []byte) (int, error) {
	w.Body.Write(data)
	return w.ResponseWriter.Write(data)
}

var _ gin.ResponseWriter = (*WriterWarp)(nil)

// 计算Key时忽略的查询参数
var CacheKeyIgnoreQuery = []string{
	"api_key",

	// Fileball
	"starttimeticks",
	"x-playback-session-id",

	// Emby
	"playsessionid",
	"tag", // 忽略 ImageTag 参数，防止 Emby 重启后缓存失效
}

// 计算Key时忽略的请求头
// var cacheKeyIgnoreHeaders = []string{
// 	"Range",
// 	"Host",
// 	"Referrer",
// 	"Connection",
// 	"Accept",
// 	"Accept-Encoding",
// 	"Accept-Language",
// 	"Cache-Control",
// 	"Upgrade-Insecure-Requests",
// 	"Referer",
// 	"Origin",

// 	// StreamMusic
// 	"X-Streammusic-Audioid",
// 	"X-Streammusic-Savepath",

// 	// IP
// 	"X-Forwarded-For",
// 	"X-Real-IP",
// 	"Forwarded",
// 	"Client-IP",
// 	"True-Client-IP",
// 	"CF-Connecting-IP",
// 	"X-Cluster-Client-IP",
// 	"Fastly-Client-IP",
// 	"X-Client-IP",
// 	"X-ProxyUser-IP",
// 	"Via",
// 	"Forwarded-For",
// 	"X-From-Cdn",
// }

func getCacheKey(ctx *gin.Context) string {
	var (
		path  string     = ctx.Request.URL.Path    // 请求路径
		query url.Values = ctx.Request.URL.Query() // 查询参数
		// header    http.Header = ctx.Request.Header      // 请求头
		// headerStr string      = ""                      // 请求头字符串
	)

	// 将查询参数转化为字符串
	for _, key := range CacheKeyIgnoreQuery {
		query.Del(key)
	}
	for key, values := range query { // 对查询参数的值进行排序
		sort.Strings(values)
		query[key] = values
	}

	// 将请求头转化为字符串
	// for _, key := range cacheKeyIgnoreHeaders {
	// 	header.Del(key)
	// }
	// headKeys := make([]string, 0, len(header))
	// for key := range header {
	// 	headKeys = append(headKeys, key)
	// }
	// sort.Strings(headKeys) // 对请求头的键进行排序
	// for _, key := range headKeys {
	// 	sort.Strings(header[key]) // 对请求头的值进行排序
	// 	headerStr += fmt.Sprintf("%s=%s;", key, strings.Join(header[key], "|"))
	// }

	return path + query.Encode() // + headerStr
}

func getCacheBaseFunc(cachePool *bigcache.BigCache, cacheName string, reg string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cacheKey := getCacheKey(ctx)
		logging.AccessDebugf(ctx, "命中 %s 缓存正则表达式: %s, CacheKey: %s", cacheName, reg, cacheKey)
		if cacheByte, err := cachePool.Get(cacheKey); err == nil {
			if cacheData, err := ParseCacheData(cacheByte); err == nil {
				logging.AccessDebugf(ctx, "命中 %s 缓存: %s", cacheName, cacheKey)
				cacheData.WriteResponse(ctx)
				ctx.Abort()
				return
			} else {
				logging.AccessWarningf(ctx, "解析 %s 缓存失败: %v", cacheName, err)
			}
		}

		writer := &WriterWarp{
			ResponseWriter: ctx.Writer,
			Body:           bytes.Buffer{},
		}
		ctx.Writer = writer

		ctx.Next() // 处理请求

		code := ctx.Writer.Status()
		if code >= http.StatusOK && code < http.StatusMultipleChoices { // 响应是2xx的成功响应，更新缓存记录
			cacheData := &CacheData{ // 创建缓存数据
				StatusCode: code, //ctx.Request.Response.StatusCode,
				Header:     ctx.Writer.Header().Clone(),
				Body:       writer.Body.Bytes(),
			}

			if cacheByte, err := cacheData.Json(); err == nil {
				err = cachePool.Set(cacheKey, cacheByte)
				if err != nil {
					logging.AccessWarningf(ctx, "写入 %s 缓存失败: %v", cacheName, err)
				} else {
					logging.AccessDebugf(ctx, "写入 %s 缓存成功", cacheName)
				}
			}
		} else {
			logging.AccessDebugf(ctx, "响应码为: %d, 不进行 %s 缓存", code, cacheName)
		}
	}
}
