package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo"
)

type AppInfoRequest struct {
	ID int64 `json:"id" xml:"id" form:"id" query:"id"`
}

type AppInfoResponse struct {
	Name       string `json:"name" xml:"name"`
	InfoURL    string `json:"info_url" xml:"info_url"`
	ArtworkURL string `json:"artwork_url" xml:"artwork_url"`
	ArtistName string `json:"artist_name" xml:"artist_name"`
	ArtistURL  string `json:"artist_url" xml:"artist_url"`
	Copyright  string `json:"copyright" xml:"copyright"`
}

type AppRankRequest struct {
	ID      int64     `json:"id" xml:"id" form:"id" query:"id"`
	Country string    `json:"country" xml:"country" form:"country" query:"country"`
	Kind    int       `json:"kind" xml:"kind" form:"kind" query:"kind"`
	Start   time.Time `json:"start" xml:"start" form:"start" query:"start"`
	End     time.Time `json:"end" xml:"end" form:"end" query:"end"`
}

type AppRankResponse struct {
	Apps []AppRankAppsResponse `json:"apps" xml:"apps"`
}

type AppRankAppsResponse struct {
	Rank    int       `json:"rank" xml:"rank" form:"rank" query:"rank"`
	Updated time.Time `json:"updated" xml:"updated" form:"updated" query:"updated"`
}

func AppInfoHandler(c echo.Context) error {
	req := new(AppInfoRequest)
	if err := c.Bind(req); err != nil {
		return err
	}
	fmt.Printf("AppInfoHandler request=%v\n", *req)

	//とりあえずダミーをセット
	response := new(AppInfoResponse)
	response.Name = "アイドルマスター シンデレラガールズ スターライトステージ"
	response.InfoURL = "https://itunes.apple.com/jp/app/%E3%82%A2%E3%82%A4%E3%83%89%E3%83%AB%E3%83%9E%E3%82%B9%E3%82%BF%E3%83%BC-%E3%82%B7%E3%83%B3%E3%83%87%E3%83%AC%E3%83%A9%E3%82%AC%E3%83%BC%E3%83%AB%E3%82%BA-%E3%82%B9%E3%82%BF%E3%83%BC%E3%83%A9%E3%82%A4%E3%83%88%E3%82%B9%E3%83%86%E3%83%BC%E3%82%B8/id1016318735?mt=8&app=itunes"
	response.ArtistName = "BANDAI NAMCO Entertainment Inc."
	response.ArtistURL = "https://itunes.apple.com/jp/developer/bandai-namco-entertainment-inc/id352305770?mt=8"
	response.Copyright = "©2015 BANDAI NAMCO Entertainment Inc."

	return c.JSON(http.StatusOK, response)
}

func AppRankHandler(c echo.Context) error {
	req := new(AppRankRequest)
	if err := c.Bind(req); err != nil {
		return err
	}
	fmt.Printf("AppRankHandler request=%v\n", *req)

	//とりあえずダミーをセット
	response := new(AppRankResponse)
	apps := make([]AppRankAppsResponse, 0)

	updated1, _ := time.Parse("2006-01-02 15:04:05", "2017-09-01 10:00:00")
	apps = append(apps, AppRankAppsResponse{
		Rank:    1,
		Updated: updated1,
	})
	response.Apps = apps

	return c.JSON(http.StatusOK, response)
}
