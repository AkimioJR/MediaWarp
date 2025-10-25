package router

import "github.com/gin-gonic/gin"

func DisableCompression(internalFunc gin.HandlerFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Request.Header.Del("Accept-Encoding")
		internalFunc(ctx)
	}
}

var _ MiddlewareFunc = DisableCompression
