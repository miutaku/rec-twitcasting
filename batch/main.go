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

type DBConfig struct {
	Host      string
	Port      string
	User      string
	Password  string
	DbName    string
	TableName string
}

func getDBConfig() DBConfig {
	return DBConfig{
		Host:      os.Getenv("DB_HOST"),
		Port:      os.Getenv("DB_PORT"),
		User:      os.Getenv("DB_USER"),
		Password:  os.Getenv("DB_PASSWORD"),
		DbName:    os.Getenv("DB_NAME"),
		TableName: os.Getenv("DB_TABLE_NAME"),
	}
}

func task() {
	config := getDBConfig()
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.DbName,
	))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	checkQuery := fmt.Sprintf("SELECT username FROM %s WHERE recording_state = false", config.TableName)
	rows, err := db.Query(checkQuery)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Query failed: %v\n", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to scan row: %v\n", err)
			return
		}
		resp, err := http.Get(fmt.Sprintf("http://api-rec-twitcasting:8080/check-live?username=%s", name))
		if err != nil {
			fmt.Fprintf(os.Stderr, "HTTP request failed: %v\n", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			fmt.Printf("User %s is live\n", name)
		} else {
			fmt.Printf("User %s is not live\n", name)
		}
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
