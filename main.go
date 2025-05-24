package main

import (
	"database/sql"
	"fmt"
	_ "modernc.org/sqlite"
)

func main() {
	db, err := sql.Open("sqlite", "file:cashflow.db?cache=shared&mode=rwc")
	if err != nil {
		panic(fmt.Sprintf("failed to open database: %v", err))
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS example (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL
	)`)
	if err != nil {
		panic(fmt.Sprintf("failed to create table: %v", err))
	}

	fmt.Println("SQLite database initialized without CGO.")
}
