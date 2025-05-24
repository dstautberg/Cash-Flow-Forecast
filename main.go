package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	_ "modernc.org/sqlite"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"os"
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

	// Create transactions table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS transactions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		account_number TEXT,
		description TEXT,
		transaction_date TEXT,
		transaction_type TEXT,
		transaction_amount REAL,
		balance REAL
	)`)
	if err != nil {
		panic(fmt.Sprintf("failed to create transactions table: %v", err))
	}

	// Locale-aware number formatting example
	p := message.NewPrinter(language.English)
	amount := 1234567.89
	fmt.Println("Locale-aware formatted number:")
	p.Printf("%v\n", amount)

	// Find the newest CSV file in the current directory
	dirEntries, err := os.ReadDir(".")
	if err != nil {
		panic(fmt.Sprintf("failed to read directory: %v", err))
	}
	var newestFile string
	var newestModTime int64
	for _, entry := range dirEntries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if len(name) > 4 && name[len(name)-4:] == ".csv" {
			info, err := entry.Info()
			if err != nil {
				continue
			}
			modTime := info.ModTime().Unix()
			if modTime > newestModTime || newestFile == "" {
				newestModTime = modTime
				newestFile = name
			}
		}
	}
	if newestFile == "" {
		panic("no CSV files found in the current directory")
	}

	fmt.Printf("Loading newest CSV file: %s\n", newestFile)
	csvFile, err := os.Open(newestFile)
	if err != nil {
		panic(fmt.Sprintf("failed to open CSV file: %v", err))
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	records, err := reader.ReadAll()
	if err != nil {
		panic(fmt.Sprintf("failed to read CSV file: %v", err))
	}

	fmt.Printf("Read %d rows from CSV file.\n", len(records)-1) // Exclude header
	// Optionally print the first few rows
	for i, row := range records {
		if i > 5 { // Print only first 5 rows for preview
			break
		}
		fmt.Println(row)
	}

	// Truncate the transactions table at the beginning
	_, err = db.Exec("DELETE FROM transactions")
	if err != nil {
		panic(fmt.Sprintf("failed to truncate transactions table: %v", err))
	}

	// Insert each CSV row into the transactions table (skip header)
	for i, row := range records {
		if i == 0 {
			continue // skip header
		}
		if len(row) < 6 {
			continue // skip malformed rows
		}
		_, err := db.Exec(`INSERT INTO transactions (account_number, description, transaction_date, transaction_type, transaction_amount, balance) VALUES (?, ?, ?, ?, ?, ?)`,
			row[0], row[1], row[2], row[3], row[4], row[5],
		)
		if err != nil {
			fmt.Printf("failed to insert row %d: %v\n", i, err)
		}
	}
	fmt.Printf("Inserted %d transactions into the database.\n", len(records)-1)

	// Query the transactions table, ordered by transaction_date
	fmt.Println("\nTransactions ordered by date:")
	rows, err := db.Query(`SELECT account_number, description, transaction_date, transaction_type, transaction_amount, balance FROM transactions ORDER BY transaction_date`)
	if err != nil {
		panic(fmt.Sprintf("failed to query transactions: %v", err))
	}
	defer rows.Close()

	for rows.Next() {
		var accountNumber, description, transactionDate, transactionType string
		var transactionAmount, balance float64
		err := rows.Scan(&accountNumber, &description, &transactionDate, &transactionType, &transactionAmount, &balance)
		if err != nil {
			fmt.Printf("failed to scan row: %v\n", err)
			continue
		}
		p.Printf("%s | %s | %s | %s | %v | %v\n", accountNumber, description, transactionDate, transactionType, transactionAmount, balance)
	}
	if err := rows.Err(); err != nil {
		fmt.Printf("row iteration error: %v\n", err)
	}
}
