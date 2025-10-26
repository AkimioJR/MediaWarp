package middleware

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/gin-gonic/gin"
)

func SubtitleCache(ttl time.Duration, reg *regexp.Regexp) gin.HandlerFunc {
	cachePool, err := bigcache.New(context.Background(), bigcache.DefaultConfig(ttl))
	if err != nil {
		panic(fmt.Sprintf("create subtitle cache pool failed: %v", err))
	}
	cacheFunc := getCacheBaseFunc(cachePool, "字幕", reg.String())

	return func(ctx *gin.Context) {
		if ctx.Request.Method != http.MethodGet || !reg.MatchString(ctx.Request.URL.Path) {
			ctx.Next()
			return
		}
		cacheFunc(ctx)
	}
}
