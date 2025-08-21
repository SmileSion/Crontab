package config

import (
	"log"

	"github.com/BurntSushi/toml"
	"crontab/utils"
)

type LogConfig struct {
	Filepath   string `toml:"filepath"`
	MaxSize    int    `toml:"max_size"`
	MaxBackups int    `toml:"max_backups"`
	MaxAge     int    `toml:"max_age"`
	Compress   bool   `toml:"compress"`
}

type MsgConfig struct {
	Addr      string `toml:"addr"`
	AppKey    string `toml:"app_key"`
	SecretKey string `toml:"secret_key"`
	Phone     string `toml:"phone"`
	UserId    string `toml:"userid"`
}

// Program 配置只保留 Name 和 Path
type Program struct {
	Name string `toml:"name"`  // 程序名称，用于 pgrep 查找和启动
	Path string `toml:"path"`  // 程序所在路径
}

type Config struct {
	Log     LogConfig   `toml:"log"`
	Msg     MsgConfig   `toml:"msgconfig"`
	Program []Program   `toml:"program"`
}

var Conf Config

func InitConfig() {
	if _, err := toml.DecodeFile("config/config.toml", &Conf); err != nil {
		panic(err)
	}

	log.Println("配置文件加载成功")
	log.Printf("配置文件加载成功，共读取 %d 个程序配置项\n", len(Conf.Program))

	// 解密敏感字段
	if decrypted, err := utils.Decrypt(Conf.Msg.AppKey); err == nil {
		Conf.Msg.AppKey = decrypted
	} else {
		log.Fatalf("AppKey解密失败: %v", err)
	}

	if decrypted, err := utils.Decrypt(Conf.Msg.SecretKey); err == nil {
		Conf.Msg.SecretKey = decrypted
	} else {
		log.Fatalf("SecretKey解密失败: %v", err)
	}
}
