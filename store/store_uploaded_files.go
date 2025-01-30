package store

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/DeanWard/erugo/config"
)

func StoreUploadedFiles(baseStoragePath string, files []*multipart.FileHeader, longId string) string {
	if len(files) == 0 {
		return ""
	}
	folderPath := filepath.Join(baseStoragePath, "shares", longId)
	os.MkdirAll(folderPath, 0755)
	for _, file := range files {
		filePath := filepath.Join(folderPath, file.Filename)
		f, err := os.Create(filePath)
		if err != nil {
			log.Printf("Failed to save file: %v", err)
			continue
		}
		defer f.Close()
		src, _ := file.Open()
		defer src.Close()
		_, err = io.Copy(f, src)
		if err != nil {
			log.Printf("Failed to write file: %v", err)
		}
	}
	return folderPath
}

// CreateZipFile creates a zip file containing the contents of the given folder
func CreateZipFile(folderPath string) (string, error) {

	//check that the folder exists and isn't above the base storage path
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		return "", fmt.Errorf("folder does not exist: %s", folderPath)
	}
	if !strings.HasPrefix(folderPath, config.AppConfig.BaseStoragePath) {
		return "", fmt.Errorf("folder is above the base storage path: %s / %s", folderPath, config.AppConfig.BaseStoragePath)
	}

	// Name of the zip file to be created
	zipFilePath := folderPath + ".zip"

	//check that the zip file doesn't already exist - if it does just return the path
	if _, err := os.Stat(zipFilePath); !os.IsNotExist(err) {
		return zipFilePath, nil
	}

	// Create the zip file
	zipFile, err := os.Create(zipFilePath)
	if err != nil {
		return "", err
	}
	defer zipFile.Close()

	// Create a new zip writer
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Walk the folder and add files to the zip
	err = filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the root directory itself
		if path == folderPath {
			return nil
		}

		// Create zip header
		relPath, err := filepath.Rel(folderPath, path)
		if err != nil {
			return err
		}
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name = relPath

		// If it's a directory, add it as is
		if info.IsDir() {
			header.Name += "/"
			_, err := zipWriter.CreateHeader(header)
			return err
		}

		// If it's a file, write its contents
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(writer, file)
		return err
	})

	if err != nil {
		return "", err
	}

	return zipFilePath, nil
}

func GetFilesInFolder(folderPath string) []string {
	files, err := os.ReadDir(folderPath)
	if err != nil {
		return []string{}
	}
	//return the names of the files in the folder
	fileNames := []string{}
	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}
	return fileNames
}

func GetTotalSize(folderPath string) int64 {
	files, err := os.ReadDir(folderPath)
	if err != nil {
		return 0
	}
	totalSize := int64(0)
	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			continue
		}
		totalSize += info.Size()
	}
	return totalSize
}
