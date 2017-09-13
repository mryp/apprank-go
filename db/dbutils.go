package db

import (
	_ "github.com/go-sql-driver/mysql" //dbrで使用する
	"github.com/gocraft/dbr"

	"github.com/mryp/apprank-go/config"
)

const (
	dbUserID     = "root"
	dbPassword   = "root"
	dbHostName   = "127.0.0.1"
	dbPortNumber = "3306"
	dbName       = "apprank"
)

//DBAccess はDBアクセス用構造体
type DBAccess struct {
	session *dbr.Session
}

//NewDBAccess はDBアクセス構造体を初期化して返す
func NewDBAccess() (*DBAccess, error) {
	access := new(DBAccess)
	access.session = nil
	return access, nil
}

//Open はDBの接続を行いセッションをDBアクセスに保存する
func (access *DBAccess) Open() error {
	dbConfig := config.Now().DB
	db, err := dbr.Open("mysql", dbConfig.UserID+":"+dbConfig.Password+"@tcp("+dbConfig.HostName+":"+dbConfig.PortNumber+")/"+dbConfig.Name+"?parseTime=true", nil)
	if err != nil {
		return err
	}

	access.session = db.NewSession(nil)
	return nil
}

//Close はDBの切断を行う
func (access *DBAccess) Close() {
	if access.session != nil {
		access.session.Close()
		access.session = nil
	}
}
