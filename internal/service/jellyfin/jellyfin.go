package jellyfin

import (
	"encoding/json"
	"fmt"
	"io"

	"net/http"
	"net/url"
	"strconv"

	"github.com/AkimioJR/MediaWarp/constants"
	"github.com/AkimioJR/MediaWarp/utils"
)

type Client struct {
	baseURL *url.URL
	apiKey  string // 认证方式：APIKey；获取方式：Jellyfin 控制台 -> 高级 -> API密钥
	client  *http.Client
}

// 获取媒体服务器类型
func (client *Client) GetType() constants.MediaServerType {
	return constants.JELLYFIN
}

// 获取 Jellyfin 连接地址
//
// 包含协议、服务器域名（IP）、端口号
// 示例：return "http://jellyfin.example.com:8096"
func (client *Client) GetBaseURLString() string {
	return client.baseURL.String()
}

// 获取 Jellyfin 的API Key
func (client *Client) GetAPIKey() string {
	return client.apiKey
}

// ItemsService
// /Items
func (client *Client) ItemsServiceQueryItem(ids string, limit int, fields string) (*Response, error) {
	var (
		params       = url.Values{}
		itemResponse = &Response{}
	)
	params.Add("Ids", ids)
	params.Add("Limit", strconv.Itoa(limit))
	params.Add("Fields", fields)

	api := client.baseURL.JoinPath("Items")
	api.RawQuery = params.Encode()

	req, err := http.NewRequest(http.MethodGet, api.String(), nil)
	if err != nil {
		return nil, err
	}
	// Jellyfin 12 移除了小写的 api_key 查询参数（返回 401），改用标准的 Authorization 头（兼容 10.11+）
	req.Header.Set("Authorization", fmt.Sprintf("MediaBrowser Token=\"%s\"", client.GetAPIKey()))

	resp, err := client.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("查询 Items 失败，HTTP 状态码：%d，响应体：%s", resp.StatusCode, string(body))
	}
	if len(body) == 0 {
		return nil, fmt.Errorf("查询 Items 返回空响应体")
	}

	if err = json.Unmarshal(body, itemResponse); err != nil {
		return nil, err
	}
	return itemResponse, nil
}

// 获取 Jellyfin 实例
func New(baseURL string, apiKey string) (*Client, error) {
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("解析 baseURL 失败: %w", err)
	}
	client := &Client{
		baseURL: parsedURL,
		apiKey:  apiKey,
		client:  utils.GetHTTPClient(),
	}
	return client, nil
}
