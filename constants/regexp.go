package constants

import "regexp"

type CacheRegexps struct {
	Image    *regexp.Regexp // 图片缓存匹配
	Subtitle *regexp.Regexp // 字幕缓存匹配
}

type EmbyRegexps struct {
	Router RouterRegexps
	Others OthersRegexps
	Cache  CacheRegexps
}

type RouterRegexps struct {
	VideosHandler        *regexp.Regexp // 普通视频处理接口匹配
	ModifyBaseHtmlPlayer *regexp.Regexp // 修改 Web 的 basehtmlplayer.js
	ModifyIndex          *regexp.Regexp // Web 首页
	ModifyPlaybackInfo   *regexp.Regexp // 播放信息处理接口
	ModifySubtitles      *regexp.Regexp // 字幕处理接口
}

type OthersRegexps struct {
	VideoRedirectReg *regexp.Regexp // 视频重定向匹配，统一视频请求格式
}

var EmbyRegexp = &EmbyRegexps{
	Router: RouterRegexps{
		VideosHandler:        regexp.MustCompile(`(?i)^(/emby)?/Videos/\d+/(stream|original)(\.\w+)?$`),
		ModifyBaseHtmlPlayer: regexp.MustCompile(`(?i)^/web/modules/htmlvideoplayer/basehtmlplayer.js$`),
		ModifyIndex:          regexp.MustCompile(`^/web/index.html$`),
		ModifyPlaybackInfo:   regexp.MustCompile(`(?i)^(/emby)?/Items/\d+/PlaybackInfo$`),
		ModifySubtitles:      regexp.MustCompile(`(?i)^(/emby)?/Videos/\d+/\w+/subtitles$`),
	},
	Others: OthersRegexps{
		VideoRedirectReg: regexp.MustCompile(`(?i)^(/emby)?/videos/(.*)/stream/(.*)`),
	},
	// /emby/Items/6/Images/Primary
	// /emby/Items/13/Images/Primary
	Cache: CacheRegexps{
		Image: regexp.MustCompile(`(?i)^(/emby)?/Items/\d+/Images`),
	},
}

type JellyfinRouterRegexps struct {
	VideosHandler      *regexp.Regexp // 普通视频处理接口匹配
	ModifyIndex        *regexp.Regexp // Web 首页
	ModifyPlaybackInfo *regexp.Regexp // 播放信息处理接口
	ModifySubtitles    *regexp.Regexp // 字幕处理接口
}
type JellyfinRegexps struct {
	Router JellyfinRouterRegexps
	Cache  CacheRegexps
}

var JellyfinRegexp = &JellyfinRegexps{
	Router: JellyfinRouterRegexps{
		VideosHandler:      regexp.MustCompile(`/Videos/\w+/(stream|original)(\.\w+)?$`), // /Videos/813a630bcf9c3f693a2ec8c498f868d2/stream /Videos/205953b114bb8c9dc2c7ba7e44b8024c/stream.mp4
		ModifyIndex:        regexp.MustCompile(`^/web/$`),
		ModifyPlaybackInfo: regexp.MustCompile(`^/Items/\w+/PlaybackInfo$`),
		ModifySubtitles:    regexp.MustCompile(`/Videos/\d+/\w+/subtitles$`),
	},
	Cache: CacheRegexps{
		// /Items/19ba9e43f0db12e2eea4294609ec1a0c/Images/Primary
		// /Items/20524938b33d516922ccea207555315b/Images/Backdrop/0
		Image: regexp.MustCompile(`(?i)/Items/\w+/Images`),
	},
}
