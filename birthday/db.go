package birthday

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var db *sqlx.DB

func ConnectToSQLite(dbFilePath string) {
	var err error
	db, err = sqlx.Connect("sqlite3", dbFilePath)
	if err != nil {
		log.Fatalln(err)
	}
}