package monitor

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"crontab/config"
	// "crontab/handler"
	"crontab/middleware"
)

// 判断程序是否运行，只匹配程序名
func isRunning(p config.Program) bool {
    pid := os.Getpid() // 守护程序自身 PID

    // 只匹配程序名，不匹配路径
    cmdStr := fmt.Sprintf(
        "ps -eo pid,cmd | grep '%s' | grep -v '^ *%d ' | grep -v 'bash -c' | grep -v 'grep'",
        p.Name, pid,
    )

    output, _ := runCommandOutput(cmdStr)
	middleware.Logger.Printf("[%s] isRunning 输出: %q", p.Name, output)
    return strings.TrimSpace(output) != ""
}

// 执行命令并返回输出，用于 isRunning
func runCommandOutput(cmdStr string) (string, error) {
    middleware.Logger.Printf("[Debug] 执行命令: %s", cmdStr)
    cmd := exec.Command("bash", "-c", cmdStr)
    out, err := cmd.Output()
    return string(out), err
}

// 执行命令，不收集 stdout/stderr，只返回 error
func runCommand(cmdStr string) error {
    middleware.Logger.Printf("[Debug] 执行命令: %s", cmdStr)
    cmd := exec.Command("bash", "-c", cmdStr)
    return cmd.Run()
}


// 检查并自动重启逻辑
func checkAndRestart(p config.Program) {
	defer func() {
		if r := recover(); r != nil {
			middleware.Logger.Printf("[%s] 异常恢复: %v", p.Name, r)
		}
	}()

	if isRunning(p) {
		middleware.Logger.Printf("[%s] 正常运行。", p.Name)
		return
	}

	// 程序未运行，发送告警短信
	// contentCheck := fmt.Sprintf("告警：程序 [%s] 未运行，尝试自动重启...", p.Name)
	// success, msg, uid := handler.SendSmsWithContent(contentCheck)
	// if success {
	// 	middleware.Logger.Printf("短信发送成功 UID: %s, 内容: %s", uid, contentCheck)
	// } else {
	// 	middleware.Logger.Printf("短信发送失败 UID: %s, 内容: %s, 失败原因: %s", uid, contentCheck, msg)
	// }

	middleware.Logger.Printf("[%s] 未运行，开始重启...", p.Name)

	// 使用完整路径，支持目录+程序名配置
	programPath := filepath.Join(p.Path, p.Name)
	dir := filepath.Dir(programPath)
	file := filepath.Base(programPath)
	logPath := filepath.Join(dir, "run.log")

	// nohup 后台启动，追加日志
	startCmd := fmt.Sprintf("cd %s && nohup ./%s >> %s 2>&1 &", dir, file, logPath)

	startErr := runCommand(startCmd)
	// var contentResult string
	if startErr != nil {
		middleware.Logger.Printf("[%s] 启动失败: %v", p.Name, startErr)
		// contentResult = fmt.Sprintf("重启失败：程序 [%s] 启动失败。\n错误信息：%v", p.Name, startErr)
	} else {
		middleware.Logger.Printf("[%s] 启动成功", p.Name)
		// contentResult = fmt.Sprintf("恢复通知：程序 [%s] 启动成功，系统已完成自动重启。", p.Name)
	}

	// 发送启动结果短信
	// success, msg, uid = handler.SendSmsWithContent(contentResult)
	// if success {
	// 	middleware.Logger.Printf("短信发送成功 UID: %s, 内容: %s", uid, contentResult)
	// } else {
	// 	middleware.Logger.Printf("短信发送失败 UID: %s, 内容: %s, 失败原因: %s", uid, contentResult, msg)
	// }
}

// 扫描所有配置程序，并发检查
func Run() {
	var wg sync.WaitGroup
	for _, p := range config.Conf.Program {
		wg.Add(1)
		go func(pr config.Program) {
			defer wg.Done()
			checkAndRestart(pr)
		}(p)
	}
	wg.Wait()
}
