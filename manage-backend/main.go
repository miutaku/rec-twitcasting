package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/lib/pq"
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
	// CORSミドルウェアを追加
	http.ListenAndServe(":8888", corsMiddleware(http.DefaultServeMux))
}

// CORSミドルウェア
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// curl "localhost:8888/list-casting-users"
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
			"recording_state":  recordingState,
			"action_date_time": createdDateTime,
			"action":           "listCastingUser",
			"target_username":  username,
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

// curl "localhost:8888/add-casting-user?username=<username>"
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
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" { // 23505 is the unique_violation error code
			response := map[string]interface{}{
				"error":            "User already exists",
				"action_date_time": time.Now().UTC().Format(time.RFC3339),
				"action":           "addCastingUser",
				"target_username":  username,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}
		http.Error(w, "Failed to insert user", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"action_date_time": time.Now().UTC().Format(time.RFC3339),
		"action":           "addCastingUser",
		"target_username":  username,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// curl "localhost:8888/del-casting-user?username=<username>"
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

	response := map[string]interface{}{
		"action_date_time": time.Now().UTC().Format(time.RFC3339),
		"action":           "deleteCastingUser",
		"target_username":  username,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// curl "localhost:8888/check-recording-state?username=<username>"
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

	response := map[string]interface{}{
		"recording_state":  recordingState,
		"action_date_time": time.Now().UTC().Format(time.RFC3339),
		"action":           "checkRecordingState",
		"target_username":  username,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// curl "localhost:8888/update-recording-state?username=<username>&recording_state=<false/true>"
func updateRecordingStateHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "username is required", http.StatusBadRequest)
		return
	}

	recordingState := r.URL.Query().Get("recording_state")
	if recordingState == "" {
		http.Error(w, "recording_state is required", http.StatusBadRequest)
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

	query := fmt.Sprintf("UPDATE %s SET recording_state = $1 WHERE username = $2", config.TableName)
	_, err = db.Exec(query, recordingState, username)
	if err != nil {
		http.Error(w, "Failed to update recording state", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"action_date_time": time.Now().UTC().Format(time.RFC3339),
		"action":           "updateRecordingState",
		"target_username":  username,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
