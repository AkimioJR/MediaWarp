package emby

type EmbyResponse struct {
	Items            []BaseItemDto `json:"Items,omitempty"`
	TotalRecordCount *int64        `json:"TotalRecordCount,omitempty"`
}

// /Items/:itemID/PlaybackInfo的响应
type PlaybackInfoResponse struct {
	ErrorCode     *PlaybackErrorCode `json:"ErrorCode,omitempty"`
	MediaSources  []MediaSourceInfo  `json:"MediaSources,omitempty"`
	PlaySessionID *string            `json:"PlaySessionId,omitempty"`
	Chapters      []ChapterInfo      `json:"Chapters,omitempty"`
}

// PlaybackErrorCode
type PlaybackErrorCode string

const (
	NoCompatibleStream PlaybackErrorCode = "NoCompatibleStream"
	NotAllowed         PlaybackErrorCode = "NotAllowed"
	RateLimitExceeded  PlaybackErrorCode = "RateLimitExceeded"
)

// BaseItemDto
type BaseItemDto struct {
	AffiliateCallSign            *string                  `json:"AffiliateCallSign,omitempty"`
	AirDays                      []DayOfWeek              `json:"AirDays,omitempty"`
	Album                        *string                  `json:"Album,omitempty"`
	AlbumArtist                  *string                  `json:"AlbumArtist,omitempty"`
	AlbumArtists                 []NameIDPair             `json:"AlbumArtists,omitempty"`
	AlbumCount                   *int64                   `json:"AlbumCount"`
	AlbumID                      *string                  `json:"AlbumId,omitempty"`
	AlbumPrimaryImageTag         *string                  `json:"AlbumPrimaryImageTag,omitempty"`
	Altitude                     *float64                 `json:"Altitude"`
	Aperture                     *float64                 `json:"Aperture"`
	ArtistItems                  []NameIDPair             `json:"ArtistItems,omitempty"`
	Artists                      []string                 `json:"Artists,omitempty"`
	AsSeries                     *bool                    `json:"AsSeries"`
	BackdropImageTags            []string                 `json:"BackdropImageTags,omitempty"`
	Bitrate                      *int64                   `json:"Bitrate"`
	CameraMake                   *string                  `json:"CameraMake,omitempty"`
	CameraModel                  *string                  `json:"CameraModel,omitempty"`
	CanDelete                    *bool                    `json:"CanDelete"`
	CanDownload                  *bool                    `json:"CanDownload"`
	CanEditItems                 *bool                    `json:"CanEditItems"`
	CanLeaveContent              *bool                    `json:"CanLeaveContent"`
	CanMakePublic                *bool                    `json:"CanMakePublic"`
	CanManageAccess              *bool                    `json:"CanManageAccess"`
	ChannelID                    *string                  `json:"ChannelId,omitempty"`
	ChannelName                  *string                  `json:"ChannelName,omitempty"`
	ChannelNumber                *string                  `json:"ChannelNumber,omitempty"`
	ChannelPrimaryImageTag       *string                  `json:"ChannelPrimaryImageTag,omitempty"`
	Chapters                     []ChapterInfo            `json:"Chapters,omitempty"`
	ChildCount                   *int64                   `json:"ChildCount"`
	CollectionType               *string                  `json:"CollectionType,omitempty"`
	CommunityRating              *float64                 `json:"CommunityRating"`
	CompletionPercentage         *float64                 `json:"CompletionPercentage"`
	Composers                    []NameIDPair             `json:"Composers,omitempty"`
	Container                    *string                  `json:"Container,omitempty"`
	CriticRating                 *float64                 `json:"CriticRating"`
	CurrentProgram               *BaseItemDto             `json:"CurrentProgram,omitempty"`
	CustomRating                 *string                  `json:"CustomRating,omitempty"`
	DateCreated                  *string                  `json:"DateCreated"`
	Disabled                     *bool                    `json:"Disabled"`
	DisplayOrder                 *string                  `json:"DisplayOrder,omitempty"`
	DisplayPreferencesID         *string                  `json:"DisplayPreferencesId,omitempty"`
	EndDate                      *string                  `json:"EndDate"`
	EpisodeTitle                 *string                  `json:"EpisodeTitle,omitempty"`
	Etag                         *string                  `json:"Etag,omitempty"`
	ExposureTime                 *float64                 `json:"ExposureTime"`
	ExternalUrls                 []ExternalURL            `json:"ExternalUrls,omitempty"`
	ExtraType                    *string                  `json:"ExtraType,omitempty"`
	FileName                     *string                  `json:"FileName,omitempty"`
	FocalLength                  *float64                 `json:"FocalLength"`
	ForcedSortName               *string                  `json:"ForcedSortName,omitempty"`
	GameSystem                   *string                  `json:"GameSystem,omitempty"`
	GameSystemID                 *int64                   `json:"GameSystemId"`
	GenreItems                   []NameLongIDPair         `json:"GenreItems,omitempty"`
	Genres                       []string                 `json:"Genres,omitempty"`
	GUID                         *string                  `json:"Guid,omitempty"`
	Height                       *int64                   `json:"Height"`
	ID                           *string                  `json:"Id,omitempty"`
	ImageOrientation             *DrawingImageOrientation `json:"ImageOrientation,omitempty"`
	ImageTags                    map[string]string        `json:"ImageTags,omitempty"`
	IndexNumber                  *int64                   `json:"IndexNumber"`
	IndexNumberEnd               *int64                   `json:"IndexNumberEnd"`
	IsFolder                     *bool                    `json:"IsFolder"`
	IsKids                       *bool                    `json:"IsKids"`
	IsLive                       *bool                    `json:"IsLive"`
	IsMovie                      *bool                    `json:"IsMovie"`
	IsNew                        *bool                    `json:"IsNew"`
	IsNews                       *bool                    `json:"IsNews"`
	ISOSpeedRating               *int64                   `json:"IsoSpeedRating"`
	IsPremiere                   *bool                    `json:"IsPremiere"`
	IsRepeat                     *bool                    `json:"IsRepeat"`
	IsSeries                     *bool                    `json:"IsSeries"`
	IsSports                     *bool                    `json:"IsSports"`
	Latitude                     *float64                 `json:"Latitude"`
	ListingsChannelID            *string                  `json:"ListingsChannelId,omitempty"`
	ListingsChannelName          *string                  `json:"ListingsChannelName,omitempty"`
	ListingsChannelNumber        *string                  `json:"ListingsChannelNumber,omitempty"`
	ListingsID                   *string                  `json:"ListingsId,omitempty"`
	ListingsPath                 *string                  `json:"ListingsPath,omitempty"`
	ListingsProviderID           *string                  `json:"ListingsProviderId,omitempty"`
	LocalTrailerCount            *int64                   `json:"LocalTrailerCount"`
	LocationType                 *LocationType            `json:"LocationType,omitempty"`
	LockData                     *bool                    `json:"LockData"`
	LockedFields                 []MetadataFields         `json:"LockedFields,omitempty"`
	Longitude                    *float64                 `json:"Longitude"`
	ManagementID                 *string                  `json:"ManagementId,omitempty"`
	MediaSources                 []MediaSourceInfo        `json:"MediaSources,omitempty"`
	MediaStreams                 []MediaStream            `json:"MediaStreams,omitempty"`
	MediaType                    *string                  `json:"MediaType,omitempty"`
	MovieCount                   *int64                   `json:"MovieCount"`
	MusicVideoCount              *int64                   `json:"MusicVideoCount"`
	Name                         *string                  `json:"Name,omitempty"`
	Number                       *string                  `json:"Number,omitempty"`
	OfficialRating               *string                  `json:"OfficialRating,omitempty"`
	OriginalTitle                *string                  `json:"OriginalTitle,omitempty"`
	Overview                     *string                  `json:"Overview,omitempty"`
	ParentBackdropImageTags      []string                 `json:"ParentBackdropImageTags,omitempty"`
	ParentBackdropItemID         *string                  `json:"ParentBackdropItemId,omitempty"`
	ParentID                     *string                  `json:"ParentId,omitempty"`
	ParentIndexNumber            *int64                   `json:"ParentIndexNumber"`
	ParentLogoImageTag           *string                  `json:"ParentLogoImageTag,omitempty"`
	ParentLogoItemID             *string                  `json:"ParentLogoItemId,omitempty"`
	ParentThumbImageTag          *string                  `json:"ParentThumbImageTag,omitempty"`
	ParentThumbItemID            *string                  `json:"ParentThumbItemId,omitempty"`
	PartCount                    *int64                   `json:"PartCount"`
	Path                         *string                  `json:"Path,omitempty"`
	People                       []BaseItemPerson         `json:"People,omitempty"`
	PlaylistItemID               *string                  `json:"PlaylistItemId,omitempty"`
	PreferredMetadataCountryCode *string                  `json:"PreferredMetadataCountryCode,omitempty"`
	PreferredMetadataLanguage    *string                  `json:"PreferredMetadataLanguage,omitempty"`
	Prefix                       *string                  `json:"Prefix,omitempty"`
	PremiereDate                 *string                  `json:"PremiereDate"`
	PresentationUniqueKey        *string                  `json:"PresentationUniqueKey,omitempty"`
	PrimaryImageAspectRatio      *float64                 `json:"PrimaryImageAspectRatio"`
	PrimaryImageItemID           *string                  `json:"PrimaryImageItemId,omitempty"`
	PrimaryImageTag              *string                  `json:"PrimaryImageTag,omitempty"`
	ProductionLocations          []string                 `json:"ProductionLocations,omitempty"`
	ProductionYear               *int64                   `json:"ProductionYear"`
	ProviderIDS                  map[string]string        `json:"ProviderIds,omitempty"`
	RecursiveItemCount           *int64                   `json:"RecursiveItemCount"`
	RemoteTrailers               []MediaURL               `json:"RemoteTrailers,omitempty"`
	RunTimeTicks                 *int64                   `json:"RunTimeTicks"`
	SeasonID                     *string                  `json:"SeasonId,omitempty"`
	SeasonName                   *string                  `json:"SeasonName,omitempty"`
	SeriesCount                  *int64                   `json:"SeriesCount"`
	SeriesID                     *string                  `json:"SeriesId,omitempty"`
	SeriesName                   *string                  `json:"SeriesName,omitempty"`
	SeriesPrimaryImageTag        *string                  `json:"SeriesPrimaryImageTag,omitempty"`
	SeriesStudio                 *string                  `json:"SeriesStudio,omitempty"`
	SeriesTimerID                *string                  `json:"SeriesTimerId,omitempty"`
	ServerID                     *string                  `json:"ServerId,omitempty"`
	ShutterSpeed                 *float64                 `json:"ShutterSpeed"`
	Size                         *int64                   `json:"Size"`
	Software                     *string                  `json:"Software,omitempty"`
	SongCount                    *int64                   `json:"SongCount"`
	SortIndexNumber              *int64                   `json:"SortIndexNumber"`
	SortName                     *string                  `json:"SortName,omitempty"`
	SortParentIndexNumber        *int64                   `json:"SortParentIndexNumber"`
	SpecialFeatureCount          *int64                   `json:"SpecialFeatureCount"`
	StartDate                    *string                  `json:"StartDate"`
	Status                       *string                  `json:"Status,omitempty"`
	Studios                      []NameLongIDPair         `json:"Studios,omitempty"`
	Subviews                     []string                 `json:"Subviews,omitempty"`
	SupportsResume               *bool                    `json:"SupportsResume"`
	SupportsSync                 *bool                    `json:"SupportsSync"`
	SyncStatus                   *SyncJobItemStatus       `json:"SyncStatus,omitempty"`
	TagItems                     []NameLongIDPair         `json:"TagItems,omitempty"`
	Taglines                     []string                 `json:"Taglines,omitempty"`
	Tags                         []string                 `json:"Tags,omitempty"`
	TimerID                      *string                  `json:"TimerId,omitempty"`
	TimerType                    *LiveTvTimerType         `json:"TimerType,omitempty"`
	Type                         *string                  `json:"Type,omitempty"`
	UserData                     *UserItemDataDto         `json:"UserData,omitempty"`
	Video3DFormat                *Video3DFormat           `json:"Video3DFormat,omitempty"`
	Width                        *int64                   `json:"Width"`
}

