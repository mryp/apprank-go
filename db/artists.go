package db

import "fmt"

//テーブル名
const (
	artistsTableName = "artists"
)

//Artists は著作者情報テーブルアクセス構造体
type Artists struct {
	access *DBAccess
}

//ArtistsRecord は著作者情報テーブルのレコード構造体
type ArtistsRecord struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
	URL  string `db:"url"`
}

//NewArtists は著作者情報テーブルアクセス構造体を生成初期化して返す
func NewArtists(access *DBAccess) *Artists {
	artists := new(Artists)
	artists.access = access
	return artists
}

//Insert は著作者情報テーブルにレコードを追加する
//すでに登録されている場合はデータの更新を行う
func (artists *Artists) Insert(record ArtistsRecord) error {
	fmt.Print(fmt.Sprintf("Artists.Insert ID:%d, Name:%s, URL:%s\n",
		record.ID, record.Name, record.URL))

	hitRecord, err := artists.SelectRecord(record.ID)
	if hitRecord.Name == record.Name && hitRecord.URL == record.URL {
		fmt.Printf("変更なしのため登録しない")
		return nil
	}

	if hitRecord.ID == 0 {
		_, err = artists.access.session.InsertInto(artistsTableName).
			Columns("id", "name", "url").
			Record(record).
			Exec()
	} else {
		_, err = artists.access.session.Update(artistsTableName).
			Set("name", record.Name).
			Set("url", record.URL).
			Where("id = ?", record.ID).
			Exec()
	}
	if err != nil {
		return err
	}

	return nil
}

//SelectRecord 指定した著作者IDのレコードを取得する
func (artists *Artists) SelectRecord(id int64) (ArtistsRecord, error) {
	var resultList []ArtistsRecord
	_, err := artists.access.session.Select("*").
		From(artistsTableName).
		Where("id = ?", id).
		Limit(1).
		Load(&resultList)
	if err != nil {
		fmt.Printf("selectRecord err=%v\n", err)
		return ArtistsRecord{}, err
	}
	if len(resultList) == 0 {
		return ArtistsRecord{}, fmt.Errorf("データが見つからない")
	}

	return resultList[0], nil
}
