package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"github.com/robfig/cron/v3"
)

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	return value
}

func checkLive(name string) {
	url := fmt.Sprintf("http://api-rec-twitcasting:8080/check-live?username=%s", name)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Failed to request %s: %v", url, err)
		return
	}
	defer resp.Body.Close()
	log.Printf("Checked live status for %s, response status: %s", name, resp.Status)
}

func fetchAndCheckLive(db *sql.DB) {
	rows, err := db.Query("SELECT name FROM twitcasting.speakers WHERE recording=false")
	if err != nil {
		log.Printf("Failed to execute query: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			log.Printf("Failed to scan row: %v", err)
			continue
		}
		checkLive(name)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating rows: %v", err)
	}
}

func main() {
	connStr := getEnv("DATABASE_URL", "postgres://rec-twitcasting-user:rec-twitcasting-pass@postgres-rec-twitcasting/dbname?sslmode=disable")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	c := cron.New()
	interval := getEnv("FETCH_INTERVAL_SEC", "60")
	c.AddFunc("@every "+interval+"s", func() { fetchAndCheckLive(db) })
	c.Start()
}
