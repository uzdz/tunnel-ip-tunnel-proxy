package config

import (
	"ip-tunnel-proxy/pkg/utils"
	"sync"
	"time"
)

// 代理是否启动
var ProxyRun bool

var DownloadNow = false

var StartMutex sync.Mutex

var ShutdownMutex sync.Mutex

var OsDesktop = "C:\\Users\\Administrator\\Desktop\\"

var LogPath = OsDesktop + "sys.log"

var ServerHost = "www.xxx.com:25435"

// 数据上报服务端地址
var RemoteUrl = "http://" + ServerHost + "/xxx-boot/ip/report/tunnel"

// 当前服务器时间
var TimeUrl = "http://" + ServerHost + "/xxx-boot/ip/time"

// 选择请求转发代理IP地址
var SwitchRemoteUrl = "http://" + ServerHost + "/xxx-boot/ip/tunnel/switch"

// 检查新版本服务端地址
var FileRequestUrl = "http://" + ServerHost + "/xxx-boot/ip/client/version"

// 是否切换IP
var ProxySwitchIp = "Proxy-Switch-Ip"

// 代理服务器授权Key
var ProxyAuthorizationKey = "Proxy-Authorization"

// 代理隧道随机转发授权KEY
var XTunnelForwardedFor = "X-Tunnel-Forwarded-For"

// 隧道代理服务需添加如下Header
var TunnelXForKey = "X-Tunnel-Key"
var TunnelXForValue = "266aXa2WWe"

// 命令队列
var CommandChan = make(chan string, 10)

// 机器编号
var Number string

var IgnoreHeaderMap = []string{
	ProxySwitchIp,
	XTunnelForwardedFor,
	ProxyAuthorizationKey,
	TunnelXForKey,
}

var TunnelUserExp = utils.NewConcurrentMap()

// 启动代理服务器
var StartCommand = "start"

// 关闭代理服务器
var ShutdownCommand = "shutdown"

// 定义一个任务触发器
var DialTicker = time.NewTicker(time.Second * time.Duration(9999))
