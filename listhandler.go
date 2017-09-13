package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/mryp/apprank-go/db"
)

//NowRequest は最新ランキング一覧RESTのリクエスト
type NowRequest struct {
	Country string `json:"country" xml:"country" form:"country" query:"country"`
	Kind    int    `json:"kind" xml:"index" form:"kind" query:"kind"`
}

//NowResponse は最新ランキング一覧RESTのレスポンス
type NowResponse struct {
	Updated time.Time         `json:"updated" xml:"updated"`
	Apps    []NowAppsResponce `json:"apps" xml:"apps"`
}

//NowAppsResponce は最新ランキング一覧RESTのアプリ情報部のレスポンス
type NowAppsResponce struct {
	ID         int64  `json:"id" xml:"id"`
	Name       string `json:"name" xml:"name"`
	ArtworkURL string `json:"artwork_url" xml:"artwork_url"`
}

//NowHandler は最新ランキング一覧ハンドラ
func NowHandler(c echo.Context) error {
	req := new(NowRequest)
	if err := c.Bind(req); err != nil {
		return err
	}
	fmt.Printf("NowHandler request=%v\n", *req)

	//DBアクセスオブジェクト生成
	access, _ := db.NewDBAccess()
	err := access.Open()
	if err != nil {
		return err
	}
	defer access.Close()

	//ランキング一覧取得
	ranks := db.NewRanks(access)
	updated, err := ranks.SelectLatestUpdated(req.Country, req.Kind)
	if err != nil {
		return err
	}

	//レスポンス生成
	response := new(NowResponse)
	response.Updated = updated
	rankList, err := ranks.SelectRankList(updated, req.Country, req.Kind)
	if err != nil {
		fmt.Printf("ランキング一覧取得失敗 err=%s\n", err)
		return err
	}
	if rankList == nil {
		return fmt.Errorf("ランキングデータが見つかりません")
	}

	//アプリ一覧用情報生成
	apps := db.NewApps(access)
	appsResponse := make([]NowAppsResponce, 0)
	for _, data := range rankList {
		record, _ := apps.SelectRecord(data.AppID)
		name := record.Name
		artworkURL := record.ArtworkURL
		appsResponse = append(appsResponse, NowAppsResponce{
			ID:         data.AppID,
			Name:       name,
			ArtworkURL: artworkURL,
		})
	}
	response.Apps = appsResponse

	return c.JSON(http.StatusOK, response)
}
