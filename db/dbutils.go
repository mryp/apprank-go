package db

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql" //dbrで使用する
	"github.com/gocraft/dbr"
)

const (
	dbUserID     = "root"
	dbPassword   = "root"
	dbHostName   = "127.0.0.1"
	dbPortNumber = "3306"
	dbName       = "apprank"
)

//ConnectDB DB接続
func ConnectDB() (*dbr.Session, error) {
	db, err := dbr.Open("mysql", dbUserID+":"+dbPassword+"@tcp("+dbHostName+":"+dbPortNumber+")/"+dbName+"?parseTime=true", nil)
	if err != nil {
		fmt.Printf("connectDB err=%v\n", err)
		return nil, err
	}

	dbsession := db.NewSession(nil)
	return dbsession, nil
}

func ConnectDBRecheck(session *dbr.Session) (*dbr.Session, error) {
	if session == nil {
		newSession, err := ConnectDB()
		if err != nil {
			return nil, err
		}
		session = newSession
	}

	return session, nil
}