// NameIdPair
type NameIDPair struct {
	ID   *string `json:"Id,omitempty"`
	Name *string `json:"Name,omitempty"`
}

// ChapterInfo
type ChapterInfo struct {
	ChapterIndex       *int64      `json:"ChapterIndex,omitempty"`
	ImageTag           *string     `json:"ImageTag,omitempty"`
	MarkerType         *MarkerType `json:"MarkerType,omitempty"`
	Name               *string     `json:"Name,omitempty"`
	StartPositionTicks *int64      `json:"StartPositionTicks,omitempty"`
}

// ExternalUrl
type ExternalURL struct {
	Name *string `json:"Name,omitempty"`
	URL  *string `json:"Url,omitempty"`
}

// NameLongIdPair
type NameLongIDPair struct {
	ID   *int64  `json:"Id,omitempty"`
	Name *string `json:"Name,omitempty"`
}

// MediaSourceInfo
type MediaSourceInfo struct {
	AddAPIKeyToDirectStreamURL *bool                     `json:"AddApiKeyToDirectStreamUrl,omitempty"`
	AnalyzeDurationMS          *int64                    `json:"AnalyzeDurationMs,omitempty"`
	Bitrate                    *int64                    `json:"Bitrate,omitempty"`
	BufferMS                   *int64                    `json:"BufferMs,omitempty"`
	Chapters                   []ChapterInfo             `json:"Chapters,omitempty"`
	Container                  *string                   `json:"Container,omitempty"` // AlistStrm 显示 strm，普通视频和 HTTPStrm
	ContainerStartTimeTicks    *int64                    `json:"ContainerStartTimeTicks,omitempty"`
	DefaultAudioStreamIndex    *int64                    `json:"DefaultAudioStreamIndex,omitempty"`
	DefaultSubtitleStreamIndex *int64                    `json:"DefaultSubtitleStreamIndex,omitempty"`
	DirectStreamURL            *string                   `json:"DirectStreamUrl,omitempty"`
	EncoderPath                *string                   `json:"EncoderPath,omitempty"`
	EncoderProtocol            *MediaProtocol            `json:"EncoderProtocol,omitempty"`
	Formats                    []string                  `json:"Formats,omitempty"`
	HasMixedProtocols          *bool                     `json:"HasMixedProtocols,omitempty"`
	ID                         *string                   `json:"Id,omitempty"` // mediasource_45
	IsInfiniteStream           *bool                     `json:"IsInfiniteStream,omitempty"`
	IsRemote                   *bool                     `json:"IsRemote,omitempty"` // HTTPStrm 会被设置为 true
	ItemID                     *string                   `json:"ItemId,omitempty"`
	LiveStreamID               *string                   `json:"LiveStreamId,omitempty"`
	MediaStreams               []MediaStream             `json:"MediaStreams"` // 避免当 MediaStreams 为空数组的情况移除该字段导致导致部分客户端（目前已知：Yamby）报错
	Name                       *string                   `json:"Name,omitempty"`
	OpenToken                  *string                   `json:"OpenToken,omitempty"`
	Path                       *string                   `json:"Path,omitempty"` // 本地视频文件则是正常的本地路径，Strm 则是 Strm 文件的内容
	ProbePath                  *string                   `json:"ProbePath,omitempty"`
	ProbeProtocol              *MediaProtocol            `json:"ProbeProtocol,omitempty"`
	Protocol                   *MediaProtocol            `json:"Protocol,omitempty"`
	ReadAtNativeFramerate      *bool                     `json:"ReadAtNativeFramerate,omitempty"`
	RequiredHTTPHeaders        map[string]string         `json:"RequiredHttpHeaders,omitempty"`
	RequiresClosing            *bool                     `json:"RequiresClosing,omitempty"`
	RequiresLooping            *bool                     `json:"RequiresLooping,omitempty"`
	RequiresOpening            *bool                     `json:"RequiresOpening,omitempty"`
	RunTimeTicks               *int64                    `json:"RunTimeTicks,omitempty"`
	ServerID                   *string                   `json:"ServerId,omitempty"`
	Size                       *int64                    `json:"Size,omitempty"`
	SortName                   *string                   `json:"SortName,omitempty"`
	SupportsDirectPlay         *bool                     `json:"SupportsDirectPlay,omitempty"`
	SupportsDirectStream       *bool                     `json:"SupportsDirectStream,omitempty"`
	SupportsProbing            *bool                     `json:"SupportsProbing,omitempty"`
	SupportsTranscoding        *bool                     `json:"SupportsTranscoding,omitempty"`
	Timestamp                  *TransportStreamTimestamp `json:"Timestamp,omitempty"`
	TrancodeLiveStartIndex     *int64                    `json:"TrancodeLiveStartIndex,omitempty"`
	TranscodingContainer       *string                   `json:"TranscodingContainer,omitempty"`
	TranscodingSubProtocol     *string                   `json:"TranscodingSubProtocol,omitempty"`
	TranscodingURL             *string                   `json:"TranscodingUrl,omitempty"`
	Type                       *MediaSourceType          `json:"Type,omitempty"`
	Video3DFormat              *Video3DFormat            `json:"Video3DFormat,omitempty"`
	WallClockStart             *string                   `json:"WallClockStart,omitempty"`
}

