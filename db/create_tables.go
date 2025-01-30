package db

import (
	"database/sql"
	"io/fs"
	"log"
	"sort"
	"time"

	"github.com/DeanWard/erugo/utils"
)

func Migrate(db *sql.DB, migrations fs.FS) {

	utils.Log("Running database migrations...", utils.ColorBlue)

	db.Exec("CREATE TABLE IF NOT EXISTS migrations (filename TEXT, applied_at TEXT)")

	//order the migrations by name
	files, err := fs.ReadDir(migrations, ".")
	if err != nil {
		log.Fatalf("Failed to read migrations directory: %v", err)
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	utils.Log("Found %d migrations to apply", utils.ColorBlue, len(files))

	//execute the migrations
	for _, file := range files {
		utils.Log("Checking migration %s:", utils.ColorBlue, file.Name())
		if checkIfMigrationApplied(db, file.Name()) {
			utils.Log("Already applied.", utils.ColorYellow)
			continue
		}
		utils.Log("Applying.", utils.ColorGreen)
		applyMigration(db, migrations, file.Name())
		logMigration(db, file.Name())
	}

	utils.Log("Database migrations complete.", utils.ColorGreen)

}

func checkIfMigrationApplied(db *sql.DB, filename string) bool {
	var appliedAt string
	err := db.QueryRow("SELECT applied_at FROM migrations WHERE filename = ?", filename).Scan(&appliedAt)
	if err != nil {
		return false
	}
	return appliedAt != ""
}

func applyMigration(db *sql.DB, migration fs.FS, filename string) {
	sql, err := fs.ReadFile(migration, filename)
	if err != nil {
		utils.Log("Failed to read migration file: %v", utils.ColorRed, err)
	}
	db.Exec(string(sql))
}

func logMigration(db *sql.DB, filename string) {
	timeNow := time.Now().Format(time.RFC3339)
	utils.Log("Logging migration as applied at %s", utils.ColorGreen, timeNow)
	db.Exec("INSERT INTO migrations (filename, applied_at) VALUES (?, ?)", filename, timeNow)
}
