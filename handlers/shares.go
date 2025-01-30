package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/DeanWard/erugo/config"
	"github.com/DeanWard/erugo/db"
	"github.com/DeanWard/erugo/middleware"
	"github.com/DeanWard/erugo/models"
	"github.com/DeanWard/erugo/responses/file_response"
	"github.com/DeanWard/erugo/responses/json_response"
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
		json_response.New(json_response.ErrorStatus, "Validation error", payload, http.StatusBadRequest).Send(w)
		//log the error
		utils.Log(fmt.Sprintf("Validation error: %v", payload), utils.ColorRed)
		return
	}

	maxShareSize := config.GetMaxShareSize()
	utils.Log(fmt.Sprintf("Max share size: %d", maxShareSize), utils.ColorGreen)
	if maxShareSize == 0 {
		json_response.New(json_response.ErrorStatus, "Max share size not set", nil, http.StatusInternalServerError).Send(w)
		utils.Log("Max share size not set", utils.ColorRed)
		return
	}

	//check if the total size of the files in the request is greater than the max share size
	totalRequestSize := r.ContentLength
	if totalRequestSize > maxShareSize {
		json_response.New(json_response.ErrorStatus, "Total size of files is greater than the max share size", nil, http.StatusBadRequest).Send(w)
		utils.Log("Total size of files is greater than the max share size", utils.ColorRed)
		return
	}

	longId := generateUniqueLongId(database)
	folderPath := store.StoreUploadedFiles(config.AppConfig.BaseStoragePath, files, longId)
	if folderPath == "" {
		json_response.New(json_response.ErrorStatus, "There was an error saving your files.", nil, http.StatusBadRequest).Send(w)
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
	savedShare, err := db.CreateShare(database, share)
	if err != nil {
		json_response.New(json_response.ErrorStatus, "Failed to create share", nil, http.StatusInternalServerError).Send(w)
		utils.Log("Failed to create share", utils.ColorRed)
		return
	}
	payload := map[string]interface{}{
		"share": savedShare.ToShareResponse(),
	}
	json_response.New(json_response.SuccessStatus, "Share created", payload, http.StatusOK).Send(w)
}

func handleGetShare(database *sql.DB, w http.ResponseWriter, longId string) {
	share := db.GetShareByLongId(database, longId)
	if share == nil {
		json_response.New(json_response.ErrorStatus, "Share not found", nil, http.StatusNotFound).Send(w)
		return
	}

	if IsShareExpired(share) {
		json_response.New(json_response.ErrorStatus, "Share has expired", nil, http.StatusGone).Send(w)
		return
	}
	payload := map[string]interface{}{
		"share": share.ToShareResponse(),
	}
	json_response.New(json_response.SuccessStatus, "Share found", payload, http.StatusOK).Send(w)
}

func generateUniqueLongId(database *sql.DB) string {
	haikunator := haikunator.New(time.Now().UnixNano())
	longId := haikunator.Haikunate() + "-" + haikunator.Haikunate()
	for db.GetShareByLongId(database, longId) != nil {
		longId = haikunator.Haikunate() + "-" + haikunator.Haikunate()
	}
	return longId
}

func handleDownloadShare(database *sql.DB, w http.ResponseWriter, r *http.Request, longId string) {
	share := db.GetShareByLongId(database, longId)
	if share == nil {
		json_response.New(json_response.ErrorStatus, "Share not found", nil, http.StatusNotFound).Send(w)
		utils.Log("Share not found", utils.ColorRed)
		return
	}

	if IsShareExpired(share) {
		json_response.New(json_response.ErrorStatus, "Share has expired", nil, http.StatusGone).Send(w)
		utils.Log("Share has expired", utils.ColorRed)
		return
	}

	if _, err := os.Stat(share.FilePath); os.IsNotExist(err) {
		json_response.New(json_response.ErrorStatus, "The file you are looking for has been deleted or expired.", nil, http.StatusNotFound).Send(w)
		utils.Log("The file you are looking for has been deleted or expired.", utils.ColorRed)
		return
	}

	userEmail := r.FormValue("user_email")
	db.CreateShareAccessLogEntry(database, &models.ShareAccessLogEntry{
		ShareId:    share.Id,
		UserEmail:  userEmail,
		UserIp:     r.RemoteAddr,
		UserAgent:  r.UserAgent(),
		AccessDate: time.Now().Format(time.RFC3339),
	})

	//now we know the folder exists, let's put them in a zip file and serve that
	zipFilePath, err := store.CreateZipFile(share.FilePath)
	if err != nil {
		json_response.New(json_response.ErrorStatus, "Failed to create zip file: "+err.Error(), nil, http.StatusInternalServerError).Send(w)
		utils.Log("Failed to create zip file: "+err.Error(), utils.ColorRed)
		return
	}
	fileName := fmt.Sprintf("%s.zip", share.LongId)
	file_response.New(zipFilePath, fileName).Send(w, r)
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
