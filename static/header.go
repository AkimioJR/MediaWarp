package static

// 通用头部

const basicPath = "/MediaWarp/static/"

const (
	WebCustomHeaderStart  = `<!-- MediaWarp Web 用户自定义额外 Header Start -->`
	WebCustomHeaderEnd    = `<!-- MediaWarp Web 用户自定义额外 Header End -->`
	WebEmbededHeaderStart = `<!-- MediaWarp Web 内嵌脚本插入 Start -->`
	WebEmbededHeaderEnd   = `<!-- MediaWarp Web 内嵌脚本插入 End -->`
)

// 通用嵌入脚本
const (
	VideoTogether     = `<script src="https://2gether.video/release/extension.website.user.js"></script>`
	ExternalPlayerUrl = `<script src="` + basicPath + `embyExternalUrl/embyWebAddExternalUrl/embyLaunchPotplayer.js"></script>`
	ActorPlus         = `<script src="` + basicPath + `emby-web-mod/actorPlus/actorPlus.js"></script>`
	FanartShow        = `<script src="` + basicPath + `emby-web-mod/fanart_show/fanart_show.js"></script>`
)

// Emby 相关嵌入脚本
const (
	EmbyCrx = `<link rel="stylesheet" id="theme-css" href="` + basicPath + `emby-crx/static/css/style.css" type="text/css" media="all" />
    <script src="` + basicPath + `emby-crx/static/js/common-utils.js"></script>
    <script src="` + basicPath + `emby-crx/static/js/jquery-3.6.0.min.js"></script>
    <script src="` + basicPath + `emby-crx/static/js/md5.min.js"></script>
    <script src="` + basicPath + `emby-crx/content/main.js"></script>`
	EmbyDanmaku = `<script src="` + basicPath + `dd-danmaku/ede.js" defer></script>`
)

// Jellyfin 相关嵌入脚本
const (
	JellyfinCrx = `<link rel="stylesheet" id="theme-css" href="` + basicPath + `jellyfin-crx/static/css/style.css" type="text/css" media="all" />
    <script src="` + basicPath + `jellyfin-crx/static/js/common-utils.js"></script>
    <script src="` + basicPath + `jellyfin-crx/static/js/jquery-3.6.0.min.js"></script>
    <script src="` + basicPath + `jellyfin-crx/static/js/md5.min.js"></script>
    <script src="` + basicPath + `jellyfin-crx/content/main.js"></script>`
	JellyfinDanmaku = `<script src="` + basicPath + `jellyfin-danmaku/ede.js" defer></script>`
)

// 飞牛影视相关嵌入脚本
const (
	FNTVPlayInfoHook = `<script src="` + basicPath + `fn-danmaku/fntv-play-info-hook.js"></script>`
	FNTVDanmakuOnly  = `<script src="` + basicPath + `fn-danmaku/fn-danmaku.js" defer></script>`
	FNTVDanmaku      = FNTVPlayInfoHook + "\n" + FNTVDanmakuOnly
)
