package db

import (
	"database/sql"
	"encoding/json"
	"log"

	"github.com/DeanWard/erugo/models"
)

func CreateShare(db *sql.DB, share *models.Share) (*models.Share, error) {
	//convert the files slice to json string
	filesJson, err := json.Marshal(share.Files)
	if err != nil {
		log.Printf("Error marshalling files: %v", err)
		return nil, err
	}
	_, err = db.Exec("INSERT INTO shares (file_path, expiration_date, long_id, num_files, total_size, files) VALUES (?, ?, ?, ?, ?, ?)", share.FilePath, share.ExpirationDate, share.LongId, share.NumFiles, share.TotalSize, filesJson)
	if err != nil {
		log.Printf("Error creating share: %v", err)
		return nil, err
	}
	//get the share we just created using LAST_INSERT_ROWID()
	row := db.QueryRow("SELECT * FROM shares WHERE id = LAST_INSERT_ROWID()")
	err = row.Scan(&share.Id, &share.FilePath, &share.ExpirationDate, &share.LongId, &share.NumFiles, &share.TotalSize, &filesJson)
	if err != nil {
		log.Printf("Error getting share: %v", err)
		return nil, err
	}
	return share, nil
}