// MediaStream
type MediaStream struct {
	AspectRatio                     *string                 `json:"AspectRatio,omitempty"`
	AttachmentSize                  *int64                  `json:"AttachmentSize"`
	AverageFrameRate                *float64                `json:"AverageFrameRate"`
	BitDepth                        *int64                  `json:"BitDepth"`
	BitRate                         *int64                  `json:"BitRate"`
	ChannelLayout                   *string                 `json:"ChannelLayout,omitempty"`
	Channels                        *int64                  `json:"Channels"`
	Codec                           *string                 `json:"Codec,omitempty"`
	CodecTag                        *string                 `json:"CodecTag,omitempty"`
	ColorPrimaries                  *string                 `json:"ColorPrimaries,omitempty"`
	ColorSpace                      *string                 `json:"ColorSpace,omitempty"`
	ColorTransfer                   *string                 `json:"ColorTransfer,omitempty"`
	Comment                         *string                 `json:"Comment,omitempty"`
	DeliveryMethod                  *SubtitleDeliveryMethod `json:"DeliveryMethod,omitempty"`
	DeliveryURL                     *string                 `json:"DeliveryUrl,omitempty"`
	DisplayLanguage                 *string                 `json:"DisplayLanguage,omitempty"`
	DisplayTitle                    *string                 `json:"DisplayTitle,omitempty"`
	ExtendedVideoSubType            *ExtendedVideoSubTypes  `json:"ExtendedVideoSubType,omitempty"`
	ExtendedVideoSubTypeDescription *string                 `json:"ExtendedVideoSubTypeDescription,omitempty"`
	ExtendedVideoType               *ExtendedVideoTypes     `json:"ExtendedVideoType,omitempty"`
	Extradata                       *string                 `json:"Extradata,omitempty"`
	Height                          *int64                  `json:"Height"`
	Index                           *int64                  `json:"Index,omitempty"`
	IsAnamorphic                    *bool                   `json:"IsAnamorphic"`
	IsAVC                           *bool                   `json:"IsAVC"`
	IsDefault                       *bool                   `json:"IsDefault,omitempty"`
	IsExternal                      *bool                   `json:"IsExternal,omitempty"`
	IsExternalURL                   *bool                   `json:"IsExternalUrl"`
	IsForced                        *bool                   `json:"IsForced,omitempty"`
	IsHearingImpaired               *bool                   `json:"IsHearingImpaired,omitempty"`
	IsInterlaced                    *bool                   `json:"IsInterlaced,omitempty"`
	IsTextSubtitleStream            *bool                   `json:"IsTextSubtitleStream,omitempty"`
	ItemID                          *string                 `json:"ItemId,omitempty"`
	Language                        *string                 `json:"Language,omitempty"`
	Level                           *float64                `json:"Level"`
	MIMEType                        *string                 `json:"MimeType,omitempty"`
	NalLengthSize                   *string                 `json:"NalLengthSize,omitempty"`
	Path                            *string                 `json:"Path,omitempty"`
	PixelFormat                     *string                 `json:"PixelFormat,omitempty"`
	Profile                         *string                 `json:"Profile,omitempty"`
	Protocol                        *MediaProtocol          `json:"Protocol,omitempty"`
	RealFrameRate                   *float64                `json:"RealFrameRate"`
	RefFrames                       *int64                  `json:"RefFrames"`
	Rotation                        *int64                  `json:"Rotation"`
	SampleRate                      *int64                  `json:"SampleRate"`
	ServerID                        *string                 `json:"ServerId,omitempty"`
	StreamStartTimeTicks            *int64                  `json:"StreamStartTimeTicks"`
	SubtitleLocationType            *SubtitleLocationType   `json:"SubtitleLocationType,omitempty"`
	SupportsExternalStream          *bool                   `json:"SupportsExternalStream,omitempty"`
	TimeBase                        *string                 `json:"TimeBase,omitempty"`
	Title                           *string                 `json:"Title,omitempty"`
	Type                            *MediaStreamType        `json:"Type,omitempty"`
	VideoRange                      *string                 `json:"VideoRange,omitempty"`
	Width                           *int64                  `json:"Width"`
}

