package transformfile

import "mime/multipart"

type FileTransformer interface {
	Transform(input multipart.File) ([]byte, error)
}
