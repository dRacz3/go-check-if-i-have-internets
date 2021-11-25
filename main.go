package main

import (
	"database/sql"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type LinkChecker struct {
	target_url   string
	check_string string
}

func (l LinkChecker) checkLink() bool {
	resp, err := http.Get(l.target_url)
	if err != nil {
		return false
	}
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return false
	}
	return strings.Contains(string(content), l.check_string)
}

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

func main() {
	database_filename := "sqlitedb.db"

	os.Remove(database_filename)

	db, err := sql.Open("sqlite3", database_filename)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	writer := ResultStreamWriter{*db}
	writer.CreateTable()

	UPTIME_CHECK_LINK := "http://9gag.com"
	UPTIME_CHECK_STRING := "9gag"
	TIME_BETWEEN_CHECKS := time.Second * 2

	link_checker := LinkChecker{UPTIME_CHECK_LINK, UPTIME_CHECK_STRING}
	other_args := os.Args[1:]
	log.Println("Starting up with args: ")
	for _, arg := range other_args {
		log.Println(string(arg))
	}

	for {
		suc := link_checker.checkLink()
		if suc {
			log.Println("Link is up")
		} else {
			log.Println("Link is down")
		}
		writer.writeResult(UPTIME_CHECK_LINK, suc)
		time.Sleep(TIME_BETWEEN_CHECKS)
	}
}
