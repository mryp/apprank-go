package db

import (
	"fmt"
	"time"
)

const (
	appsTableName = "apps"
)

//Apps はアプリ情報テーブルアクセス構造体
type Apps struct {
	access *DBAccess
}

//AppsRecord はアプリ情報テーブルレコード構造体
type AppsRecord struct {
	ID          int64     `db:"id"`
	Name        string    `db:"name"`
	URL         string    `db:"url"`
	ArtworkURL  string    `db:"artwork_url"`
	Kind        string    `db:"kind"`
	Copyright   string    `db:"copyright"`
	ArtistsID   int64     `db:"artist_id"`
	ReleaseDate time.Time `db:"release_date"`
}

//NewApps はアプリ情報テーブルを初期返して返す
func NewApps(access *DBAccess) *Apps {
	apps := new(Apps)
	apps.access = access
	return apps
}

//Insert はアプリ情報テーブルにレコードを追加する
//既に登録されている場合はデータの更新を行う
func (apps *Apps) Insert(record AppsRecord) error {
	fmt.Print(fmt.Sprintf("Apps.Insert ID:%d, Name:%s, URL:%s, ArtworkURL:%s, Kind:%s, Copyright:%s, ArtistsID:%d, ReleaseDate:%s\n",
		record.ID, record.Name, record.URL, record.ArtworkURL, record.Kind, record.Copyright, record.ArtistsID, record.ReleaseDate))

	hitRecord, err := apps.SelectRecord(record.ID)
	if hitRecord.Name == record.Name &&
		hitRecord.URL == record.URL &&
		hitRecord.ArtworkURL == record.ArtworkURL &&
		hitRecord.Kind == record.Kind &&
		hitRecord.Copyright == record.Copyright &&
		hitRecord.ArtistsID == record.ArtistsID &&
		hitRecord.ReleaseDate == record.ReleaseDate {
		fmt.Printf("変更なしのため登録しない")
		return nil
	}

	if hitRecord.ID == 0 {
		_, err = apps.access.session.InsertInto(appsTableName).
			Columns("id", "name", "url", "artwork_url", "kind", "copyright", "artist_id", "release_date").
			Record(record).
			Exec()
	} else {
		_, err = apps.access.session.Update(appsTableName).
			Set("name", record.Name).
			Set("url", record.URL).
			Set("artwork_url", record.ArtworkURL).
			Set("kind", record.Kind).
			Set("copyright", record.Copyright).
			Set("artist_id", record.ArtistsID).
			Set("release_date", record.ReleaseDate).
			Where("id = ?", record.ID).
			Exec()
	}
	if err != nil {
		return err
	}

	return nil
}

//SelectRecord は指定したアプリIDからアプリ詳細情報レコードを取得する
func (apps *Apps) SelectRecord(id int64) (AppsRecord, error) {
	var resultList []AppsRecord
	_, err := apps.access.session.Select("*").
		From(appsTableName).
		Where("id = ?", id).
		Limit(1).
		Load(&resultList)
	if err != nil {
		fmt.Printf("selectRecord err=%v\n", err)
		return AppsRecord{}, err
	}
	if len(resultList) == 0 {
		return AppsRecord{}, fmt.Errorf("データが見つからない")
	}

	return resultList[0], nil
}
