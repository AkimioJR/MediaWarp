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

	ginR.NoRoute(RegexpRouterHandler)
	return ginR
}

var middlewareChain = NewMiddlewareChain().
	Add(QueryKeyCaseInsensitive).
	Add(DisableCompression)

// 正则表达式路由处理器
//
// 从媒体服务器处理结构体中获取正则路由规则
// 依次匹配请求, 找到对应的处理器
func RegexpRouterHandler(ctx *gin.Context) {
	mediaServerHandler := handler.GetMediaServer()

	for _, rule := range mediaServerHandler.GetRegexpRouteRules() {
		if rule.Regexp.MatchString(ctx.Request.URL.Path) { // 不带查询参数的字符串：/emby/Items/54/Images/Primary
			logging.Debugf("URL: %s 匹配成功 -> %s", ctx.Request.URL.Path, rule.Regexp.String())

			middlewareChain.Execute(rule.Handler)(ctx)
			return
		}
	}

	// 未匹配路由
	mediaServerHandler.ReverseProxy(ctx.Writer, ctx.Request)
}
