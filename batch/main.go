package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-co-op/gocron"
	_ "github.com/lib/pq"
)

func task() {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	rows, err := db.Query("SELECT name FROM twitcasting.speakers WHERE recording=false")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Query failed: %v\n", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Row scan failed: %v\n", err)
			continue
		}

		fmt.Printf("Checking live status for user: %s\n", name)
		resp, err := http.Get(fmt.Sprintf("http://rec-twitcasting:8080/check-live?username=%s", name))
		if err != nil {
			fmt.Fprintf(os.Stderr, "HTTP request failed: %v\n", err)
			continue
		}
		resp.Body.Close()
		fmt.Printf("Checked live status for user: %s\n", name)
	}
}

func main() {
	s1 := gocron.NewScheduler(time.Local)

	s1.Every(2).Seconds().Do(task)
	s1.StartAsync()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
}
