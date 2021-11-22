package ticker

import (
	"ip-tunnel-proxy/pkg/config"
	"ip-tunnel-proxy/pkg/out"
	"ip-tunnel-proxy/pkg/utils"
	"log"
	"os/exec"
	"runtime"
	"runtime/debug"
	"strconv"
	"time"
)

func CleanYesterdayLog() {
	// 初始化
	out.LogFileOut(config.LogPath)
	for {
		now := time.Now()
		next := now.Add(time.Hour * 24)
		next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
		t := time.NewTimer(next.Sub(now))

		<-t.C
		out.LogFileOut(config.LogPath)
	}
}

func TimingDialTask() {
	for range config.DialTicker.C {
		config.CommandChan <- config.StartCommand
	}
}

func MemoryLimit() {

	var maxUseHeap uint64 = 40960

	for {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		use := m.Alloc / 1024

		if use >= maxUseHeap {
			log.Printf("内存超过最大使用限制：%d Kb 最大 %d Kb\n", use, maxUseHeap)
			debug.FreeOSMemory()
		}

		time.Sleep(time.Second * 30)
	}
}

func TimeResetTicker() {
	for {
		now := time.Now()
		next := now.Add(time.Hour * 24)
		next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
		t := time.NewTimer(next.Sub(now))

		<-t.C

		TimeLoad()
	}
}

func TimeLoad() {
	_, value := utils.SampleGet(config.TimeUrl)
	timer, err := strconv.ParseInt(value, 10, 64)
	if err == nil {
		// 进行时间设置
		TimeReset(timer)
	}
}

func TimeReset(sec int64) {
	// 传入当前unix秒级时间戳，重制当前系统时间
	time := time.Unix(sec, 0)

	hms := time.Format("15:04:05")
	cmd := exec.Command("cmd", "/c", "time", hms)
	utils.Output(cmd, "GB18030")

	date := time.Format("2006-01-02")
	cmd = exec.Command("cmd", "/c", "date", date)
	utils.Output(cmd, "GB18030")
}
