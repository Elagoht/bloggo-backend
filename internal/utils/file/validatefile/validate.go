package validatefile

import "mime/multipart"

type FileValidator interface {
	Validate(file multipart.File, header *multipart.FileHeader) error
}