// BaseItemPerson
type BaseItemPerson struct {
	ID              *string     `json:"Id,omitempty"`
	Name            *string     `json:"Name,omitempty"`
	PrimaryImageTag *string     `json:"PrimaryImageTag,omitempty"`
	Role            *string     `json:"Role,omitempty"`
	Type            *PersonType `json:"Type,omitempty"`
}

// MediaUrl
type MediaURL struct {
	Name *string `json:"Name,omitempty"`
	URL  *string `json:"Url,omitempty"`
}

// UserItemDataDto
type UserItemDataDto struct {
	IsFavorite            *bool    `json:"IsFavorite,omitempty"`
	ItemID                *string  `json:"ItemId,omitempty"`
	Key                   *string  `json:"Key,omitempty"`
	LastPlayedDate        *string  `json:"LastPlayedDate"`
	PlaybackPositionTicks *int64   `json:"PlaybackPositionTicks,omitempty"`
	PlayCount             *int64   `json:"PlayCount"`
	Played                *bool    `json:"Played,omitempty"`
	PlayedPercentage      *float64 `json:"PlayedPercentage"`
	Rating                *float64 `json:"Rating"`
	ServerID              *string  `json:"ServerId,omitempty"`
	UnplayedItemCount     *int64   `json:"UnplayedItemCount"`
}

