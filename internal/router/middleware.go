package router

import (
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

type MiddlewareFunc func(internalFunc gin.HandlerFunc) gin.HandlerFunc

func QueryKeyCaseInsensitive(internalFunc gin.HandlerFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		queryParams := make(url.Values)
		for key, values := range ctx.Request.URL.Query() {
			queryParams.Add(strings.ToLower(key), strings.Join(values, ","))
		}
		ctx.Request.URL.RawQuery = queryParams.Encode()

		internalFunc(ctx)
	}
}

var _ MiddlewareFunc = QueryKeyCaseInsensitive

func DisableCompression(internalFunc gin.HandlerFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Request.Header.Del("Accept-Encoding")
		internalFunc(ctx)
	}
}

var _ MiddlewareFunc = DisableCompression
