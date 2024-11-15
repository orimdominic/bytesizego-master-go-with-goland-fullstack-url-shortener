package db

import (
	"database/sql"
)

func Connect(url string) (*sql.DB, error) {
	conn, err := sql.Open("sqlite3", url)
	if err != nil {
		return nil, err
	}

	err = conn.Ping()
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func CreateUrlTable(db *sql.DB) error {
	statement := `CREATE TABLE IF NOT EXISTS urls (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		hash VARCHAR(10) NOT NULL,
		url VARCHAR TEXT NOT NULL
	);`

	_, err := db.Exec(statement)
	if err != nil {
		return err
	}

	return nil
}

func CreateUrlRecord(db *sql.DB, url string, hash string) (res sql.Result, err error) {
	// TODO check if exists, if exists, no need, else insert
	res, err = db.Exec(`INSERT INTO urls (hash, url) VALUES (?, ?) RETURNING *;`, hash, url)
	return
}

func GetUrlRecord(db *sql.DB, hash string) (string, error) {
	var url string

	err := db.QueryRow(`SELECT url FROM urls WHERE hash = ?`, hash).Scan(&url)
	if err != nil {
		return "", err
	}

	return url, nil
}