package file_response

import (
	"fmt"
	"net/http"
)

type FileResponse struct {
	FilePath string `json:"file_path"`
	FileName string `json:"file_name"`
}

func New(filePath string, fileName string) *FileResponse {
	return &FileResponse{
		FilePath: filePath,
		FileName: fileName,
	}
}

func (f *FileResponse) Send(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", f.FileName))
	http.ServeFile(w, r, f.FilePath)
}
