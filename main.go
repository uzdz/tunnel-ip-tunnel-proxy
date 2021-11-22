package main

import (
	"ip-tunnel-proxy/pkg"
	"ip-tunnel-proxy/pkg/config"
	"ip-tunnel-proxy/pkg/gui"
	"ip-tunnel-proxy/pkg/reload"
	"ip-tunnel-proxy/pkg/ticker"
)

func main() {

	// 系统时间设置
	ticker.TimeLoad()
	go ticker.TimeResetTicker()

	// 设置开机自启动脚本
	go reload.AutoStart()

	// 先关闭拨号任务
	config.DialTicker.Stop()
	// 开启拨号任务监听
	go ticker.TimingDialTask()

	// 记录日志文件，并且每天24点定时清空文件
	go ticker.CleanYesterdayLog()

	// 监控堆使用大小
	go ticker.MemoryLimit()

	// 写入更新VPS脚本（覆盖）
	reload.UpdateVbsInit()

	// 定时检查是否有新版本更新
	go reload.ReportVAndDownload()

	// 写入热启动配置文件
	go reload.InputReloadFile()

	// 任务执行池
	go pkg.CommandServerThread()

	// 查找文件，存在则自启动
	reload.FindReloadFileAndLoad()

	// GUI窗口启动
	gui.Start()
}
