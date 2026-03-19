package middleware

import (
	"net/http"
	"strings"

	"github.com/AkimioJR/MediaWarp/constants"
	"github.com/AkimioJR/MediaWarp/internal/config"
	"github.com/AkimioJR/MediaWarp/internal/logging"

	"github.com/gin-gonic/gin"
)

// 客户端过滤器
func ClientFilter() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userAgent := ctx.Request.UserAgent()
		var allowed bool
		if userAgent == "" { // 开启了客户端过滤器后禁止所有未提供User-Agent的链接
			allowed = false
		} else {
			switch config.ClientFilter.Mode {
			case constants.WHITELIST: // 白名单模式
				allowed = false
				for _, ua := range config.ClientFilter.ClientList {
					if strings.Contains(userAgent, ua) {
						allowed = true
						break
					}
				}
			case constants.BLACKLIST: // 黑名单模式
				allowed = true
				for _, ua := range config.ClientFilter.ClientList {
					if strings.Contains(userAgent, ua) {
						allowed = false
						break
					}
				}
			}
		}
		if !allowed {
			ctx.AbortWithStatus(http.StatusForbidden) // 禁止访问
			logging.Info("客户端过滤器拦截了请求，User-Agent: ", userAgent)
			return
		}
		logging.Debug("客户端过滤器放行了请求，User-Agent: ", userAgent)
		ctx.Next()

	}
}
