package alist

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type Request interface {
	GetMethod() string
	GetAPIPath() string
	NeedAuth() bool
	GetCacheKey() string
}

func getReqBody(r Request) io.Reader {
	b, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	return bytes.NewBuffer(b)
}

func newHTTPReq(endpoint string, r Request) *http.Request {
	var (
		req *http.Request
		err error
	)
	if r.GetMethod() == http.MethodGet {
		req, err = http.NewRequest(r.GetMethod(), endpoint+r.GetAPIPath(), nil)
	} else {
		req, err = http.NewRequest(r.GetMethod(), endpoint+r.GetAPIPath(), getReqBody(r))
	}
	if err != nil {
		panic(fmt.Errorf("创建请求失败: %w", err))
	}
	return req
}

type FsGetRequest struct {
	Path     string `json:"path"`
	Password string `json:"password"`
	Page     uint32 `json:"page"`
	PerPage  uint32 `json:"per_page"`
	Refresh  bool   `json:"refresh"`
}

func (FsGetRequest) GetMethod() string {
	return http.MethodPost
}

func (FsGetRequest) GetAPIPath() string {
	return "/api/fs/get"
}

func (FsGetRequest) NeedAuth() bool {
	return true
}

func (req *FsGetRequest) GetCacheKey() string {
	return req.GetAPIPath() + req.Path + req.Password + strconv.Itoa(int(req.Page)) + strconv.Itoa(int(req.PerPage)) + strconv.FormatBool(req.Refresh)
}

type FsOtherRequest struct {
	Path     string `json:"path"`
	Method   string `json:"method"`
	Password string `json:"password"`
}

func (FsOtherRequest) GetMethod() string {
	return http.MethodPost
}

func (FsOtherRequest) GetAPIPath() string {
	return "/api/fs/other"
}

func (FsOtherRequest) NeedAuth() bool {
	return true
}

func (req *FsOtherRequest) GetCacheKey() string {
	return req.GetAPIPath() + req.Path + req.Method + req.Password
}

type AuthLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (AuthLoginRequest) GetMethod() string {
	return http.MethodPost
}

func (AuthLoginRequest) GetAPIPath() string {
	return "/api/auth/login"
}

func (AuthLoginRequest) NeedAuth() bool {
	return false
}

func (req *AuthLoginRequest) GetCacheKey() string {
	return ""
}

type MeRequest struct{}

func (MeRequest) GetMethod() string {
	return http.MethodGet
}

func (MeRequest) GetAPIPath() string {
	return "/api/me"
}

func (MeRequest) NeedAuth() bool {
	return true
}

func (req MeRequest) GetCacheKey() string {
	return req.GetAPIPath()
}
