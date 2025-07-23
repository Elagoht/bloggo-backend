package bucket

import "mime/multipart"

type Bucket interface {
	// Upsert file
	Save(file []byte, name string) error
	SaveMultiPart(file multipart.File, name string) error
	// Read file
	Get(name string) ([]byte, error)
	// Destroy data
	Delete(name string) error
	DeleteMatching(pattern string, blacklist string) error
}
