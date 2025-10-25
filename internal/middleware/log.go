package middleware

import (
	"MediaWarp/constants"
	"MediaWarp/internal/logging"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// 记录访问日志
func Logger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		method := ctx.Request.Method
		path := ctx.Request.URL.Path
		query := ctx.Request.URL.RawQuery
		if query != "" {
			path = path + "?" + query
		}

		startTime := time.Now()
		ctx.Next()
		wasteTime := time.Since(startTime)

		clientIP := ctx.ClientIP()
		statusCode := ctx.Writer.Status()

		statusColor, methodColor := getColor(statusCode, method)

		logging.AccessLog(
			"【Access】 %s |\033[4%dm %d \033[0m| %-10s |\033[4%dm %-7s \033[0m| %s \"%s\"",
			startTime.Format(time.DateTime),
			statusColor, statusCode,
			wasteTime,
			methodColor, method,
			clientIP,
			path,
		)
	}
}

// 根据Http状态码和Http请求方法获取颜色
func getColor(statusCode int, method string) (constants.Color, constants.Color) {
	var statusColor, methodColor constants.Color
	switch {
	case statusCode >= http.StatusOK && statusCode < http.StatusMultipleChoices:
		statusColor = constants.StatusCode200Color
	case statusCode >= http.StatusMultipleChoices && statusCode < http.StatusBadRequest:
		statusColor = constants.StatusCode300Color
	case statusCode >= http.StatusBadRequest && statusCode < http.StatusInternalServerError:
		statusColor = constants.StatusCode400Color
	case statusCode >= http.StatusInternalServerError:
		statusColor = constants.StatusCode500Color
	default:
		statusColor = constants.ColorBlack
	}
	switch method {
	case http.MethodGet:
		methodColor = constants.MethodGetColor
	case http.MethodPost:
		methodColor = constants.MethodPostColor
	case http.MethodPut:
		methodColor = constants.MethodPutColor
	case http.MethodPatch:
		methodColor = constants.MethodPatchColor
	case http.MethodDelete:
		methodColor = constants.MethodDeleteColor
	case http.MethodHead:
		methodColor = constants.MethodHeadColor
	case http.MethodOptions:
		methodColor = constants.MethodOptionsColor
	default:
		methodColor = constants.ColorBlack
	}
	return statusColor, methodColor
}
