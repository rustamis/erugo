package db

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/DeanWard/erugo/models"
)

func GetShareByLongId(db *sql.DB, longId string) *models.Share {
	row := db.QueryRow("SELECT * FROM shares WHERE long_id = ?", longId)
	var share models.Share
	var filesJson string
	err := row.Scan(&share.Id, &share.FilePath, &share.ExpirationDate, &share.LongId, &share.NumFiles, &share.TotalSize, &filesJson)
	if err == nil {
		err = json.Unmarshal([]byte(filesJson), &share.Files)
	}
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return &share
}
