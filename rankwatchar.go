package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/mryp/apprank-go/config"
	"github.com/mryp/apprank-go/db"
)

const (
	//RssBaseURL はRSS取得のベースURL
	RssBaseURL = "https://rss.itunes.apple.com/api/v1/ios-apps/"
	//RssAPIName はRSS取得の実行API名
	RssAPIName = "explicit.json"

	//RssKindGrossing はRSS取得種別のセールスランキングの設定値
	RssKindGrossing = "top-grossing"
	//RssKindGrossingIpad はRSS取得種別のiPadセールスランキングの設定値
	RssKindGrossingIpad = "top-grossing-ipad"
	//RssKindPaid はRSS取得種別の有料ランキングの設定値
	RssKindPaid = "top-paid"
	//RssKindPaidIpad はRSS取得種別のiPad有料ランキングの設定値
	RssKindPaidIpad = "top-paid-ipad"
	//RssKindFree はRSS取得種別の無料ランキングの設定値
	RssKindFree = "top-free"
	//RssKindFreeIpad はRSS取得種別のiPad無料ランキングの設定値
	RssKindFreeIpad = "top-free-ipad"
)

var (
	//rankWatcher はインスタンス化したものを保持しておく
	rankWatcher *RankWatcher
)

//RssFeed はRSSデータの中身オブジェクトを表す構造体
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

//RankWatcher はランキング監視用構造体
type RankWatcher struct {
	isBusy bool
}

//NewRankWatcher はランキングオブジェクトを生成して返す
//既に作成済みの時はそれを返す
func NewRankWatcher() *RankWatcher {
	if rankWatcher != nil {
		return rankWatcher
	}

	rankWatcher := new(RankWatcher)
	rankWatcher.isBusy = false
	return rankWatcher
}

//StartBgTask はランキングの取得とDB設定タスクを非同期で実行する
func (watcher *RankWatcher) StartBgTask() {
	go func() {
		fmt.Printf("RankWatcher.StartBgTask 処理開始")
		watcher.UpdateRanking("jp", RssKindGrossing)
	}()
}

//UpdateRanking はランキングを取得しDBを更新する
func (watcher *RankWatcher) UpdateRanking(country string, kind string) {
	fmt.Printf("RankWatcher.UpdateRanking 処理開始")

	//ランキングRSSを取得
	url := fmt.Sprintf("https://rss.itunes.apple.com/api/v1/%s/ios-apps/%s/%d/%s",
		country, kind, config.Now().Watch.MaxCount, RssAPIName)
	fmt.Printf("url=%s\n", url)
	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("データ取得失敗 %v\n", err.Error())
		return
	}
	defer response.Body.Close()

	//RSSをオブジェクトに変換
	rss := new(RssFeed)
	err = json.NewDecoder(response.Body).Decode(rss)
	if err != nil {
		fmt.Printf("データ変換 %v\n", err.Error())
		return
	}

	//DBアクセスオブジェクトを作成
	access, _ := db.NewDBAccess()
	err = access.Open()
	if err != nil {
		fmt.Printf("DBオープンエラー %v\n", err.Error())
		return
	}
	defer access.Close()
	ranks := db.NewRanks(access)
	artists := db.NewArtists(access)
	apps := db.NewApps(access)

	//最新のデータかどうか確認
	updated := stringToTime(rss.Feed.Updated)
	dbKind := rssKindToDBKind(kind)
	latestUpdated, _ := ranks.SelectLatestUpdated(country, dbKind)
	fmt.Printf("updated=%s\n", updated.UTC())
	fmt.Printf("latestUpdated=%s\n", latestUpdated.UTC())
	if updated.UTC() == latestUpdated.UTC() {
		fmt.Printf("対象時刻データは登録済み\n")
		return
	}

	for i, data := range rss.Feed.Results {
		//ランキングを登録
		rank := i + 1
		ranksRecord := db.RanksRecord{Updated: updated,
			Country: country,
			Kind:    dbKind,
			Rank:    rank,
			AppID:   stringToInt64(data.ID)}
		err = ranks.Insert(ranksRecord)
		if err != nil {
			fmt.Printf("Rankテーブル登録エラー %v\n", err.Error())
			continue
		}

		//著作者情報を登録・更新
		artistsRecord := db.ArtistsRecord{ID: stringToInt64(data.ArtistID),
			Name: data.ArtistName,
			URL:  data.ArtistURL}
		err = artists.Insert(artistsRecord)
		if err != nil {
			fmt.Printf("Artistsテーブル登録エラー %v\n", err.Error())
			continue
		}

		//アプリ情報を登録・更新
		appsRecord := db.AppsRecord{ID: stringToInt64(data.ID),
			Name:        data.Name,
			URL:         data.URL,
			ArtworkURL:  data.ArtworkURL100,
			Kind:        data.Kind,
			Copyright:   data.Copyright,
			ArtistsID:   stringToInt64(data.ArtistID),
			ReleaseDate: stringToDate(data.ReleaseDate)}
		err = apps.Insert(appsRecord)
		if err != nil {
			fmt.Printf("Appsテーブル登録エラー %v\n", err.Error())
			continue
		}
	}
}

//stringToTime はRSSの時刻文字列（RFC3339）を時刻データに変換する
func stringToTime(rssTime string) time.Time {
	t, _ := time.Parse(time.RFC3339, rssTime)
	return t
}

//stringToDate はRSSの日付文字列（YYYY-MM-DD）を時刻データに変換する
func stringToDate(rssTime string) time.Time {
	t, _ := time.Parse("2006-01-02", rssTime)
	return t
}

//rssKindToDBKind はRSSの取得種別をDBのランキング種別に変換する
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

//stringToInt64 はRSS文字列をint64に変換する
func stringToInt64(text string) int64 {
	i64, _ := strconv.ParseInt(text, 10, 64)
	return i64
}
