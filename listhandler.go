package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo"
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

	//とりあえずダミーをセット
	response := new(NowResponse)
	response.Updated, _ = time.Parse("2006-01-02 15:04:05", "2017-09-01 10:00:00")
	apps := make([]NowAppsResponce, 0)
	apps = append(apps, NowAppsResponce{
		ID:         1016318735,
		Name:       "アイドルマスター シンデレラガールズ スターライトステージ",
		ArtworkURL: "http://is3.mzstatic.com/image/thumb/Purple118/v4/6d/e7/40/6de74073-31d0-5dbb-4962-6f3ebb288514/source/200x200bb.png",
	})
	apps = append(apps, NowAppsResponce{
		ID:         1015521325,
		Name:       "Fate/Grand Order",
		ArtworkURL: "http://is2.mzstatic.com/image/thumb/Purple118/v4/dd/2e/a4/dd2ea4d0-4435-c073-dfa6-f1df8c26bfc5/source/200x200bb.png",
	})
	apps = append(apps, NowAppsResponce{
		ID:         658511662,
		Name:       "モンスターストライク",
		ArtworkURL: "http://is1.mzstatic.com/image/thumb/Purple128/v4/9b/5e/61/9b5e61fc-8b30-7ab4-3e0f-c0bb82e3a8eb/source/200x200bb.png",
	})

	response.Apps = apps
	return c.JSON(http.StatusOK, response)
}
