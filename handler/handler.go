package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"crontab/config"
	"crontab/middleware"
	"crontab/utils"
)

// 短信请求结构体
type SmsRequest struct {
	SecretKey string `json:"secret_key"`
	Uid       string `json:"uid"`
	AppKey    string `json:"app_key"`
	Phone     string `json:"phone"`
	Sign      string `json:"sign"`
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
	UserId    string `json:"userId"`
}

// 短信响应结构体
type SmsResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}


// 只传content的短信发送函数
func SendSmsWithContent(content string) (success bool, message string, uid string) {
	fixedUserId := config.Conf.Msg.UserId
	fixedPhone  := config.Conf.Msg.Phone
	content = strings.TrimSpace(content)
	uid = utils.GenerateUID()
	timestamp := time.Now().Unix()

	params := map[string]string{
		"app_key":    config.Conf.Msg.AppKey,
		"content":    content,
		"phone":      fixedPhone,
		"secret_key": config.Conf.Msg.SecretKey,
		"timestamp":  fmt.Sprintf("%d", timestamp),
		"uid":        uid,
		"userId":     fixedUserId,
		"sign":       "",
	}
	params["sign"] = utils.GenerateSM3Sign(params)

	requestBody := SmsRequest{
		SecretKey: params["secret_key"],
		Uid:       params["uid"],
		AppKey:    params["app_key"],
		Phone:     params["phone"],
		Sign:      params["sign"],
		Content:   params["content"],
		Timestamp: params["timestamp"],
		UserId:    params["userId"],
	}

	jsonData, _ := json.Marshal(requestBody)
	resp, err := http.Post(config.Conf.Msg.Addr, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		middleware.Logger.Printf("短信发送失败: %v", err)
		return false, "短信接口调用失败", uid
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var smsResp SmsResponse
	_ = json.Unmarshal(body, &smsResp)

	middleware.Logger.Printf("[短信日志] UID: %s | UserId: %s | Phone: %s | Content: %s | 状态: %s | 返回: %s",
		uid,
		fixedUserId,
		fixedPhone,
		content,
		func() string {
			if smsResp.Success {
				return "发送成功"
			}
			return "发送失败"
		}(),
		smsResp.Message,
	)

	return smsResp.Success, smsResp.Message, uid
}
