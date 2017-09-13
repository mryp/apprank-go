package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/mryp/apprank-go/config"
	"github.com/robfig/cron"
)

func main() {
	if !config.Init() {
		log.Println("設定ファイル読み込み失敗（デフォルト値動作）")
	}
	startCrontab()
	startEchoServer()
}

//startCrontab はCRONの設定と初回起動を行う
func startCrontab() {
	interval := config.Now().Watch.Interval
	rankWatcher = NewRankWatcher()
	if interval > 0 {
		c := cron.New()
		format := fmt.Sprintf("0 */%d * * * *", interval) //分単位指定
		c.AddFunc(format, func() {
			rankWatcher.StartBgTask()
		})
		c.Start()
	}

	rankWatcher.StartBgTask()
}

//startEchoServer はHTTPサーバーの初期化と起動を行う
func startEchoServer() {
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.CORS()) //CORS対応（他ドメインからAJAX通信可能にする）
	switch config.Now().Log.Output {
	case "stream":
		e.Use(middleware.Logger())
	case "file":
		//未実装
	}

	//ルーティング
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "apprank-go")
	})
	v1Group := e.Group("/v1")
	v1Group.GET("/now", NowHandler)
	v1Group.GET("/appinfo", AppInfoHandler)
	v1Group.GET("/apprank", AppRankHandler)

	//開始
	e.Logger.Fatal(e.Start(":" + strconv.Itoa(config.Now().Server.PortNumber)))
}
