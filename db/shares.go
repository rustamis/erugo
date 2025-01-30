package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/DeanWard/erugo/models"
)

func ShareByLongId(db *sql.DB, longId string) *models.Share {
	row := db.QueryRow("SELECT * FROM shares WHERE long_id = ?", longId)
	var share models.Share
	var filesJson string
	err := row.Scan(&share.Id, &share.FilePath, &share.ExpirationDate, &share.LongId, &share.NumFiles, &share.TotalSize, &filesJson, &share.UserId)
	if err == nil {
		err = json.Unmarshal([]byte(filesJson), &share.Files)
	}
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return &share
}

func ShareCreate(db *sql.DB, share *models.Share) (*models.Share, error) {
	//convert the files slice to json string
	filesJson, err := json.Marshal(share.Files)
	if err != nil {
		log.Printf("Error marshalling files: %v", err)
		return nil, err
	}
	_, err = db.Exec(
		"INSERT INTO shares (file_path, expiration_date, long_id, num_files, total_size, files, user_id) VALUES (?, ?, ?, ?, ?, ?, ?)",
		share.FilePath, share.ExpirationDate, share.LongId, share.NumFiles, share.TotalSize, filesJson, share.UserId,
	)
	if err != nil {
		log.Printf("Error creating share: %v", err)
		return nil, err
	}
	//get the share we just created using LAST_INSERT_ROWID()
	row := db.QueryRow("SELECT * FROM shares WHERE id = LAST_INSERT_ROWID()")
	err = row.Scan(&share.Id, &share.FilePath, &share.ExpirationDate, &share.LongId, &share.NumFiles, &share.TotalSize, &filesJson, &share.UserId)
	if err != nil {
		log.Printf("Error getting share: %v", err)
		return nil, err
	}
	return share, nil
}
