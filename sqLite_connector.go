package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

func load_database() *sql.DB {
	database, e := sql.Open("sqlite3", "./try.db")
	checkError(e)
	statement, e := database.Prepare("CREATE TABLE IF NOT EXISTS HSK1 (english TEXT PRIMARY KEY, " +
		"chinese TEXT, pinyin TEXT)")
	checkError(e)
	statement.Exec()
	return database
}

func insert(chinese string, english string, pinyin string) {
	database := load_database()
	statement, e := database.Prepare("INSERT INTO HSK1 (english, chinese, pinyin) VALUES (?,?,?)")
	checkError(e)
	statement.Exec(english, chinese, pinyin)
}

func query(chinese string) (string, string) {
	database := load_database()
	row := database.QueryRow("SELECT english, pinyin FROM HSK1 WHERE chinese=?", chinese)
	var english string
	var pinyin string
	row.Scan(&english, &pinyin)
	return english, pinyin
}

func checkError(e error) {
	if e != nil {
		fmt.Println("Encounter Error: " + e.Error())
		os.Exit(1)
	}
}
