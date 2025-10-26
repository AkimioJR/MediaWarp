package middleware

import (
	"MediaWarp/internal/logging"
	"bytes"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/allegro/bigcache"
	"github.com/gin-gonic/gin"
)

func ImageCache(ttl time.Duration, reg *regexp.Regexp) gin.HandlerFunc {
	cachePool, err := bigcache.NewBigCache(bigcache.DefaultConfig(ttl))
	if err != nil {
		panic(fmt.Sprintf("create image cache pool failed: %v", err))
	}
	logging.Debugf("图片缓存中间件已启用, TTL: %s", ttl.String())

	return func(ctx *gin.Context) {
		if ctx.Request.Method != http.MethodGet || !reg.MatchString(ctx.Request.URL.Path) {
			ctx.Next()
			return
		}

		logging.AccessDebug(ctx, "命中图片缓存正则表达式: "+reg.String())

		cacheKey := getCacheKey(ctx)
		logging.AccessDebug(ctx, "Cache Key: "+cacheKey)
		if cacheByte, err := cachePool.Get(cacheKey); err == nil {
			if cacheData, err := ParseCacheData(cacheByte); err == nil {
				logging.AccessDebug(ctx, "命中缓存: "+cacheKey)
				cacheData.WriteResponse(ctx)
				ctx.Abort()
				return
			} else {
				logging.AccessWarningf(ctx, "解析缓存失败: %v", err)
			}
		}

		writer := &WriterWarp{
			ResponseWriter: ctx.Writer,
			Body:           &bytes.Buffer{},
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
					logging.AccessWarningf(ctx, "写入缓存失败: %v", err)
				} else {
					logging.AccessDebug(ctx, "写入缓存成功")
				}
			}
		} else {
			logging.AccessDebugf(ctx, "响应码为: %d, 不进行缓存", code)
		}
	}
}
