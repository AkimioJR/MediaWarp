package handler

import (
	"MediaWarp/constants"
	"MediaWarp/internal/config"
	"errors"
	"net/http"
	"regexp"
)

// 媒体服务器处理接口
type MediaServerHandler interface {
	ReverseProxy(http.ResponseWriter, *http.Request) // 转发请求至上游服务器
	GetRegexpRouteRules() []RegexpRouteRule          // 获取正则路由表
	GetImageCacheRegexp() *regexp.Regexp             // 获取图片缓存正则表达式
}

var mediaServerHandler MediaServerHandler
var ErrInvalidMediaServerType = errors.New("错误的媒体服务器类型")

// 初始化媒体服务器处理器
func Init() error {
	var err error
	switch config.MediaServer.Type {
	case constants.EMBY:
		mediaServerHandler, err = NewEmbyServerHandler(config.MediaServer.ADDR, config.MediaServer.AUTH)
	case constants.JELLYFIN:
		mediaServerHandler, err = NewJellyfinHander(config.MediaServer.ADDR, config.MediaServer.AUTH)
	default:
		err = ErrInvalidMediaServerType
	}
	return err
}

// 获取媒体服务器接口
func GetMediaServer() MediaServerHandler {
	return mediaServerHandler
}
