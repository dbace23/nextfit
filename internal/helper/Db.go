package helper

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func GetMysqlConnection() *sql.DB {
	db, err := sql.Open("mysql", "root:@/nextfit?parseTime=true")
	if err != nil {
		panic(err)
	}

	return db
}
