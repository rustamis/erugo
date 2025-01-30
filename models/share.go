package models

type Share struct {
	Id             int      `json:"id"`
	FilePath       string   `json:"file_path"`
	ExpirationDate string   `json:"expiration_date"`
	LongId         string   `json:"long_id"`
	NumFiles       int      `json:"num_files"`
	TotalSize      int64    `json:"total_size"`
	Files          []string `json:"files"`
}
