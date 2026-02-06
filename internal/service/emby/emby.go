package emby

import (
	"MediaWarp/constants"
	"MediaWarp/utils"
	"encoding/json"
	"io"
	"net/url"
	"strconv"
)

type Client struct {
	endpoint string
	apiKey   string // 认证方式：APIKey；获取方式：Emby控制台 -> 高级 -> API密钥
}

// 获取媒体服务器类型
func (client *Client) GetType() constants.MediaServerType {
	return constants.EMBY
}

// 获取Emby连接地址
//
// 包含协议、服务器域名（IP）、端口号
// 示例：return "http://emby.example.com:8096"
func (client *Client) GetEndpoint() string {
	return client.endpoint
}

// 获取Emby的API Key
func (client *Client) GetAPIKey() string {
	return client.apiKey
}

// ItemsService
// /Items
func (client *Client) ItemsServiceQueryItem(ids string, limit int, fields string) (*EmbyResponse, error) {
	var (
		params       = url.Values{}
		itemResponse = &EmbyResponse{}
	)
	params.Add("Ids", ids)
	params.Add("Limit", strconv.Itoa(limit))
	params.Add("Fields", fields)
	params.Add("Recursive", "true")
	params.Add("api_key", client.GetAPIKey())
	api := client.GetEndpoint() + "/Items?" + params.Encode()
	resp, err := utils.GetHTTPClient().Get(api)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, itemResponse)
	if err != nil {
		return nil, err
	}
	return itemResponse, nil
}

// 获取index.html内容 API：/web/index.html
func (client *Client) GetIndexHtml() ([]byte, error) {
	resp, err := utils.GetHTTPClient().Get(client.GetEndpoint() + "/web/index.html")
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	htmlContent, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return htmlContent, nil
}

// 获取Emby实例
func New(addr string, apiKey string) *Client {
	client := &Client{
		endpoint: utils.GetEndpoint(addr),
		apiKey:   apiKey,
	}
	return client
}
