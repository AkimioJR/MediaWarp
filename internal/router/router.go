package router

import (
	"MediaWarp/constants"
	"MediaWarp/internal/config"
	"MediaWarp/internal/handler"
	"MediaWarp/internal/logging"
	"MediaWarp/internal/middleware"
	"MediaWarp/static"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	ginR := gin.New()
	ginR.Use(
		middleware.Logger(),
		middleware.Recovery(),
		middleware.SetRefererPolicy(constants.SameOrigin),
	)

	if config.ClientFilter.Enable {
		ginR.Use(middleware.ClientFilter())
		logging.Info("客户端过滤中间件已启用")
	} else {
		logging.Info("客户端过滤中间件未启用")
	}

	mediawarpRouter := ginR.Group("/MediaWarp")
	{
		mediawarpRouter.Any("/version", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, config.Version())
		})
		if config.Web.Enable { // 启用 Web 页面修改相关设置
			mediawarpRouter.StaticFS("/static", http.FS(static.EmbeddedStaticAssets))
			if config.Web.Custom { // 用户自定义静态资源目录
				mediawarpRouter.Static("/custom", config.CostomDir())
			}
			if config.Web.Robots != "" { // 自定义 robots.txt
				ginR.GET(
					"/robots.txt",
					func(ctx *gin.Context) {
						ctx.String(http.StatusOK, config.Web.Robots)
					},
				)
			}
		}
	}

	handlers := make(gin.HandlersChain, 0, 3)
	if config.Cache.Enable {
		mediaServerHandler := handler.GetMediaServer()
		{
			if config.Cache.ImageTTL > 0 {
				if mediaServerHandler.GetImageCacheRegexp() != nil {
					logging.Infof("图片缓存中间件已启用, TTL: %s", config.Cache.ImageTTL.String())
					handlers = append(handlers, middleware.ImageCache(config.Cache.ImageTTL, mediaServerHandler.GetImageCacheRegexp()))
				} else {
					logging.Warningf("媒体服务器 %s 不支持图片缓存, 未添加图片缓存中间件", config.MediaServer.Type.String())
				}
			} else {
				logging.Infof("图片缓存中间件未启用, TTL: %s", config.Cache.ImageTTL.String())
			}
		}

		{
			if config.Cache.SubtitleTTL > 0 {
				if mediaServerHandler.GetSubtitleCacheRegexp() != nil {
					logging.Infof("字幕缓存中间件已启用, TTL: %s", config.Cache.SubtitleTTL.String())
					handlers = append(handlers, middleware.SubtitleCache(config.Cache.SubtitleTTL, mediaServerHandler.GetSubtitleCacheRegexp()))
				} else {
					logging.Warningf("媒体服务器 %s 不支持字幕缓存, 未添加字幕缓存中间件", config.MediaServer.Type.String())
				}
			} else {
				logging.Infof("字幕缓存中间件未启用, TTL: %s", config.Cache.SubtitleTTL.String())
			}
		}
	} else {
		logging.Info("全局缓存未启用, 未添加缓存中间件")
	}

	handlers = append(handlers, getRegexpRouterHandler())
	ginR.NoRoute(handlers...)
	return ginR
}

// 正则表达式路由处理器
//
// 从媒体服务器处理结构体中获取正则路由规则
// 依次匹配请求, 找到对应的处理器
func getRegexpRouterHandler() gin.HandlerFunc {
	mediaServerHandler := handler.GetMediaServer()
	middlewareChain := NewMiddlewareChain().
		Add(QueryKeyCaseInsensitive).
		Add(DisableCompression)

	return func(ctx *gin.Context) {
		for _, rule := range mediaServerHandler.GetRegexpRouteRules() {
			if rule.Regexp.MatchString(ctx.Request.URL.Path) { // 不带查询参数的字符串：/emby/Items/54/Images/Primary
				logging.AccessDebugf(ctx, "匹配成功正则表达式: %s", rule.Regexp.String())

				middlewareChain.Execute(rule.Handler)(ctx)
				return
			}
		}

		// 未匹配路由
		mediaServerHandler.ReverseProxy(ctx.Writer, ctx.Request)
	}
}
