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

// 运行命令并返回输出
func runCommand(cmdStr string) (string, error) {
	middleware.Logger.Printf("[Debug] 实际执行命令: %s", cmdStr)
	cmd := exec.Command("bash", "-c", cmdStr)  // 通过bash -c执行整条命令字符串
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

// 判断输出中是否包含任意一个关键字
func containsAny(s string, keywords []string) bool {
	for _, kw := range keywords {
		if strings.Contains(s, kw) {
			return true
		}
	}
	return false
}

// 检查并自动重启逻辑
func checkAndRestart(p config.Program) {
	statusOutput, err := runCommand(p.StatusCmd)

	// 清洗并小写化输出
	cleanOutput := strings.ToLower(strings.TrimSpace(statusOutput))

	// 输出调试信息
	middleware.Logger.Printf("[%s] 状态命令原始输出: %q", p.Name, statusOutput)
	middleware.Logger.Printf("[%s] 状态命令清洗后输出: %q", p.Name, cleanOutput)

	// 定义异常关键词
	keywords := []string{"not running", "inactive", "failed", "stopped"}

	// 判断是否异常
	isNotRunning := err != nil || containsAny(cleanOutput, keywords)

	if isNotRunning {
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
		middleware.Logger.Printf("[%s] 正常运行。\n", p.Name)
	}
}

// 扫描所有配置程序并检查
func Run() {
	for _, p := range config.Conf.Program {
		checkAndRestart(p)
	}
}
