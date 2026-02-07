package alist

import (
	"MediaWarp/internal/config"
	"MediaWarp/utils"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/allegro/bigcache/v3"
)

type alistToken struct {
	value    string       // 令牌 Token
	expireAt time.Time    // 令牌过期时间
	mutex    sync.RWMutex // 令牌锁
}
type AlistClient struct {
	endpoint string // 服务器入口 URL
	username string // 用户名
	password string // 密码

	userInfo UserInfoData

	token  alistToken
	client *http.Client
	cache  *bigcache.BigCache
}

// 获得AlistClient实例
func NewAlistClient(addr string, username string, password string, token *string) (*AlistClient, error) {
	client := AlistClient{
		endpoint: utils.GetEndpoint(addr),
		username: username,
		password: password,
		client:   utils.GetHTTPClient(),
	}
	if token != nil {
		client.token = alistToken{
			value:    *token,
			expireAt: time.Time{},
		}
	}

	if config.Cache.Enable && config.Cache.AlistAPITTL > 0 {
		cache, err := bigcache.New(context.Background(), bigcache.DefaultConfig(config.Cache.AlistAPITTL))
		if err == nil {
			client.cache = cache
		} else {
			return nil, fmt.Errorf("创建 Alist API 缓存失败: %w", err)
		}
	}

	userInfo, err := client.Me()
	if err != nil {
		return nil, fmt.Errorf("获取用户当前信息失败：%w", err)
	}
	client.userInfo = *userInfo

	return &client, nil
}

// 得到服务器入口
//
// 避免直接访问 endpoint 字段
func (client *AlistClient) GetEndpoint() string {
	return client.endpoint
}

// 得到用户名
//
// 避免直接访问 username 字段
func (client *AlistClient) GetUsername() string {
	return client.username
}

// 得到一个可用的 Token
//
// 先从缓存池中读取，若过期或者未找到则重新生成
func (client *AlistClient) getToken() (string, error) {
	var tokenDuration = 2*24*time.Hour - 5*time.Minute // Token 有效期为 2 天，提前 5 分钟刷新

	client.token.mutex.RLock()
	if client.token.value != "" && (client.token.expireAt.IsZero() || time.Now().Before(client.token.expireAt)) {
		// 零值表示永不过期
		defer client.token.mutex.RUnlock()
		return client.token.value, nil
	}

	loginData, err := client.authLogin() // 重新生成一个token
	client.token.mutex.RUnlock()
	if err != nil {
		return "", err
	}

	client.token.mutex.Lock()
	defer client.token.mutex.Unlock()
	client.token.value = loginData.Token
	client.token.expireAt = time.Now().Add(tokenDuration) // Token 有效期为30分钟

	return loginData.Token, nil
}

func doRequest[T any](client *AlistClient, r Request) (*T, error) {
	var resp AlistResponse[T]
	cacheKey := r.GetCacheKey()
	if cacheKey != "" && client.cache != nil {
		if data, err := client.cache.Get(cacheKey); err == nil {
			if json.Unmarshal(data, &resp) == nil {
				return &resp.Data, nil
			}
		}
	}

	req := newHTTPReq(client.GetEndpoint(), r)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	if r.NeedAuth() {
		token, err := client.getToken()
		if err != nil {
			return nil, err
		}
		req.Header.Add("Authorization", token)
	}

	res, err := client.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败: %w", err)
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, fmt.Errorf("解析响应体失败: %w", err)
	}

	if resp.Code != http.StatusOK {
		return nil, fmt.Errorf("请求失败，HTTP 状态码: %d, 响应状态码: %d, 响应信息: %s", res.StatusCode, resp.Code, resp.Message)
	}

	if cacheKey != "" && client.cache != nil {
		err = client.cache.Set(cacheKey, data)
		if err != nil {
			return nil, fmt.Errorf("缓存响应体失败: %w", err)
		}
	}

	return &resp.Data, nil
}

// ==========Alist API(v3) 相关操作==========

// 登录Alist（获取一个新的Token）
func (client *AlistClient) authLogin() (*AuthLoginData, error) {
	req := AuthLoginRequest{
		Username: client.GetUsername(),
		Password: client.password,
	}
	data, err := doRequest[AuthLoginData](client, &req)
	if err != nil {
		return nil, fmt.Errorf("登录失败: %w", err)
	}

	return data, nil
}

// 获取某个文件/目录信息
func (client *AlistClient) FsGet(req *FsGetRequest) (*FsGetData, error) {
	respData, err := doRequest[FsGetData](client, req)
	if err != nil {
		return nil, fmt.Errorf("获取文件/目录信息失败: %w", err)
	}
	return respData, nil
}

func (client *AlistClient) Me() (*UserInfoData, error) {
	data, err := doRequest[UserInfoData](client, &MeRequest{})
	if err != nil {
		return nil, fmt.Errorf("获取用户信息失败: %w", err)
	}
	return data, nil
}

// GetFileURL 获取文件的可访问 URL
func (client *AlistClient) GetFileURL(p string, isRawURL bool) (string, error) {
	fileData, err := client.FsGet(&FsGetRequest{Path: p, Page: 1})
	if err != nil {
		return "", fmt.Errorf("获取文件信息失败：%w", err)
	}
	if isRawURL {
		return fileData.RawURL, nil
	}
	var url strings.Builder
	url.WriteString(client.GetEndpoint())
	if fileData.Sign != "" {
		url.WriteString("?sign=" + fileData.Sign)
	}
	url.WriteString(path.Join("/d", client.userInfo.BasePath, p))
	return url.String(), nil
}

func (client *AlistClient) GetFsOther(req *FsOtherRequest, result any) error {
	respData, err := doRequest[any](client, req)
	if err != nil {
		return fmt.Errorf("请求失败: %w", err)
	}
	respBytes, err := json.Marshal(respData)
	if err != nil {
		return fmt.Errorf("序列化响应数据失败: %w", err)
	}
	err = json.Unmarshal(respBytes, &result)
	if err != nil {
		return fmt.Errorf("反序列化响应数据失败: %w", err)
	}
	return nil
}
