package config

import (
	"ip-tunnel-proxy/pkg/config/po"
	"sync"
)

var mutex sync.Mutex

func Fill(rp po.Rp) {
	mutex.Lock()
	defer mutex.Unlock()

	Port = rp.Port
	NoAuth = rp.NoAuth
	HeartbeatInterval = rp.HeartbeatInterval
	UserConnectLimit = rp.ConfigData.UserConnectLimit
	AuthListWithIpMap = rp.ConfigData.AuthListWithIpMap
	AuthListWithUNameMap = rp.ConfigData.AuthListWithUNameMap
	AppId = rp.AppId
	PoolId = rp.PoolId
}