// DayOfWeek
type DayOfWeek string

const (
	Friday    DayOfWeek = "Friday"
	Monday    DayOfWeek = "Monday"
	Saturday  DayOfWeek = "Saturday"
	Sunday    DayOfWeek = "Sunday"
	Thursday  DayOfWeek = "Thursday"
	Tuesday   DayOfWeek = "Tuesday"
	Wednesday DayOfWeek = "Wednesday"
)

// MarkerType
type MarkerType string

const (
	Chapter      MarkerType = "Chapter"
	CreditsStart MarkerType = "CreditsStart"
	IntroEnd     MarkerType = "IntroEnd"
	IntroStart   MarkerType = "IntroStart"
)

// Drawing.ImageOrientation
type DrawingImageOrientation string

const (
	BottomLeft  DrawingImageOrientation = "BottomLeft"
	BottomRight DrawingImageOrientation = "BottomRight"
	LeftBottom  DrawingImageOrientation = "LeftBottom"
	LeftTop     DrawingImageOrientation = "LeftTop"
	RightBottom DrawingImageOrientation = "RightBottom"
	RightTop    DrawingImageOrientation = "RightTop"
	TopLeft     DrawingImageOrientation = "TopLeft"
	TopRight    DrawingImageOrientation = "TopRight"
)

