package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/mryp/apprank-go/db"
)

const (
	RssBaseURL          = "https://rss.itunes.apple.com/api/v1/ios-apps/"
	RssKindGrossing     = "top-grossing"
	RssKindGrossingIpad = "top-grossing-ipad"
	RssKindPaid         = "top-paid"
	RssKindPaidIpad     = "top-paid-ipad"
	RssKindFree         = "top-free"
	RssKindFreeIpad     = "top-free-ipad"
	RssListCountMax     = 10
	RssAPIName          = "explicit.json"
)

var (
	rankWatcher *RankWatcher
)

//RssReedItem はRSSデータの中身オブジェクトを表す構造体
//変換ツール：https://mholt.github.io/json-to-go/
type RssFeed struct {
	Feed struct {
		Title  string `json:"title"`
		ID     string `json:"id"`
		Author struct {
			Name string `json:"name"`
			URI  string `json:"uri"`
		} `json:"author"`
		Links []struct {
			Self      string `json:"self,omitempty"`
			Alternate string `json:"alternate,omitempty"`
		} `json:"links"`
		Copyright string `json:"copyright"`
		Country   string `json:"country"`
		Icon      string `json:"icon"`
		Updated   string `json:"updated"`
		Results   []struct {
			ArtistID      string `json:"artistId"`
			ArtistName    string `json:"artistName"`
			ArtistURL     string `json:"artistUrl"`
			ArtworkURL100 string `json:"artworkUrl100"`
			Copyright     string `json:"copyright"`
			Genres        []struct {
				GenreID string `json:"genreId"`
				Name    string `json:"name"`
				URL     string `json:"url"`
			} `json:"genres"`
			ID          string `json:"id"`
			Kind        string `json:"kind"`
			Name        string `json:"name"`
			ReleaseDate string `json:"releaseDate"`
			URL         string `json:"url"`
		} `json:"results"`
	} `json:"feed"`
}

type RankWatcher struct {
	isBusy bool
}

func NewRankWatcher() *RankWatcher {
	if rankWatcher != nil {
		return rankWatcher
	}

	rankWatcher := new(RankWatcher)
	rankWatcher.isBusy = false
	return rankWatcher
}

func (watcher *RankWatcher) StartBgTask() {
	go func() {
		fmt.Printf("RankWatcher.StartBgTask 処理開始")
		watcher.UpdateRanking("jp", RssKindGrossing)
	}()
}

func (watcher *RankWatcher) UpdateRanking(country string, kind string) {
	fmt.Printf("RankWatcher.UpdateRanking 処理開始")
	url := fmt.Sprintf("https://rss.itunes.apple.com/api/v1/%s/ios-apps/%s/%d/%s",
		country, kind, RssListCountMax, RssAPIName)
	fmt.Printf("url=%s\n", url)
	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("データ取得失敗 %v\n", err.Error())
		return
	}
	defer response.Body.Close()

	rss := new(RssFeed)
	err = json.NewDecoder(response.Body).Decode(rss)
	if err != nil {
		fmt.Printf("データ変換 %v\n", err.Error())
		return
	}

	access, _ := db.NewDBAccess()
	err = access.Open()
	if err != nil {
		fmt.Printf("DBオープンエラー %v\n", err.Error())
		return
	}
	defer access.Close()

	ranks := db.NewRanks(access)
	artists := db.NewArtists(access)

	updated := stringToTime(rss.Feed.Updated)
	dbKind := rssKindToDBKind(kind)
	latestUpdated, _ := ranks.SelectLatestUpdated(country, dbKind)
	if updated.UTC() == latestUpdated.UTC() {
		fmt.Printf("対象時刻データは登録済み\n")
		return
	}

	for i, data := range rss.Feed.Results {
		//ランキングを登録
		rank := i + 1
		ranksRecord := db.RanksTable{Updated: updated, Country: country, Kind: dbKind, Rank: rank, AppID: strToInt64(data.ID)}
		err = ranks.Insert(ranksRecord)
		if err != nil {
			fmt.Printf("Rankテーブル登録エラー %v\n", err.Error())
			continue
		}

		//著作者情報を登録・更新
		artistsRecord := db.ArtistsTable{ID: strToInt64(data.ArtistID), Name: data.ArtistName, URL: data.ArtistURL}
		artists.Insert(artistsRecord)
		if err != nil {
			fmt.Printf("Artistsテーブル登録エラー %v\n", err.Error())
			continue
		}

		//アプリ情報を登録・更新

	}
}

func stringToTime(rssTime string) time.Time {
	//2017-09-04T01:39:08.000-07:00
	t, _ := time.Parse(time.RFC3339, rssTime)
	return t
}

func rssKindToDBKind(rssKind string) int {
	kind := 0
	switch rssKind {
	case RssKindGrossing:
		kind = db.RanksKindGrossing
	case RssKindGrossingIpad:
		kind = db.RanksKindGrossingIpad
	case RssKindPaid:
		kind = db.RanksKindPaid
	case RssKindPaidIpad:
		kind = db.RanksKindPaidIpad
	case RssKindFree:
		kind = db.RanksKindFree
	case RssKindFreeIpad:
		kind = db.RanksKindFreeIpad
	}
	return kind
}

func strToInt64(text string) int64 {
	i64, _ := strconv.ParseInt(text, 10, 64)
	return i64
}
