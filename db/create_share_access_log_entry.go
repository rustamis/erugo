package db

import (
	"database/sql"

	"github.com/DeanWard/erugo/models"
)

func CreateShareAccessLogEntry(db *sql.DB, shareAccessLogEntry *models.ShareAccessLogEntry) {
	db.Exec("INSERT INTO share_access_logs (share_id, user_email, user_ip, user_agent, access_date) VALUES (?, ?, ?, ?, ?)", shareAccessLogEntry.ShareId, shareAccessLogEntry.UserEmail, shareAccessLogEntry.UserIp, shareAccessLogEntry.UserAgent, shareAccessLogEntry.AccessDate)
}
