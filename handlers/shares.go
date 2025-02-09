package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/DeanWard/erugo/config"
	"github.com/DeanWard/erugo/db"
	"github.com/DeanWard/erugo/middleware"
	"github.com/DeanWard/erugo/models"
	"github.com/DeanWard/erugo/responses"
	"github.com/DeanWard/erugo/responses/file_response"
	"github.com/DeanWard/erugo/store"
	"github.com/DeanWard/erugo/utils"
	"github.com/go-playground/validator/v10"
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

func handleCreateShare(database *sql.DB, w http.ResponseWriter, r *http.Request) {

	utils.Log("Creating share", utils.ColorGreen)

	userID := r.Context().Value(middleware.ContextKey("userID")).(int)
	utils.Log(fmt.Sprintf("User ID: %d", userID), utils.ColorGreen)

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	validate := validator.New()
	files := r.MultipartForm.File["files"]
	req := models.CreateShareRequest{
		Files: files,
	}

	if err := validate.Struct(req); err != nil {
		payload := map[string]interface{}{
			"errors": extractValidationErrors(err),
		}
		responses.SendResponse(w, responses.StatusError, "Validation error", payload, http.StatusBadRequest)
		//log the error
		utils.Log(fmt.Sprintf("Validation error: %v", payload), utils.ColorRed)
		return
	}

	maxShareSize := config.GetMaxShareSize()
	utils.Log(fmt.Sprintf("Max share size: %d", maxShareSize), utils.ColorGreen)
	if maxShareSize == 0 {
		responses.SendResponse(w, responses.StatusError, "Max share size not set", nil, http.StatusInternalServerError)
		utils.Log("Max share size not set", utils.ColorRed)
		return
	}

	//check if the total size of the files in the request is greater than the max share size
	totalRequestSize := r.ContentLength
	if totalRequestSize > maxShareSize {
		responses.SendResponse(w, responses.StatusError, "Total size of files is greater than the max share size", nil, http.StatusBadRequest)
		utils.Log("Total size of files is greater than the max share size", utils.ColorRed)
		return
	}

	longId := generateUniqueLongId(database)
	folderPath := store.StoreUploadedFiles(config.AppConfig.BaseStoragePath, files, longId)
	if folderPath == "" {
		responses.SendResponse(w, responses.StatusError, "There was an error saving your files.", nil, http.StatusBadRequest)
		utils.Log("There was an error saving users files.", utils.ColorRed)
		return
	}
	totalSize := store.GetTotalSize(folderPath)
	fileNames := store.GetFilesInFolder(folderPath)

	expirationDate := time.Now().AddDate(0, 0, 7).Format(time.RFC3339)
	share := &models.Share{
		FilePath:       folderPath,
		ExpirationDate: expirationDate,
		LongId:         longId,
		NumFiles:       len(files),
		TotalSize:      totalSize,
		Files:          fileNames,
		UserId:         userID,
	}
	savedShare, err := db.ShareCreate(database, share)
	if err != nil {
		responses.SendResponse(w, responses.StatusError, "Failed to create share", nil, http.StatusInternalServerError)
		utils.Log("Failed to create share", utils.ColorRed)
		return
	}
	payload := map[string]interface{}{
		"share": savedShare.ToShareResponse(),
	}
	responses.SendResponse(w, responses.StatusSuccess, "Share created", payload, http.StatusOK)
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

func extractValidationErrors(err error) map[string]string {
	errorMap := make(map[string]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			errorMap[strings.ToLower(e.Field())] = e.Tag() // e.Tag() gives the failed rule (e.g., required, email)
		}
	}

	return errorMap
}

func IsShareExpired(share *models.Share) bool {
	expirationDate, err := time.Parse(time.RFC3339, share.ExpirationDate)
	if err != nil {
		fmt.Println("Error parsing expiration date:", err)
		return false
	}
	return time.Now().After(expirationDate)
}
