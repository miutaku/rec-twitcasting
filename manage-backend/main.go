package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type DBConfig struct {
	Host         string
	Port         string
	User         string
	Password     string
	DbName       string
	TableName    string
	ApiTableName string
}

func getDBConfig() DBConfig {
	return DBConfig{
		Host:         getEnv("DB_HOST", "localhost"),
		Port:         getEnv("DB_PORT", "5432"), // Default to 5432 if DB_PORT is not set
		User:         getEnv("DB_USER", "user"),
		Password:     getEnv("DB_PASSWORD", "password"),
		DbName:       getEnv("DB_NAME", "dbname"),
		TableName:    getEnv("DB_TABLE_NAME", "tablename"),
		ApiTableName: getEnv("DB_API_TABLE_NAME", "api_key"),
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
	http.HandleFunc("/get-twitcasting-code", getOauthCodeHandler) // Handle Twitcasting API OAuth2 callback endpoint
	http.HandleFunc("/alert-expire-token", alertExpireToken)
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

// Handle OAuth callback
func getOauthCodeHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code") // Get the code from the query parameter
	if code == "" {
		http.Error(w, "Failed to get twitcasting api code.", http.StatusInternalServerError)
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

	if code != "" {
		// Prepare the request payload
		data := url.Values{}
		data.Set("code", code)
		data.Set("grant_type", "authorization_code")
		data.Set("client_id", os.Getenv("CLIENT_ID"))
		data.Set("client_secret", os.Getenv("CLIENT_SECRET"))
		data.Set("redirect_uri", os.Getenv("BACKEND_SERVER"))

		// Make the POST request
		resp, err := http.PostForm("https://apiv2.twitcasting.tv/oauth2/access_token", data)
		if err != nil {
			http.Error(w, "Failed to request access token", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		// Parse the response
		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			http.Error(w, "Failed to parse access token response", http.StatusInternalServerError)
			return
		}

		// Extract the access token and expires_in
		expiresIn, ok := result["expires_in"].(float64)
		if !ok {
			http.Error(w, "Failed to get expires_in", http.StatusInternalServerError)
			return
		}
		accessToken, ok := result["access_token"].(string)
		if !ok {
			http.Error(w, "Failed to get access token", http.StatusInternalServerError)
			return
		}

		// Convert expires_in to timestamp
		expirationTime := time.Now().Add(time.Duration(expiresIn) * time.Second).Unix()

		// Store the access token and expiration time in the database
		query := fmt.Sprintf("TRUNCATE TABLE %s; INSERT INTO %s (code, access_token, expires_in) VALUES ($1, $2, $3)", config.ApiTableName, config.ApiTableName)
		_, err = db.Exec(query, code, accessToken, expirationTime)
		if err != nil {
			http.Error(w, "Failed to store access token in database", http.StatusInternalServerError)
			return
		}
		html := fmt.Sprintf(`
				<!DOCTYPE html>
				<html>
					<body>
						<h2>Getting your twitcasting api token successful. Code is %s.</h2>
						<p>Code has been stored your DB server in %s.(DB: %s, Table: %s)</p>
					</body>
				</html>`, code, getDBConfig().Host, getDBConfig().DbName, getDBConfig().ApiTableName)
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, html)
		return
	}
	if code == "" {
		html := fmt.Sprintf(`
				<!DOCTYPE html>
				<html>
					<body>
						<h2>Failed to get twitcasting api code. Code is null.</h2>
						<p>Code has not been stored your DB server in %s.(DB: %s, Table: %s)</p>
					</body>
				</html>`, getDBConfig().Host, getDBConfig().DbName, getDBConfig().ApiTableName)
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, html)
		return
	}
}

func alertExpireToken(w http.ResponseWriter, r *http.Request) {
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

	var expiresIn int64
	query := fmt.Sprintf("SELECT expires_in FROM %s", config.ApiTableName)
	err = db.QueryRow(query).Scan(&expiresIn)
	if err != nil {
		http.Error(w, "Failed to query expires_in", http.StatusInternalServerError)
		return
	}

	// Check if the token is expiring in a week
	if time.Now().Unix() > expiresIn-7*24*60*60 {
		if os.Getenv("LINE_CHANNEL_ACCESS_TOKEN") != "" && os.Getenv("LINE_USER_ID") != "" {
			alertToLine(w, r)
		}
	}

	response := map[string]interface{}{
		"action_date_time": time.Now().UTC().Format(time.RFC3339),
		"action":           "alertExpireToken",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func alertToLine(w http.ResponseWriter, _ *http.Request) {
	// Prepare the webhook payload
	webhookURL := "https://api.line.me/v2/bot/message/push"
	channelAccessToken := os.Getenv("LINE_CHANNEL_ACCESS_TOKEN")
	toUserID := os.Getenv("LINE_USER_ID")
	message := map[string]interface{}{
		"to": toUserID,
		"messages": []map[string]string{
			{
				"type": "text",
				"text": "[Action Required] The access token is expiring in a week.",
			},
		},
	}
	payload, err := json.Marshal(message)
	if err != nil {
		http.Error(w, "Failed to marshal webhook payload", http.StatusInternalServerError)
		return
	}

	// Send the webhook notification
	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(payload))
	if err != nil {
		http.Error(w, "Failed to create webhook request", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+channelAccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to send webhook request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Webhook request failed", http.StatusInternalServerError)
		return
	}
}
