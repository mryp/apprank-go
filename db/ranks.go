package db

import (
	"fmt"
	"time"
)

const (
	//ranksTableName はランキングテーブル名
	ranksTableName = "ranks"

	//RanksKindGrossing は種別のセールスランキング
	RanksKindGrossing = 1
	//RanksKindGrossingIpad は種別のiPadセールスランキング
	RanksKindGrossingIpad = 2
	//RanksKindPaid は種別の有料ランキング
	RanksKindPaid = 3
	//RanksKindPaidIpad は種別のiPad有料ランキング
	RanksKindPaidIpad = 4
	//RanksKindFree は種別の無料ランキング
	RanksKindFree = 5
	//RanksKindFreeIpad は種別のiPad無料スランキング
	RanksKindFreeIpad = 6
)

//Ranks はRanksテーブルアクセス用構造体
type Ranks struct {
	access *DBAccess
}

//RanksRecord はRanksテーブルレコードを表す構造体
type RanksRecord struct {
	ID      int64     `db:"id"`
	Updated time.Time `db:"updated"`
	Country string    `db:"country"`
	Kind    int       `db:"kind"`
	Rank    int       `db:"rank"`
	AppID   int64     `db:"app_id"`
}

//NewRanks はDBアクセス情報使用してアクセス用構造体を生成して返す
func NewRanks(access *DBAccess) *Ranks {
	ranks := new(Ranks)
	ranks.access = access
	return ranks
}

//Insert はランキング情報を登録する
//既に登録済みの時は登録しない
func (ranks *Ranks) Insert(record RanksRecord) error {
	if record.ID != 0 {
		return fmt.Errorf("パラメーターエラー")
	}
	fmt.Print(fmt.Sprintf("Ranks.Insert ID:%d, Updated:%s, Country:%s, Kind:%d, Rank:%d, AppID:%d\n",
		record.ID, record.Updated, record.Country, record.Kind, record.Rank, record.AppID))

	//すでに同じ時刻で登録されている時は何もしない
	rank, err := ranks.selectAppRank(record.Updated, record.Country, record.Kind, record.AppID)
	if err != nil {
		return err
	}
	if rank != 0 {
		return fmt.Errorf("既にデータ登録ずみ")
	}

	//登録
	_, err = ranks.access.session.InsertInto(ranksTableName).
		Columns("updated", "country", "kind", "rank", "app_id").
		Record(record).
		Exec()
	if err != nil {
		return err
	}

	return nil
}

//selectAppRank は指定した時刻のアプリランキング順位を取得する
func (ranks *Ranks) selectAppRank(updated time.Time, country string, kind int, appID int64) (int, error) {
	var resultList []RanksRecord
	_, err := ranks.access.session.Select("*").
		From(ranksTableName).
		Where("updated = ? AND country = ? AND kind = ? AND app_id = ?",
			updated, country, kind, appID).
		Limit(1).
		Load(&resultList)
	if err != nil {
		fmt.Printf("selectAppRank err=%v\n", err)
		return 0, err
	}
	if len(resultList) == 0 {
		return 0, nil
	}

	return resultList[0].Rank, nil
}

//SelectAppRankList は指定した時間範囲のアプリランキング順位リストを取得する
func (ranks *Ranks) SelectAppRankList(start time.Time, end time.Time, country string, kind int, appID int64) ([]RanksRecord, error) {
	var resultList []RanksRecord
	_, err := ranks.access.session.Select("*").
		From(ranksTableName).
		Where("updated >= ? AND updated < ? AND country = ? AND kind = ? AND app_id = ?",
			start, end, country, kind, appID).
		OrderDir("rank", true).
		Load(&resultList)
	if err != nil {
		return nil, err
	}

	return resultList, nil
}

//SelectLatestUpdated は指定した種別の最後の更新日時を取得する
func (ranks *Ranks) SelectLatestUpdated(country string, kind int) (time.Time, error) {
	var resultList []RanksRecord
	_, err := ranks.access.session.Select("*").
		From(ranksTableName).
		Where("country = ? AND kind = ?", country, kind).
		OrderDir("updated", false).
		Limit(1).
		Load(&resultList)
	if err != nil {
		return time.Time{}, err
	}
	if len(resultList) == 0 {
		return time.Time{}, fmt.Errorf("検索対象データなし")
	}

	return resultList[0].Updated, nil
}

//SelectRankList は指定した時刻に登録されているアプリランキング一覧を取得する
func (ranks *Ranks) SelectRankList(updated time.Time, country string, kind int) ([]RanksRecord, error) {
	var resultList []RanksRecord
	_, err := ranks.access.session.Select("*").
		From(ranksTableName).
		Where("updated = ? AND country = ? AND kind = ?", updated, country, kind).
		OrderDir("rank", true).
		Load(&resultList)
	if err != nil {
		return nil, err
	}

	return resultList, nil
}
