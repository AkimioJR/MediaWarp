package router

import (
	"github.com/gin-gonic/gin"
)

type MiddlewareFunc func(internalFunc gin.HandlerFunc) gin.HandlerFunc

// MiddlewareChain 中间件链构建器
type MiddlewareChain struct {
	middlewares []MiddlewareFunc
}

// NewMiddlewareChain 创建新的中间件链
func NewMiddlewareChain(middlewares ...MiddlewareFunc) *MiddlewareChain {
	return &MiddlewareChain{
		middlewares: middlewares,
	}
}

// Add 添加中间件到链中
func (mc *MiddlewareChain) Add(middleware MiddlewareFunc) *MiddlewareChain {
	mc.middlewares = append(mc.middlewares, middleware)
	return mc
}

// Execute 执行中间件链
func (mc *MiddlewareChain) Execute(handler gin.HandlerFunc) gin.HandlerFunc {
	// 从最后一个中间件开始，向前包装
	result := handler
	for i := len(mc.middlewares) - 1; i >= 0; i-- {
		result = mc.middlewares[i](result)
	}
	return result
}
