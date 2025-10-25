package config

import (
	"MediaWarp/constants"
	"time"
)

// 程序版本信息
type VersionInfo struct {
	AppVersion string `json:"app_version"` // 程序版本号
	CommitHash string `json:"commit_hash"` // GIt Commit Hash
	BuildDate  string `json:"build_date"`  // 编译时间
	GoVersion  string `json:"go_version"`  // 编译 Golang 版本
	OS         string `json:"os"`          // 操作系统
	Arch       string `json:"arch"`        //  架构
}

// 上游媒体服务器相关设置
type MediaServerSetting struct {
	Type constants.MediaServerType `yaml:"type"` // 媒体服务器类型
	ADDR string                    `yaml:"addr"` // 地址
	AUTH string                    `yaml:"auth"` // 认证授权KEY
}

// 日志设置
type LoggerSetting struct {
	AccessLogger  BaseLoggerSetting `yaml:"access"`  // 访问日志相关配置
	ServiceLogger BaseLoggerSetting `yaml:"service"` // 服务日志相关配置
}

// 基础日志配置字段
type BaseLoggerSetting struct {
	Console bool `yaml:"console"` // 是否将日志输出到终端中
	File    bool `yaml:"file"`    // 是否将日志输出到文件中
}

// 缓存设置
type CacheSetting struct {
	Enable      bool          `yaml:"enable"`
	HTTPStrmTTL time.Duration `yaml:"http_strm_ttl"`
	AlistAPITTL time.Duration `yaml:"alist_api_ttl"`
	ImageTTL    time.Duration `yaml:"image_ttl"`
}

// Web前端自定义设置
type WebSetting struct {
	Enable            bool   `yaml:"enable"`              // 启用自定义前端设置
	Custom            bool   `yaml:"custom"`              // 启用用户自定义静态资源
	Index             bool   `yaml:"index"`               // 是否从 custom 目录读取 index.html 文件作为首页
	Head              string `yaml:"head"`                // 添加到 index.html 的 HEAD 中
	Robots            string `yaml:"robots"`              // 自定义 robots.txt，若为空表示不修改
	ExternalPlayerUrl bool   `yaml:"external_player_url"` // 是否开启外置播放器
	Crx               bool   `yaml:"crx"`                 // crx 美化
	ActorPlus         bool   `yaml:"actor_plus"`          // 过滤没有头像的演员和制作人员
	FanartShow        bool   `yaml:"fanart_show"`         // 显示同人图（fanart图）
	Danmaku           bool   `yaml:"danmaku"`             // Web 弹幕
	VideoTogether     bool   `yaml:"video_together"`      // VideoTogether
}

// 客户端User-Agent过滤设置
type ClientFilterSetting struct {
	Enable     bool                 `yaml:"enable"`
	Mode       constants.FliterMode `yaml:"mode"`
	ClientList []string             `yaml:"list"`
}

// HTTPStrm播放设置
type HTTPStrmSetting struct {
	Enable     bool     `yaml:"enable"`
	TransCode  bool     `yaml:"transcode"` // false->强制关闭转码 true->保持原有转码设置
	FinalURL   bool     `yaml:"final_url"` // 对 URL 进行重定向判断，找到非重定向地址再重定向给客户端，减少客户端重定向次数
	PrefixList []string `yaml:"prefix_list"`
}

// AlistStrm具体设置
type AlistSetting struct {
	ADDR       string   `yaml:"addr"`
	Username   string   `yaml:"username"`
	Password   string   `yaml:"password"`
	Token      *string  `yaml:"token"`
	PrefixList []string `yaml:"prefix_list"`
}

// AlistStrm播放设置
type AlistStrmSetting struct {
	Enable    bool           `yaml:"enable"`
	TransCode bool           `yaml:"transcode"` // false->强制关闭转码 true->保持原有转码设置
	RawURL    bool           `yaml:"raw_url"`   // 是否使用原始 URL
	List      []AlistSetting `yaml:"list"`
}

// 字幕设置
type SubtitleSetting struct {
	Enable   bool     `yaml:"enable"`
	SRT2ASS  bool     `yaml:"srt2ass"` // SRT 字幕转 ASS 字幕
	ASSStyle []string `yaml:"ass_style"`
	SubSet   bool     `yaml:"subset"` // ASS 字幕字体子集化
}

type Setting struct {
	Port         uint16              `yaml:"port"`
	MediaServer  MediaServerSetting  `yaml:"server"`
	Logger       LoggerSetting       `yaml:"log"`
	Cache        CacheSetting        `yaml:"cache"`
	Web          WebSetting          `yaml:"web"`
	ClientFilter ClientFilterSetting `yaml:"client"`
	HTTPStrm     HTTPStrmSetting     `yaml:"http_strm"`
	AlistStrm    AlistStrmSetting    `yaml:"alist_strm"`
	Subtitle     SubtitleSetting     `yaml:"subtitle"`
}
