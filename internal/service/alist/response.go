package alist

import (
	"encoding/json"
)

type AlistResponse[T any] struct {
	Code    int64  `json:"code"`    // 状态码
	Data    T      `json:"data"`    // data
	Message string `json:"message"` // 信息
}

type AuthLoginData struct {
	Token string `json:"token"` // token
}

type FsGetData struct {
	Created  string `json:"created"` // 创建时间
	HashInfo any    `json:"hash_info"`
	Hashinfo string `json:"hashinfo"`
	Header   string `json:"header"`
	IsDir    bool   `json:"is_dir"`   // 是否是文件夹
	Modified string `json:"modified"` // 修改时间
	Name     string `json:"name"`     // 文件名
	Provider string `json:"provider"`
	RawURL   string `json:"raw_url"` // 原始url
	Readme   string `json:"readme"`  // 说明
	Related  any    `json:"related"`
	Sign     string `json:"sign"`  // 签名
	Size     int64  `json:"size"`  // 大小
	Thumb    string `json:"thumb"` // 缩略图
	Type     int64  `json:"type"`  // 类型
}

type UserInfoData struct {
	BasePath   string     `json:"base_path"`  // 根目录
	Disabled   bool       `json:"disabled"`   // 是否禁用
	ID         int64      `json:"id"`         // id
	Otp        bool       `json:"otp"`        // 是否开启二步验证
	Password   string     `json:"password"`   // 密码
	Permission int64      `json:"permission"` // 权限
	Role       Int64Slice `json:"role"`       // 角色
	SsoID      string     `json:"sso_id"`     // sso id
	Username   string     `json:"username"`   // 用户名
}

type Int64Slice []int64

func (s *Int64Slice) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		*s = []int64{}
		return nil
	}
	var single int64
	if err := json.Unmarshal(data, &single); err == nil {
		*s = []int64{single}
		return nil
	}
	var arr []int64
	if err := json.Unmarshal(data, &arr); err != nil {
		return err
	}
	*s = arr
	return nil
}
