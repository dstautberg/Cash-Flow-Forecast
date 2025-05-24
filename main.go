package main

import (
	"database/sql"
	"fmt"
	_ "modernc.org/sqlite"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
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

	// Locale-aware number formatting example
	p := message.NewPrinter(language.English)
	amount := 1234567.89
	fmt.Println("Locale-aware formatted number:")
	p.Printf("%v\n", amount)
}
