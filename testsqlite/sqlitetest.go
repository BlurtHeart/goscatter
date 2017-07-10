package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

type TestItem struct {
	Id         int
	Username   string
	Departname string
	Created    *time.Time
}

func main() {
	db := InitDB("./foo.db")

	err := createTable(db)
	checkError(err)

	item := TestItem{Username: "abc", Departname: "安全部门", Created: "2017-07-10"}
	items := []TestItem{item}

	StoreItem(db, items)
	id, err := res.LastInsertId()
	checkError(err)
	fmt.Println(id)
}

func InitDB(filepath string) *sql.DB {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		panic(err)
	}
	if db == nil {
		panic("db nil")
	}
	return db
}

func createTable(db *sql.DB) error {
	sql_table := `
    CREATE TABLE IF NOT EXISTS users(
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username VARCHAR(20) NOT NULL,
        departname VARCHAR(64) NULL,
        created DATE NULL
        );
    `
	_, err := db.Exec(sql_table)
	return err
}

func StoreItem(db *sql.DB, items []TestItem) {
	sql_additem := `
    INSERT INTO users(username, departname, created) values(?,?,?);
    `
	stmt, err := db.Prepare(sql_additem)
	checkError(err)
	defer stmt.Close()

	for _, item := range items {
		_, err2 := stmt.Exec(item.Username, item.Departname, item.Created)
		checkError(err2)
	}
}

func ReadItem(db *sql.DB) []TestItem {
	sql_readall := `
    SELECT userame, departname, created FROM users
    ORDER BY datetime(created) DESC
    `

	rows, err := db.Query(sql_readall)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var result []TestItem
	for rows.Next() {
		item := TestItem{}
		err2 := rows.Scan(&item.Id, &item.Username, &item.Departname)
		if err2 != nil {
			panic(err2)
		}
		result = append(result, item)
	}
	return result
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
