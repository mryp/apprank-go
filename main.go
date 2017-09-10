package main

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	startCrontab()
	startEchoServer()
}

func startCrontab() {
	//
}

func startEchoServer() {
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.CORS()) //CORS対応（他ドメインからAJAX通信可能にする）

	//ルーティング
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "apprank-go")
	})
	v1Group := e.Group("/v1")
	v1Group.GET("/now", NowHandler)
	v1Group.GET("/appinfo", AppInfoHandler)
	v1Group.GET("/apprank", AppRankHandler)

	//開始
	e.Logger.Fatal(e.Start(":" + strconv.Itoa(8080)))
}