// LocationType
type LocationType string

const (
	FileSystem LocationType = "FileSystem"
	Virtual    LocationType = "Virtual"
)

// MetadataFields
type MetadataFields string

const (
	Cast                  MetadataFields = "Cast"
	ChannelNumber         MetadataFields = "ChannelNumber"
	Collections           MetadataFields = "Collections"
	CommunityRating       MetadataFields = "CommunityRating"
	CriticRating          MetadataFields = "CriticRating"
	Genres                MetadataFields = "Genres"
	Name                  MetadataFields = "Name"
	OfficialRating        MetadataFields = "OfficialRating"
	OriginalTitle         MetadataFields = "OriginalTitle"
	Overview              MetadataFields = "Overview"
	ProductionLocations   MetadataFields = "ProductionLocations"
	Runtime               MetadataFields = "Runtime"
	SortIndexNumber       MetadataFields = "SortIndexNumber"
	SortName              MetadataFields = "SortName"
	SortParentIndexNumber MetadataFields = "SortParentIndexNumber"
	Studios               MetadataFields = "Studios"
	Tagline               MetadataFields = "Tagline"
	Tags                  MetadataFields = "Tags"
)

// MediaProtocol
type MediaProtocol string

const (
	FTP  MediaProtocol = "Ftp"
	File MediaProtocol = "File"
	HTTP MediaProtocol = "Http"
	Mms  MediaProtocol = "Mms"
	RTP  MediaProtocol = "Rtp"
	RTSP MediaProtocol = "Rtsp"
	Rtmp MediaProtocol = "Rtmp"
	UDP  MediaProtocol = "Udp"
)

