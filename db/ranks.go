package db

import (
	"fmt"
	"time"

	"github.com/gocraft/dbr"
)

//テーブル名
const (
	ranksTableName        = "ranks"
	RanksKindGrossing     = 1
	RanksKindGrossingIpad = 2
	RanksKindPaid         = 3
	RanksKindPaidIpad     = 4
	RanksKindFree         = 5
	RanksKindFreeIpad     = 6
)

//情報テーブル
type Ranks struct {
	session *dbr.Session
}

type RanksTable struct {
	ID      int64     `db:"id"`
	Updated time.Time `db:"updated"`
	Country string    `db:"country"`
	Kind    int       `db:"kind"`
	Rank    int       `db:"rank"`
	AppID   int64     `db:"app_id"`
}

func NewRanks(session *dbr.Session) (*Ranks, error) {
	ranks := new(Ranks)
	session, err := ConnectDBRecheck(session)
	if err != nil {
		return nil, err
	}
	ranks.session = session
	return ranks, nil
}

func (ranks *Ranks) Close() {
	if ranks.session != nil {
		defer ranks.session.Close()
	}
}

func (ranks *Ranks) Insert(record RanksTable) error {
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
	_, err = ranks.session.InsertInto(ranksTableName).
		Columns("updated", "country", "kind", "rank", "app_id").
		Record(record).
		Exec()
	if err != nil {
		return err
	}

	return nil
}

func (ranks *Ranks) selectAppRank(updated time.Time, country string, kind int, appID int64) (int, error) {
	var resultList []RanksTable
	_, err := ranks.session.Select("*").
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

func (ranks *Ranks) SelectLatestUpdated(country string, kind int) (time.Time, error) {
	var resultList []RanksTable
	_, err := ranks.session.Select("*").
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

func (ranks *Ranks) SelectRankList(updated time.Time, country string, kind int) ([]RanksTable, error) {
	var resultList []RanksTable
	_, err := ranks.session.Select("*").
		From(ranksTableName).
		Where("updated = ? AND country = ? AND kind = ?", updated, country, kind).
		OrderDir("rank", true).
		Load(&resultList)
	if err != nil {
		return nil, err
	}

	return resultList, nil
}
