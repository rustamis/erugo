package db

import (
	"database/sql"
	"log"
	"os"

	"github.com/DeanWard/erugo/config"
)

func Connect() *sql.DB {
	log.Println("Connecting to database ./erugo.db")

	//does the file exist?
	if _, err := os.Stat(config.AppConfig.DatabaseFilePath); os.IsNotExist(err) {
		log.Println("Database file does not exist, creating it")
	}

	db, err := sql.Open("sqlite", config.AppConfig.DatabaseFilePath)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	log.Println("Database connected")
	return db
}
