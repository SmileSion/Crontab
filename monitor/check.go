package monitor

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"crontab/config"
	"crontab/handler"
	"crontab/middleware"
)

type Program struct {
	Name      string
	StatusCmd string
	StartCmd  string
}

func runCommand(cmdStr string) (string, error) {
	parts := strings.Fields(cmdStr)
	if len(parts) == 0 {
		return "", fmt.Errorf("命令为空：%q", cmdStr)
	}
	cmd := exec.Command(parts[0], parts[1:]...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func checkAndRestart(p config.Program) {
	statusOutput, err := runCommand(p.StatusCmd)
	middleware.Logger.Printf("检查状态输出: %s", statusOutput)

	// 状态异常：准备重启并发短信
	if err != nil || strings.Contains(statusOutput, "not running") {
		// 第一次短信：发现未运行
		contentCheck := fmt.Sprintf("告警：您的程序 [%s] 未运行，正在尝试自动重启...", p.Name)
		success, msg, uid := handler.SendSmsWithContent(contentCheck)
		if success {
			middleware.Logger.Printf("短信发送成功 UID: %s，内容: %s\n返回消息: %s\n", uid, contentCheck, msg)
		} else {
			middleware.Logger.Printf("短信发送失败 UID: %s，内容: %s\n失败原因: %s\n", uid, contentCheck, msg)
		}

		middleware.Logger.Printf("%s 未运行，正在尝试重启...", p.Name)

		// 执行启动命令
		startOutput, startErr := runCommand(p.StartCmd)
		var contentResult string

		if startErr != nil {
			middleware.Logger.Printf("启动 %s 失败：%v\n输出：%s\n", p.Name, startErr, startOutput)
			contentResult = fmt.Sprintf("重启失败：程序 [%s] 启动失败。\n错误信息：%v", p.Name, startErr)
		} else {
			middleware.Logger.Printf("启动 %s 成功。输出：%s\n", p.Name, startOutput)
			contentResult = fmt.Sprintf("恢复通知：程序 [%s] 启动成功，系统已完成自动重启。", p.Name)
		}

		// 第二次短信：发送启动结果
		success, msg, uid = handler.SendSmsWithContent(contentResult)
		if success {
			middleware.Logger.Printf("短信发送成功 UID: %s，内容: %s\n返回消息: %s\n", uid, contentResult, msg)
		} else {
			middleware.Logger.Printf("短信发送失败 UID: %s，内容: %s\n失败原因: %s\n", uid, contentResult, msg)
		}
	} else {
		middleware.Logger.Printf("%s 正常运行。\n", p.Name)
	}
}

func Run() {
	for _, p := range config.Conf.Program {
		checkAndRestart(p)
	}
}
