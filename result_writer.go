package main

import (
	"database/sql"
	"log"
	"time"
)

type ResultStreamWriter struct {
	db sql.DB
}

func (r *ResultStreamWriter) CreateTable() {
	_, err := r.db.Exec("CREATE TABLE IF NOT EXISTS results(id INTEGER PRIMARY KEY, timestamp TEXT , target TEXT ,result BOOLEAN)")
	if err != nil {
		log.Fatal(err)
	}

}

func (r *ResultStreamWriter) writeResult(target string, result bool) {
	_, err := r.db.Exec("INSERT INTO results(timestamp, target, result) VALUES(?, ?, ?)", time.Now(), target, result)
	if err != nil {
		log.Fatal(err)
	}
}
