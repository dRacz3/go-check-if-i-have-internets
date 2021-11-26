package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type ParsedArgs struct {
	poll_time              int
	purge_history_on_start bool
}

func parse_args() ParsedArgs {
	// variables declaration
	var time int

	// flags declaration using flag package
	flag.IntVar(&time, "t", 1, "Polling time in s. Default is 1s")
	purge_history_on_start := flag.Bool("p", false, "Purge history on start")
	flag.Parse()

	parsed := ParsedArgs{time, *purge_history_on_start}
	log.Println("Parsed args: ", parsed)
	return parsed
}

/*
This small program is intended to check whether a given URL is up or down.
Essentially, i'm assuming that url for 9gag is always available and contains
the name of the site, if it is not, it means
that the internet connection is down. This is used, so I can contously check
if vodafone has a shitty service with a lot of unnoticed downtime or not.
*/
func main() {
	args := parse_args()
	database_filename := "sqlitedb.db"
	if args.purge_history_on_start {
		log.Println("Purging previous data")
		os.Remove(database_filename)
	}
	db, err := sql.Open("sqlite3", database_filename)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	writer := ResultStreamWriter{*db}
	writer.CreateTable()

	UPTIME_CHECK_LINK := "http://9gag.com"
	UPTIME_CHECK_STRING := "9gag"
	TIME_BETWEEN_CHECKS := time.Second * time.Duration(args.poll_time)

	link_checker := LinkChecker{UPTIME_CHECK_LINK, UPTIME_CHECK_STRING,
		http.Client{
			Timeout: 2 * time.Second,
		}}
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
