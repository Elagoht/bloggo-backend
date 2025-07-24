package transtormfile

import "mime/multipart"

type FileTransformer interface {
	Transform(input multipart.File) ([]byte, error)
}