// SubtitleDeliveryMethod
type SubtitleDeliveryMethod string

const (
	Embed                               SubtitleDeliveryMethod = "Embed"
	Encode                              SubtitleDeliveryMethod = "Encode"
	External                            SubtitleDeliveryMethod = "External"
	HLS                                 SubtitleDeliveryMethod = "Hls"
	SubtitleDeliveryMethodVideoSideData SubtitleDeliveryMethod = "VideoSideData"
)

// ExtendedVideoSubTypes
type ExtendedVideoSubTypes string

const (
	DoviProfile02                      ExtendedVideoSubTypes = "DoviProfile02"
	DoviProfile10                      ExtendedVideoSubTypes = "DoviProfile10"
	DoviProfile22                      ExtendedVideoSubTypes = "DoviProfile22"
	DoviProfile30                      ExtendedVideoSubTypes = "DoviProfile30"
	DoviProfile42                      ExtendedVideoSubTypes = "DoviProfile42"
	DoviProfile50                      ExtendedVideoSubTypes = "DoviProfile50"
	DoviProfile61                      ExtendedVideoSubTypes = "DoviProfile61"
	DoviProfile76                      ExtendedVideoSubTypes = "DoviProfile76"
	DoviProfile81                      ExtendedVideoSubTypes = "DoviProfile81"
	DoviProfile82                      ExtendedVideoSubTypes = "DoviProfile82"
	DoviProfile83                      ExtendedVideoSubTypes = "DoviProfile83"
	DoviProfile84                      ExtendedVideoSubTypes = "DoviProfile84"
	DoviProfile85                      ExtendedVideoSubTypes = "DoviProfile85"
	DoviProfile92                      ExtendedVideoSubTypes = "DoviProfile92"
	ExtendedVideoSubTypesHdr10         ExtendedVideoSubTypes = "Hdr10"
	ExtendedVideoSubTypesHyperLogGamma ExtendedVideoSubTypes = "HyperLogGamma"
	ExtendedVideoSubTypesNone          ExtendedVideoSubTypes = "None"
	Hdr10Plus0                         ExtendedVideoSubTypes = "Hdr10Plus0"
)

