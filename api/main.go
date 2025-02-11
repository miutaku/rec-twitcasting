package main

import (
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

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func getTwitcastingAPI() (*cas.Client, error) {
	clientID := os.Getenv("TWITCASTING_CLIENT_ID")
	clientSecret := os.Getenv("TWITCASTING_CLIENT_SECRET")
	if clientID == "" || clientSecret == "" {
		return nil, fmt.Errorf("TWITCASTING_CLIENT_ID or TWITCASTING_CLIENT_SECRET is not set")
	}
	api := cas.New(clientID, clientSecret)
	return api, nil
}

func updateRecordingState(username string, state bool) error {
	manageBackendHost := getEnv("MANAGE_BACKEND_HOST", "manage-backend-rec-twitcasting:8888")
	updateURL := fmt.Sprintf("http://%s/update-recording-state?username=%s&recording_state=%t", manageBackendHost, username, state)
	resp, err := http.Get(updateURL)
	if err != nil {
		return fmt.Errorf("failed to update recording state: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update recording state, status code: %d", resp.StatusCode)
	}
	return nil
}

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

func handleCheckLive(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "username parameter is required", http.StatusBadRequest)
		return
	}

	api, err := getTwitcastingAPI()
	if err != nil {
		log.Fatal(err)
	}

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

	if err := updateRecordingState(username, true); err != nil {
		http.Error(w, fmt.Sprintf("Failed to update recording state: %v", err), http.StatusInternalServerError)
		return
	}

	title := liveInfo.Movie.Title
	fmt.Printf("User is live streaming. Title: %s\n", title)

	now := time.Now()
	dateDir := now.Format("2006-01-02")
	timePrefix := now.Format("15-04")
	safeTitle := sanitizeFilename(title)
	baseDir := os.Getenv("OUTPUT_DIR")
	if baseDir == "" {
		baseDir = "./recorded"
	}
	outputPath := filepath.Join(baseDir, username, dateDir)
	outputFile := filepath.Join(outputPath, fmt.Sprintf("%s_%s.mp4", timePrefix, safeTitle))

	if err := os.MkdirAll(outputPath, os.ModePerm); err != nil {
		http.Error(w, fmt.Sprintf("Failed to create output directory: %v", err), http.StatusInternalServerError)
		return
	}

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

	if err := updateRecordingState(username, false); err != nil {
		http.Error(w, fmt.Sprintf("Failed to update recording state: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Recording finished. Saved as: %s\n", outputFile)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/check-live", handleCheckLive).Methods(http.MethodGet)

	port := ":8080"
	fmt.Printf("Server is running at http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, r))
}
