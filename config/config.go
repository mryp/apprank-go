package config

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

const settingFileName = "config.toml"

//EnvConfig 環境設定構造体
type EnvConfig struct {
	Log    LogEnvConfig
	Server ServerEnvConfig
	DB     DBEnvConfig
	Watch  WatchEnvConfig
}

//LogEnvConfig ログ設定情報
type LogEnvConfig struct {
	Output string
}

//ServerEnvConfig HTTPサーバー設定情報
type ServerEnvConfig struct {
	PortNumber int
}

//DBEnvConfig DB接続設定情報
type DBEnvConfig struct {
	UserID     string
	Password   string
	HostName   string
	PortNumber string
	Name       string
}

//WatchEnvConfig 監視設定情報
type WatchEnvConfig struct {
	Interval int
	MaxCount int
}

// 設定情報保持変数
var envConfig = EnvConfig{
	Log:    LogEnvConfig{Output: "stream"},
	Server: ServerEnvConfig{PortNumber: 8080},
	DB:     DBEnvConfig{UserID: "root", Password: "root", HostName: "127.0.0.1", PortNumber: "3306", Name: "apprank"},
	Watch:  WatchEnvConfig{Interval: 60, MaxCount: 200},
}

//Init 設定ファイルから設定を読み込み（失敗時はfalseを返す）
func Init() bool {
	var config EnvConfig
	if _, err := toml.DecodeFile(settingFileName, &config); err != nil {
		fmt.Println(err)
		return false
	}

	fmt.Printf("config=%#v\n", config)
	envConfig = config
	return true
}

//Env 現在の設定値を取得する
func Now() EnvConfig {
	return envConfig
}
