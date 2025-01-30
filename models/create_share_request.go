package models

import "mime/multipart"

type CreateShareRequest struct {
	// Name        string                  `validate:"required"`
	// Description string                  `validate:"required"`
	Files []*multipart.FileHeader `validate:"required"`
}
