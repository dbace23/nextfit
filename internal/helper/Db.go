package helper

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func GetMysqlConnection() *sql.DB {
	// db, err := sql.Open("mysql", "root:@/nextfit?parseTime=true")
	db, err := sql.Open("mysql", "nextfit:F2598676t!@tcp(127.0.0.1:3306)/nextfit?parseTime=true&loc=Local")
	if err != nil {
		panic(err)
	}

	return db
}
