package progress

import (
	"sync"
	"time"
)

type Progress struct {
	BytesRead      int64   `json:"bytesRead"`      // Current file bytes read
	TotalSize      int64   `json:"totalSize"`      // Current file total size
	Percentage     float64 `json:"percentage"`     // Current file percentage
	TotalBytesRead int64   `json:"totalBytesRead"` // Total bytes read across all files
	TotalFileSize  int64   `json:"totalFileSize"`  // Total size of all files
	TotalProgress  float64 `json:"totalProgress"`  // Overall progress percentage
	UploadID       string  `json:"uploadId"`
	LastUpdate     time.Time
}

type ProgressTracker struct {
	mu      sync.RWMutex
	uploads map[string]chan Progress
	cleanup map[string]chan struct{}
}

var tracker *ProgressTracker
var once sync.Once

func GetTracker() *ProgressTracker {
	once.Do(func() {
		tracker = &ProgressTracker{
			uploads: make(map[string]chan Progress),
			cleanup: make(map[string]chan struct{}),
		}
	})
	return tracker
}

func (t *ProgressTracker) NewUpload(uploadID string) chan Progress {
	t.mu.Lock()
	defer t.mu.Unlock()

	progressChan := make(chan Progress, 100)
	cleanupChan := make(chan struct{})
	t.uploads[uploadID] = progressChan
	t.cleanup[uploadID] = cleanupChan

	// Automatic cleanup after 1 hour
	go func() {
		select {
		case <-cleanupChan:
			// Normal cleanup
		case <-time.After(1 * time.Hour):
			// Timeout cleanup
		}
		t.DeleteUpload(uploadID)
	}()

	return progressChan
}

func (t *ProgressTracker) GetUploadChannel(uploadID string) (chan Progress, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	ch, exists := t.uploads[uploadID]
	return ch, exists
}

func (t *ProgressTracker) DeleteUpload(uploadID string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if cleanupChan, exists := t.cleanup[uploadID]; exists {
		close(cleanupChan)
		delete(t.cleanup, uploadID)
	}

	if progressChan, exists := t.uploads[uploadID]; exists {
		close(progressChan)
		delete(t.uploads, uploadID)
	}
}
