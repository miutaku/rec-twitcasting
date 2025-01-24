package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/nobuf/cas"
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
		os.Exit(1)
	}
	defer db.Close()

	query := fmt.Sprintf("UPDATE %s SET recording_state = $1 WHERE username = $2", config.TableName)
	_, err = db.Exec(query, state, username)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	clientID := os.Getenv("TWITCASTING_CLIENT_ID")
	clientSecret := os.Getenv("TWITCASTING_CLIENT_SECRET")
	if clientID == "" || clientSecret == "" {
		log.Fatal("TWITCASTING_CLIENT_ID or TWITCASTING_CLIENT_SECRET is not set")
	}

	// TwitCastingクライアントを初期化
	api := cas.New(clientID, clientSecret)

	// 出力ディレクトリを環境変数から取得 (デフォルト: "/recorded")
	baseDir := os.Getenv("OUTPUT_DIR")
	if baseDir == "" {
		baseDir = "./recorded"
	}

	// HTTPルーターの設定
	r := mux.NewRouter()
	r.HandleFunc("/check-live", func(w http.ResponseWriter, r *http.Request) {
		// `username` パラメータを取得
		username := r.URL.Query().Get("username")
		if username == "" {
			http.Error(w, "username parameter is required", http.StatusBadRequest)
			return
		}

		// 配信状況を確認
		liveInfo, err := api.UserCurrentLive(username)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to get current live information: %v", err), http.StatusInternalServerError)
			return
		}

		if liveInfo.Movie.ID == "" {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, "User is not live streaming.")
			return
		}

		// 配信中の場合、recording_stateをTRUEに更新
		if err := updateRecordingState(username, true); err != nil {
			http.Error(w, fmt.Sprintf("Failed to update recording state: %v", err), http.StatusInternalServerError)
			return
		}

		// 配信中のタイトルを取得
		title := liveInfo.Movie.Title
		fmt.Printf("User is live streaming. Title: %s\n", title)

		// 現在の日時を取得
		now := time.Now()
		dateDir := now.Format("2006-01-02")  // YYYY-MM-DD
		timePrefix := now.Format("15-04")    // HH-mm
		safeTitle := sanitizeFilename(title) // タイトルの利用できない文字を置換
		outputPath := filepath.Join(baseDir, username, dateDir)
		outputFile := filepath.Join(outputPath, fmt.Sprintf("%s_%s.mp4", timePrefix, safeTitle))

		// 出力ディレクトリを作成 (存在しない場合)
		if err := os.MkdirAll(outputPath, os.ModePerm); err != nil {
			http.Error(w, fmt.Sprintf("Failed to create output directory: %v", err), http.StatusInternalServerError)
			return
		}

		// ffmpegで録画
		streamURL := fmt.Sprintf("https://twitcasting.tv/%s/metastream.m3u8/?video=1", username)
		fmt.Printf("Starting ffmpeg to record the stream...\n")
		cmd := exec.Command("ffmpeg", "-y", "-i", streamURL, "-c", "copy", "-map", "p:0", outputFile)
		if os.Getenv("LOG_LEVEL") == "debug" {
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
		}

		if err := cmd.Run(); err != nil {
			http.Error(w, fmt.Sprintf("Failed to execute ffmpeg: %v", err), http.StatusInternalServerError)
			return
		}
		// 録画終了後、recording_stateをFALSEに更新
		if err := updateRecordingState(username, false); err != nil {
			http.Error(w, fmt.Sprintf("Failed to update recording state: %v", err), http.StatusInternalServerError)
			return
		}
		fmt.Printf("Recording finished. Saved as: %s\n", outputFile)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Recording finished. Saved as: %s\n", outputFile)
	}).Methods(http.MethodGet)

	// サーバーの起動
	port := ":8080"
	fmt.Printf("Server is running at http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, r))
}

// sanitizeFilename は、タイトルをファイル名として安全に利用できる形式に変換します
func sanitizeFilename(name string) string {
	// 不正な文字を削除
	replacer := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|", " "}
	for _, r := range replacer {
		name = replaceAll(name, r, "_")
	}
	return name
}

func replaceAll(str, old, new string) string {
	return string([]rune(str))
}
