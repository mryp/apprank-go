package db

import "fmt"

//テーブル名
const (
	artistsTableName = "artists"
)

type Artists struct {
	access *DBAccess
}

type ArtistsTable struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
	URL  string `db:"url"`
}

func NewArtists(access *DBAccess) *Artists {
	artists := new(Artists)
	artists.access = access
	return artists
}

func (artists *Artists) Insert(record ArtistsTable) error {
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

func (artists *Artists) SelectRecord(id int64) (ArtistsTable, error) {
	var resultList []ArtistsTable
	_, err := artists.access.session.Select("*").
		From(artistsTableName).
		Where("id = ?", id).
		Limit(1).
		Load(&resultList)
	if err != nil {
		fmt.Printf("selectRecord err=%v\n", err)
		return ArtistsTable{}, err
	}
	if len(resultList) == 0 {
		return ArtistsTable{}, fmt.Errorf("データが見つからない")
	}

	return resultList[0], nil
}
