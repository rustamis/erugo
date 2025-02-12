package handlers

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/DeanWard/erugo/config"
	"github.com/DeanWard/erugo/db"
	"github.com/DeanWard/erugo/middleware"
	"github.com/DeanWard/erugo/models"
	"github.com/DeanWard/erugo/progress"
	"github.com/DeanWard/erugo/responses"
	"github.com/DeanWard/erugo/responses/file_response"
	"github.com/DeanWard/erugo/store"
	"github.com/DeanWard/erugo/utils"
	"github.com/gorilla/mux"
	"github.com/yelinaung/go-haikunator"
)

// POST /api/shares/ - create a new share
func CreateShareHandler(database *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handleCreateShare(database, w, r)
	})
}

// GET /api/shares/{longId}/ - get a share by its longId
func GetShareHandler(database *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		longId := mux.Vars(r)["longId"]
		handleGetShare(database, w, longId)
	})
}

// POST /api/shares/{longId}/download/ - download a share by its longId
func DownloadShareHandler(database *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		longId := mux.Vars(r)["longId"]
		handleDownloadShare(database, w, r, longId)
	})
}

func handleGetShare(database *sql.DB, w http.ResponseWriter, longId string) {
	share := db.ShareByLongId(database, longId)
	if share == nil {
		responses.SendResponse(w, responses.StatusError, "Share not found", nil, http.StatusNotFound)
		return
	}

	if IsShareExpired(share) {
		responses.SendResponse(w, responses.StatusError, "Share has expired", nil, http.StatusGone)
		return
	}
	payload := map[string]interface{}{
		"share": share.ToShareResponse(),
	}
	responses.SendResponse(w, responses.StatusSuccess, "Share found", payload, http.StatusOK)
}

func generateUniqueLongId(database *sql.DB) string {
	haikunator := haikunator.New(time.Now().UnixNano())
	longId := haikunator.Haikunate() + "-" + haikunator.Haikunate()
	for db.ShareByLongId(database, longId) != nil {
		longId = haikunator.Haikunate() + "-" + haikunator.Haikunate()
	}
	return longId
}

func handleDownloadShare(database *sql.DB, w http.ResponseWriter, r *http.Request, longId string) {
	share := db.ShareByLongId(database, longId)
	if share == nil {
		responses.SendResponse(w, responses.StatusError, "Share not found", nil, http.StatusNotFound)
		utils.Log("Share not found", utils.ColorRed)
		return
	}

	if IsShareExpired(share) {
		responses.SendResponse(w, responses.StatusError, "Share has expired", nil, http.StatusGone)
		utils.Log("Share has expired", utils.ColorRed)
		return
	}

	if _, err := os.Stat(share.FilePath); os.IsNotExist(err) {
		responses.SendResponse(w, responses.StatusError, "The file you are looking for has been deleted or expired.", nil, http.StatusNotFound)
		utils.Log("The file you are looking for has been deleted or expired.", utils.ColorRed)
		return
	}

	db.CreateShareAccessLogEntry(database, &models.ShareAccessLogEntry{
		ShareId:    share.Id,
		UserEmail:  "not-collected",
		UserIp:     r.RemoteAddr,
		UserAgent:  r.UserAgent(),
		AccessDate: time.Now().Format(time.RFC3339),
	})

	//now we know the folder exists, let's put them in a zip file and serve that
	downloadFilePath, err := store.CreateDownload(share.FilePath)
	if err != nil {
		responses.SendResponse(w, responses.StatusError, "Failed to create download file: "+err.Error(), nil, http.StatusInternalServerError)
		utils.Log("Failed to create download file: "+err.Error(), utils.ColorRed)
		return
	}
	fileName := filepath.Base(downloadFilePath)
	file_response.New(downloadFilePath, fileName).Send(w, r)
}

func IsShareExpired(share *models.Share) bool {
	expirationDate, err := time.Parse(time.RFC3339, share.ExpirationDate)
	if err != nil {
		fmt.Println("Error parsing expiration date:", err)
		return false
	}
	return time.Now().After(expirationDate)
}

func handleCreateShare(database *sql.DB, w http.ResponseWriter, r *http.Request) {
	utils.Log("Creating share", utils.ColorGreen)

	userID := r.Context().Value(middleware.ContextKey("userID")).(int)
	uploadID := r.URL.Query().Get("uploadId")
	if uploadID == "" {
		responses.SendResponse(w, responses.StatusError, "Upload ID required", nil, http.StatusBadRequest)
		return
	}

	reader, err := r.MultipartReader()
	if err != nil {
		responses.SendResponse(w, responses.StatusError, "Failed to create multipart reader", nil, http.StatusBadRequest)
		return
	}

	// Initialize progress tracking
	tracker := progress.GetTracker()
	tracker.NewUpload(uploadID)
	defer tracker.DeleteUpload(uploadID)

	maxShareSize := config.GetMaxShareSize()
	if maxShareSize == 0 {
		responses.SendResponse(w, responses.StatusError, "Max share size not set", nil, http.StatusInternalServerError)
		return
	}

	totalSize := r.ContentLength // Total size of the upload
	utils.Log(fmt.Sprintf("Total size: %d", totalSize), utils.ColorGreen)
	longId := generateUniqueLongId(database)
	uploadDir := filepath.Join(config.AppConfig.BaseStoragePath, longId)
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		responses.SendResponse(w, responses.StatusError, "Failed to create upload directory", nil, http.StatusInternalServerError)
		return
	}

	var fileNames []string
	var actualTotalSize int64
	var totalBytesRead int64

	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			responses.SendResponse(w, responses.StatusError, "Error reading upload", nil, http.StatusBadRequest)
			return
		}
		defer part.Close()

		if part.FileName() == "" {
			continue
		}

		destPath := filepath.Join(uploadDir, part.FileName())
		dest, err := os.Create(destPath)
		if err != nil {
			responses.SendResponse(w, responses.StatusError, "Failed to create file", nil, http.StatusInternalServerError)
			return
		}
		defer dest.Close()

		// Create progress reader with both current file and total size
		pr := NewProgressReader(part, 1, totalSize, totalBytesRead, uploadID) // per-file progress isn't working yet so just use 1 for now

		written, err := io.Copy(dest, pr)
		if err != nil {
			responses.SendResponse(w, responses.StatusError, "Failed to save file", nil, http.StatusInternalServerError)
			return
		}

		totalBytesRead += written
		actualTotalSize += written
		fileNames = append(fileNames, part.FileName())

		if actualTotalSize > maxShareSize {
			os.RemoveAll(uploadDir)
			responses.SendResponse(w, responses.StatusError, "Total size of files exceeds maximum allowed", nil, http.StatusBadRequest)
			return
		}
	}

	expirationDate := time.Now().AddDate(0, 0, 7).Format(time.RFC3339)
	share := &models.Share{
		FilePath:       uploadDir,
		ExpirationDate: expirationDate,
		LongId:         longId,
		NumFiles:       len(fileNames),
		TotalSize:      actualTotalSize,
		Files:          fileNames,
		UserId:         userID,
	}

	savedShare, err := db.ShareCreate(database, share)
	if err != nil {
		os.RemoveAll(uploadDir)
		responses.SendResponse(w, responses.StatusError, "Failed to create share", nil, http.StatusInternalServerError)
		return
	}

	payload := map[string]interface{}{
		"share": savedShare.ToShareResponse(),
	}
	responses.SendResponse(w, responses.StatusSuccess, "Share created", payload, http.StatusOK)
}
