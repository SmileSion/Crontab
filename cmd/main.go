package main

import (
	"crontab/config"
	"crontab/middleware"
	"crontab/monitor"
	"time"
)

func main() {
	// 初始化配置和日志
	config.InitConfig()
	middleware.InitLogger(config.Conf.Log.Filepath)

	middleware.Logger.Println("进程守护程序启动，开始定时检查...")

	// 先执行一次检查
	safeRunMonitor()

	// 定时器，每分钟执行一次
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			safeRunMonitor()
		}
	}
}

// 包装 monitor.Run()，捕获异常，保证守护程序不中断
func safeRunMonitor() {
	defer func() {
		if r := recover(); r != nil {
			middleware.Logger.Printf("监控任务异常恢复: %v\n", r)
		}
	}()
	monitor.Run()
}
