package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/mryp/apprank-go/db"
)

//AppInfoRequest はアプリ情報RESTのリクエスト
type AppInfoRequest struct {
	ID int64 `json:"id" xml:"id" form:"id" query:"id"`
}

//AppInfoResponse はアプリ情報RESTのレスポンス
type AppInfoResponse struct {
	Name       string `json:"name" xml:"name"`
	InfoURL    string `json:"info_url" xml:"info_url"`
	ArtworkURL string `json:"artwork_url" xml:"artwork_url"`
	ArtistName string `json:"artist_name" xml:"artist_name"`
	ArtistURL  string `json:"artist_url" xml:"artist_url"`
	Copyright  string `json:"copyright" xml:"copyright"`
}

//AppRankRequest はアプリ順位RESTのリクエスト
type AppRankRequest struct {
	ID      int64  `json:"id" xml:"id" form:"id" query:"id"`
	Country string `json:"country" xml:"country" form:"country" query:"country"`
	Kind    int    `json:"kind" xml:"kind" form:"kind" query:"kind"`
	Start   string `json:"start" xml:"start" form:"start" query:"start"`
	End     string `json:"end" xml:"end" form:"end" query:"end"`
}

//AppRankResponse はアプリ順位RESTのレスポンス
type AppRankResponse struct {
	Apps []AppRankAppsResponse `json:"apps" xml:"apps"`
}

//AppRankAppsResponse はアプリ順位RESTのランキング詳細レスポンス
type AppRankAppsResponse struct {
	Rank    int       `json:"rank" xml:"rank" form:"rank" query:"rank"`
	Updated time.Time `json:"updated" xml:"updated" form:"updated" query:"updated"`
}

//AppInfoHandler アプリ詳細RESTハンドラ
func AppInfoHandler(c echo.Context) error {
	req := new(AppInfoRequest)
	if err := c.Bind(req); err != nil {
		return err
	}
	fmt.Printf("AppInfoHandler request=%v\n", *req)

	//DBアクセスオブジェクト生成
	access, _ := db.NewDBAccess()
	err := access.Open()
	if err != nil {
		return err
	}
	defer access.Close()

	//IDからアプリ情報を取得する
	response := new(AppInfoResponse)
	apps := db.NewApps(access)
	appsRecord, _ := apps.SelectRecord(req.ID)
	if appsRecord.ID != 0 {
		artists := db.NewArtists(access)
		artistsRecord, _ := artists.SelectRecord(appsRecord.ArtistsID)

		response.Name = appsRecord.Name
		response.InfoURL = appsRecord.URL
		response.ArtworkURL = appsRecord.ArtworkURL
		response.ArtistName = artistsRecord.Name
		response.ArtistURL = artistsRecord.URL
		response.Copyright = appsRecord.Copyright
	}

	return c.JSON(http.StatusOK, response)
}

//AppRankHandler アプリ順位RESTハンドラ
func AppRankHandler(c echo.Context) error {
	req := new(AppRankRequest)
	if err := c.Bind(req); err != nil {
		return err
	}
	fmt.Printf("AppRankHandler request=%v\n", *req)

	//DBアクセスオブジェクト生成
	access, _ := db.NewDBAccess()
	err := access.Open()
	if err != nil {
		return err
	}
	defer access.Close()

	//指定アプリのランキング一覧を取得
	ranks := db.NewRanks(access)
	start := rankRangeStringToDate(req.Start)
	end := rankRangeStringToDate(req.End)
	ranksList, err := ranks.SelectAppRankList(start, end, req.Country, req.Kind, req.ID)
	if err != nil {
		return err
	}

	//レスポンス生成
	response := new(AppRankResponse)
	appsResponse := make([]AppRankAppsResponse, 0)
	for _, data := range ranksList {
		appsResponse = append(appsResponse, AppRankAppsResponse{
			Rank:    data.Rank,
			Updated: data.Updated,
		})
	}
	response.Apps = appsResponse

	return c.JSON(http.StatusOK, response)
}

//rankRangeStringToDate ランキングの範囲日文字列を時刻オブジェクトに変換する
func rankRangeStringToDate(rangeDate string) time.Time {
	t, _ := time.Parse("2006-01-02", rangeDate)
	return t
}
