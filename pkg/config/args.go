package config

// -------------- 标准参数

// 隧道心跳时间
var HeartbeatInterval int64

// 是否需要验证 0需要 1不需要
var NoAuth int64

// 端口号
var Port int64

// 用户访问限制
var UserConnectLimit map[string]string

// 授权用户列表（IP反转MAP）
var AuthListWithIpMap map[string]string

// 授权用户列表（用户名反转MAP）
var AuthListWithUNameMap map[string]string

// 应用Id
var AppId string

// 池子Id
var PoolId string
