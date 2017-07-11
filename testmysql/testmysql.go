package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func CheckErrorAndPanic(err error) {
	if err != nil {
		panic(err)
	}
}

func InitDB(dsn string) *sql.DB {
	db, err := sql.Open("mysql", dsn)
	CheckErrorAndPanic(err)
	return db
}

func main() {
	db := InitDB("root:111111@tcp(localhost:3306)/test?charset=utf8")
	CreateTable(db)
}

func CreateTable(db *sql.DB) {
	sql_table := `
    create table if not exists dev(
        id int unsigned auto_increment not null primary key,
        uid varchar(64),
        did varchar(64),
        qid varchar(64),
        name varchar(64),
        status char default 'u'
        );
    `
	stmt, err := db.Prepare(sql_table)
	CheckErrorAndPanic(err)
	defer stmt.Close()

	result, err := stmt.Exec()
	CheckErrorAndPanic(err)
	fmt.Println(result)
}
