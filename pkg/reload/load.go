package reload

import (
	"ip-tunnel-proxy/pkg/config"
	"ip-tunnel-proxy/pkg/utils"
	"log"
)

var reloadName = config.OsDesktop + "reload.csv"

// 规定隧道文件格式如下，csv格式，逗号分割
// 第一列：设备号
func FindReloadFileAndLoad() {

	defer func() {
		if p := recover(); p != nil {
			log.Printf("load#FindReloadFileAndLoad internal error: %v", p)
		}
	}()

	exist := utils.FileExist(reloadName)

	if exist {
		data := ReadOneLineWithCSV(reloadName)

		if data == nil || len(data) != 1 {
			return
		}

		config.Number = data[0]

		// 首次启动
		config.CommandChan <- config.StartCommand
	}
}
