package transformfile

import (
	"bytes"
	"fmt"
	"image"
	"io"
	"mime/multipart"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/chai2010/webp"
	"github.com/nfnt/resize"
)

type ImageTransformer struct {
	Width  uint
	Height uint
}

func NewImageTransformer(width uint, height uint) FileTransformer {
	return &ImageTransformer{
		width,
		height,
	}
}

func (transformer *ImageTransformer) Transform(
	input multipart.File,
) ([]byte, error) {
	// Make sure file seeker at start
	if seeker, ok := input.(io.Seeker); ok {
		seeker.Seek(0, io.SeekStart)
	}

	// Decode image
	img, _, err := image.Decode(input)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	// Resize
	resized := resize.Resize(
		transformer.Width,
		transformer.Height,
		img,
		resize.Lanczos3,
	)

	// Convert to webp
	var buffer bytes.Buffer
	err = webp.Encode(&buffer, resized, &webp.Options{Quality: 90})
	if err != nil {
		return nil, fmt.Errorf("failed to encode image as webp: %w", err)
	}

	return buffer.Bytes(), nil
}