// ExtendedVideoTypes
type ExtendedVideoTypes string

const (
	DolbyVision                     ExtendedVideoTypes = "DolbyVision"
	ExtendedVideoTypesHdr10         ExtendedVideoTypes = "Hdr10"
	ExtendedVideoTypesHyperLogGamma ExtendedVideoTypes = "HyperLogGamma"
	ExtendedVideoTypesNone          ExtendedVideoTypes = "None"
	Hdr10Plus                       ExtendedVideoTypes = "Hdr10Plus"
)

// SubtitleLocationType
type SubtitleLocationType string

const (
	InternalStream                    SubtitleLocationType = "InternalStream"
	SubtitleLocationTypeVideoSideData SubtitleLocationType = "VideoSideData"
)

// MediaStreamType
type MediaStreamType string

const (
	Attachment    MediaStreamType = "Attachment"
	Audio         MediaStreamType = "Audio"
	Data          MediaStreamType = "Data"
	EmbeddedImage MediaStreamType = "EmbeddedImage"
	Subtitle      MediaStreamType = "Subtitle"
	Unknown       MediaStreamType = "Unknown"
	Video         MediaStreamType = "Video"
)

// TransportStreamTimestamp
type TransportStreamTimestamp string

const (
	TransportStreamTimestampNone TransportStreamTimestamp = "None"
	Valid                        TransportStreamTimestamp = "Valid"
	Zero                         TransportStreamTimestamp = "Zero"
)

// MediaSourceType
type MediaSourceType string

const (
	Default     MediaSourceType = "Default"  // 普通视频和 AlistStrm 显示 Default
	Grouping    MediaSourceType = "Grouping" // HTTPStrm 显示 Grouping
	Placeholder MediaSourceType = "Placeholder"
)

// Video3DFormat
type Video3DFormat string

const (
	FullSideBySide   Video3DFormat = "FullSideBySide"
	FullTopAndBottom Video3DFormat = "FullTopAndBottom"
	HalfSideBySide   Video3DFormat = "HalfSideBySide"
	HalfTopAndBottom Video3DFormat = "HalfTopAndBottom"
	MVC              Video3DFormat = "MVC"
)

// PersonType
type PersonType string

const (
	Actor     PersonType = "Actor"
	Composer  PersonType = "Composer"
	Conductor PersonType = "Conductor"
	Director  PersonType = "Director"
	GuestStar PersonType = "GuestStar"
	Lyricist  PersonType = "Lyricist"
	Producer  PersonType = "Producer"
	Writer    PersonType = "Writer"
)

// SyncJobItemStatus
type SyncJobItemStatus string

const (
	Converting      SyncJobItemStatus = "Converting"
	Failed          SyncJobItemStatus = "Failed"
	Queued          SyncJobItemStatus = "Queued"
	ReadyToTransfer SyncJobItemStatus = "ReadyToTransfer"
	Synced          SyncJobItemStatus = "Synced"
	Transferring    SyncJobItemStatus = "Transferring"
)

// LiveTv.TimerType
type LiveTvTimerType string

const (
	DateTime LiveTvTimerType = "DateTime"
	Keyword  LiveTvTimerType = "Keyword"
	Program  LiveTvTimerType = "Program"
)
