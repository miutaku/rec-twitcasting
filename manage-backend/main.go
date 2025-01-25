package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

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
		Host:      getEnv("DB_HOST", "localhost"),
		Port:      getEnv("DB_PORT", "5432"), // Default to 5432 if DB_PORT is not set
		User:      getEnv("DB_USER", "user"),
		Password:  getEnv("DB_PASSWORD", "password"),
		DbName:    getEnv("DB_NAME", "dbname"),
		TableName: getEnv("DB_TABLE_NAME", "tablename"),
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func main() {
	http.HandleFunc("/list-casting-users", listCastingUsersHandler)
	http.HandleFunc("/add-casting-user", addCastingUserHandler)
	http.HandleFunc("/del-casting-user", delCastingUserHandler)
	http.HandleFunc("/check-recording-state", checkRecordingStateHandler)
	http.HandleFunc("/update-recording-state", updateRecordingStateHandler)
	log.Fatal(http.ListenAndServe(":8888", nil))
}

func listCastingUsersHandler(w http.ResponseWriter, r *http.Request) {
	config := getDBConfig()
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.DbName,
	))
	if err != nil {
		http.Error(w, "Failed to connect to the database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	query := fmt.Sprintf("SELECT username, recording_state, created_date_time FROM %s", config.TableName)
	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, "Failed to query users", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []map[string]interface{}
	for rows.Next() {
		var username string
		var recordingState bool
		var createdDateTime string
		err := rows.Scan(&username, &recordingState, &createdDateTime)
		if err != nil {
			http.Error(w, "Failed to scan row", http.StatusInternalServerError)
			return
		}
		user := map[string]interface{}{
			"username":          username,
			"recording_state":   recordingState,
			"created_date_time": createdDateTime,
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, "Error iterating rows", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
func addCastingUserHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "username is required", http.StatusBadRequest)
		return
	}

	config := getDBConfig()
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.DbName,
	))
	if err != nil {
		http.Error(w, "Failed to connect to the database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	query := fmt.Sprintf("INSERT INTO %s (username) VALUES($1)", config.TableName)
	_, err = db.Exec(query, username)
	if err != nil {
		http.Error(w, "Failed to insert user", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "User %s added successfully", username)
}

func delCastingUserHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "username is required", http.StatusBadRequest)
		return
	}

	config := getDBConfig()
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.DbName,
	))
	if err != nil {
		http.Error(w, "Failed to connect to the database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	query := fmt.Sprintf("DELETE FROM %s WHERE username = $1", config.TableName)
	_, err = db.Exec(query, username)
	if err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "User %s deleted successfully", username)
}

func checkRecordingStateHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "username is required", http.StatusBadRequest)
		return
	}

	config := getDBConfig()
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.DbName,
	))
	if err != nil {
		http.Error(w, "Failed to connect to the database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var recordingState string
	query := fmt.Sprintf("SELECT recording_state FROM %s WHERE username = $1", config.TableName)
	err = db.QueryRow(query, username).Scan(&recordingState)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to query recording state", http.StatusInternalServerError)
		}
		return
	}

	fmt.Fprintf(w, "Recording state for user %s is %s", username, recordingState)
}

func updateRecordingStateHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "username is required", http.StatusBadRequest)
		return
	}

	state := r.URL.Query().Get("state")
	if state == "" {
		http.Error(w, "state is required", http.StatusBadRequest)
		return
	}

	recordingState := state == "true"

	err := updateRecordingState(username, recordingState)
	if err != nil {
		http.Error(w, "Failed to update recording state", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Recording state for user %s updated successfully", username)
}

func updateRecordingState(username string, state bool) error {
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
		return err
	}
	defer db.Close()

	query := fmt.Sprintf("UPDATE %s SET recording_state = $1 WHERE username = $2", config.TableName)
	_, err = db.Exec(query, state, username)
	if err != nil {
		return err
	}

	return nil
}
