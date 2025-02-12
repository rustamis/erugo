package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/DeanWard/erugo/progress"
	"github.com/gorilla/mux"
)

func UploadProgressHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uploadID := mux.Vars(r)["uploadId"]
		if uploadID == "" {
			http.Error(w, "Upload ID required", http.StatusBadRequest)
			return
		}

		// Set headers for SSE
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Get the progress channel for this upload
		tracker := progress.GetTracker()
		progressChan, exists := tracker.GetUploadChannel(uploadID)
		if !exists {
			http.Error(w, "Upload not found", http.StatusNotFound)
			return
		}

		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
			return
		}

		// Use the request's context to detect client disconnect
		ctx := r.Context()

		for {
			select {
			case <-ctx.Done():
				// Client disconnected
				tracker.DeleteUpload(uploadID)
				return

			case progress, ok := <-progressChan:
				if !ok {
					// Channel closed - upload complete or failed
					return
				}

				// Send progress update
				data, _ := json.Marshal(progress)
				fmt.Fprintf(w, "data: %s\n\n", data)
				flusher.Flush()
			}
		}
	})
}

type ProgressReader struct {
	Reader        io.Reader
	Size          int64 // Current file size
	TotalFileSize int64 // Size of all files
	bytesRead     int64 // Current file bytes read
	totalRead     int64 // Total bytes read across all files
	lastUpdate    time.Time
	uploadID      string
	tracker       *progress.ProgressTracker
}

func NewProgressReader(reader io.Reader, size int64, totalSize int64, totalRead int64, uploadID string) *ProgressReader {
	return &ProgressReader{
		Reader:        reader,
		Size:          size,
		TotalFileSize: totalSize,
		totalRead:     totalRead,
		uploadID:      uploadID,
		tracker:       progress.GetTracker(),
		lastUpdate:    time.Now(),
	}
}

func (pr *ProgressReader) Read(p []byte) (int, error) {
	n, err := pr.Reader.Read(p)

	pr.bytesRead += int64(n)
	pr.totalRead += int64(n)

	// Update progress every 500ms
	if time.Since(pr.lastUpdate) > 500*time.Millisecond {
		if progressChan, exists := pr.tracker.GetUploadChannel(pr.uploadID); exists {

			totalProgress := float64(pr.totalRead) / float64(pr.TotalFileSize) * 100

			select {
			case progressChan <- progress.Progress{
				BytesRead:      pr.bytesRead,
				TotalSize:      pr.Size,
				TotalBytesRead: pr.totalRead,
				TotalFileSize:  pr.TotalFileSize,
				TotalProgress:  totalProgress,
				UploadID:       pr.uploadID,
				LastUpdate:     time.Now(),
			}:

			default:
			}
		}
		pr.lastUpdate = time.Now()
	}

	return n, err
}
