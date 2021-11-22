package pkg

import (
	"ip-tunnel-proxy/pkg/config"
	"ip-tunnel-proxy/pkg/config/po"
	"ip-tunnel-proxy/pkg/report"
	"ip-tunnel-proxy/pkg/server"
	"log"
	"strconv"
	"time"
)

func run() {
	port := ":" + strconv.Itoa(int(config.Port))
	server.P = &server.S{}
	go server.P.OpenServer(port)
}

func start() {
	config.StartMutex.Lock()
	defer config.StartMutex.Unlock()

	rq := po.Rq{}
	rq.DeviceId = config.Number
	err, reLoad := report.Report(rq)
	if err != nil {
		log.Printf("WARING 代理服务器上报失败！\n" + err.Error())

		time.Sleep(time.Second * time.Duration(5))
		config.CommandChan <- config.StartCommand
		return
	}

	if server.P == nil {
		run()
	} else if reLoad {
		shutdown()
		run()
	}

	// 拨号任务重新计时
	var waitTime int64
	if config.HeartbeatInterval == 0 {
		waitTime = 30
	} else {
		waitTime = config.HeartbeatInterval
	}
	config.DialTicker.Reset(time.Second * time.Duration(waitTime))

	config.ProxyRun = true
}

func shutdown() {
	config.ShutdownMutex.Lock()
	defer config.ShutdownMutex.Unlock()

	if server.P != nil {
		b := server.P.Shutdown()

		if b {
			server.P = nil
		}
	}

	config.DialTicker.Stop()
	config.ProxyRun = false
	log.Println("代理服务器已成功关闭...")
}

func CommandServerThread() {

	for {
		command := <-config.CommandChan

		for {
			// 如果下载状态下，阻塞等待下载完成
			if config.DownloadNow {
				time.Sleep(time.Second * time.Duration(5))
			} else {
				break
			}
		}

		switch command {
		case config.StartCommand:
			start()
		case config.ShutdownCommand:
			shutdown()
		}
	}
}
