package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/mryp/apprank-go/db"
)

type NowRequest struct {
	Country string `json:"country" xml:"country" form:"country" query:"country"`
	Kind    int    `json:"kind" xml:"index" form:"kind" query:"kind"`
}

type NowResponse struct {
	Updated time.Time         `json:"updated" xml:"updated"`
	Apps    []NowAppsResponce `json:"apps" xml:"apps"`
}

type NowAppsResponce struct {
	ID         int64  `json:"id" xml:"id"`
	Name       string `json:"name" xml:"name"`
	ArtworkURL string `json:"artwork_url" xml:"artwork_url"`
}

func NowHandler(c echo.Context) error {
	req := new(NowRequest)
	if err := c.Bind(req); err != nil {
		return err
	}
	fmt.Printf("NowHandler request=%v\n", *req)

	access, _ := db.NewDBAccess()
	err := access.Open()
	if err != nil {
		return err
	}
	defer access.Close()

	ranks := db.NewRanks(access)
	updated, err := ranks.SelectLatestUpdated(req.Country, req.Kind)
	if err != nil {
		return err
	}

	//とりあえずダミーをセット
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

	apps := make([]NowAppsResponce, 0)
	for _, data := range rankList {
		apps = append(apps, NowAppsResponce{
			ID:         data.AppID,
			Name:       "名前はまだない",
			ArtworkURL: "http://hogehoge/",
		})
	}
	response.Apps = apps

	return c.JSON(http.StatusOK, response)
}
