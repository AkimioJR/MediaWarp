package alist

type AliyunVideoData struct {
	DriveId              string `json:"drive_id"`
	FileId               string `json:"file_id"`
	VideoPreviewPlayInfo struct {
		Category                string `json:"category"`
		LiveTranscodingTaskList []struct {
			Stage          string `json:"stage"`
			Status         string `json:"status"`
			TemplateHeight int    `json:"template_height"`
			TemplateId     string `json:"template_id"`
			TemplateName   string `json:"template_name"`
			TemplateWidth  int    `json:"template_width"`
			Url            string `json:"url"` // 视频播放 URL
		} `json:"live_transcoding_task_list"`
		Meta struct {
			Duration float64 `json:"duration"`
			Height   int     `json:"height"`
			Width    int     `json:"width"`
		} `json:"meta"`
	} `json:"video_preview_play_info"`
}
