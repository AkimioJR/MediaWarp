package constants

type Color uint8

// 基础颜色枚举
const (
	ColorBlack  Color = iota // 黑色
	ColorRed                 // 红色
	ColorGreen               // 绿色
	ColorYellow              // 黄色
	ColorBlue                // 蓝色
	ColorPurple              // 紫色
	ColorCyan                // 青色
	ColorGray                // 灰色
)

// HTTP 状态码对应颜色
const (
	StatusCode200Color = ColorGreen  // HTTP 200 成功响应颜色
	StatusCode300Color = ColorGray   // HTTP 300 重定向颜色
	StatusCode400Color = ColorYellow // HTTP 400 客户端错误颜色
	StatusCode500Color = ColorRed    // HTTP 500 服务器错误颜色
)

// HTTP 方法对应颜色
const (
	MethodGetColor     = ColorBlue   // GET 方法颜色
	MethodPostColor    = ColorCyan   // POST 方法颜色
	MethodPutColor     = ColorYellow // PUT 方法颜色
	MethodPatchColor   = ColorGreen  // PATCH 方法颜色
	MethodDeleteColor  = ColorRed    // DELETE 方法颜色
	MethodHeadColor    = ColorPurple // HEAD 方法颜色
	MethodOptionsColor = ColorGray   // OPTIONS 方法颜色
)
