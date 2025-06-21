package main

import (
	"crontab/config"
	"crontab/middleware"
	"crontab/monitor"
	"time"
)

func main() {
	config.InitConfig()
	middleware.InitLogger(config.Conf.Log.Filepath)
	middleware.Logger.Println("进程守护程序启动，开始定时检查...")
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	monitor.Run()

	for {
		select {
		case <-ticker.C:
			monitor.Run()
		}
	}

}
