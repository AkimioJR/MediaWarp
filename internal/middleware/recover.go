package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/AkimioJR/MediaWarp/internal/logging"

	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if r := recover(); r != nil && r != http.ErrAbortHandler { // 忽略 http.ErrAbortHandler 的 panic（httputil.ReverseProxy 的 panic）
				stack := debug.Stack()
				logging.Errorf("[Recovery] %s panic revocered: %v\n%s", ctx.Request.URL.Path, r, string(stack))
				ctx.AbortWithStatusJSON(
					http.StatusInternalServerError,
					gin.H{
						"error": "MediaWarp Internal Server Error",
						"msg":   r,
					},
				)
			}
		}()
		ctx.Next()
	}
}
