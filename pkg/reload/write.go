package reload

import (
	"encoding/csv"
	"ip-tunnel-proxy/pkg/config"
	"ip-tunnel-proxy/pkg/utils"
	"os"
	"time"
)

// 热加载自动配置文件
func InputReloadFile() {
	for {
		if config.ProxyRun {
			WriteReloadCsv()
			return
		}

		time.Sleep(time.Second * time.Duration(5))
	}
}

// 写入自动配置文件
func WriteReloadCsv() {

	exist := utils.FileExist(reloadName)

	if exist {
		os.Remove(reloadName)
	}

	//创建文件
	f, err := os.Create(reloadName)
	if err != nil {
		return
	}

	defer f.Close()

	w := csv.NewWriter(f)
	data := [][]string{
		{config.Number},
	}

	w.WriteAll(data)
	w.Flush()
}
