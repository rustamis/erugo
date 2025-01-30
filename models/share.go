package models

type Share struct {
	Id             int      `json:"id"`
	FilePath       string   `json:"file_path"`
	ExpirationDate string   `json:"expiration_date"`
	LongId         string   `json:"long_id"`
	NumFiles       int      `json:"num_files"`
	TotalSize      int64    `json:"total_size"`
	Files          []string `json:"files"`
	UserId         int      `json:"user_id"`
}

type ShareResponse struct {
	Id             int      `json:"id"`
	ExpirationDate string   `json:"expiration_date"`
	NumFiles       int      `json:"num_files"`
	TotalSize      int64    `json:"total_size"`
	LongId         string   `json:"long_id"`
	Files          []string `json:"files"`
}

func (s *Share) ToShareResponse() *ShareResponse {
	return &ShareResponse{
		Id:             s.Id,
		ExpirationDate: s.ExpirationDate,
		NumFiles:       s.NumFiles,
		TotalSize:      s.TotalSize,
		LongId:         s.LongId,
		Files:          s.Files,
	}
}
