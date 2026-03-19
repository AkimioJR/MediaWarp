package service

import (
	"MediaWarp/internal/config"
	"MediaWarp/internal/logging"
	"MediaWarp/internal/service/alist"
	"MediaWarp/utils"
	"fmt"
	"sync"
)

var (
	alistClientMap sync.Map
)

// 初始化 Alist 客户端
func InitAlistClient() {
	if config.AlistStrm.Enable {
		for _, alist := range config.AlistStrm.List {
			registerAlistClient(alist.ADDR, alist.Username, alist.Password, alist.Token)
		}
	}
}

// 注册Alist客户端
//
// 将Alist客户端注册到全局Map中
func registerAlistClient(addr string, username string, password string, token *string) {
	c, err := alist.New(addr, username, password, token)
	if err != nil {
		logging.Warningf("注册 Alist 客户端 %s 失败：%s", addr, err)
		return
	}
	alistClientMap.Store(c.GetEndpoint(), c)
}

// 获取Alist客户端
//
// 从全局Map中获取Alist客户端
func GetAlistClient(addr string) (*alist.Client, error) {
	endpoint := utils.GetEndpoint(addr)
	if client, ok := alistClientMap.Load(endpoint); ok {
		return client.(*alist.Client), nil
	}
	return nil, fmt.Errorf("%s 未注册到 Alist 客户端列表中", endpoint)
}
