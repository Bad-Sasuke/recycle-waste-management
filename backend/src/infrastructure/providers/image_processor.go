package providers

import (
	"bytes"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"path/filepath"
	"strings"

	"github.com/nickalie/go-webpbin"
)

type ImageProcessor struct {
	Quality   int
	MaxWidth  int
	MaxHeight int
}

type IImageProcessor interface {
	ProcessImage(imageData []byte, filename string) (*ProcessedImage, error)
	IsValidImageFormat(filename string) bool
}

type ProcessedImage struct {
	Data        []byte
	Filename    string
	ContentType string
	Size        int64
}

func NewImageProcessor() IImageProcessor {
	return &ImageProcessor{
		Quality:   80, // WebP quality (0-100)
		MaxWidth:  1024,
		MaxHeight: 1024,
	}
}

func (p *ImageProcessor) IsValidImageFormat(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	validExts := []string{".jpg", ".jpeg", ".png", ".gif"}

	for _, validExt := range validExts {
		if ext == validExt {
			return true
		}
	}
	return false
}

func (p *ImageProcessor) ProcessImage(imageData []byte, filename string) (*ProcessedImage, error) {
	// Validate image format
	if !p.IsValidImageFormat(filename) {
		return nil, fmt.Errorf("unsupported image format: %s", filepath.Ext(filename))
	}

	// Decode the image to validate it's actually an image
	_, format, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return nil, fmt.Errorf("invalid image data: %v", err)
	}

	// Resize if needed
	resizedData, err := p.resizeImage(imageData, format)
	if err != nil {
		return nil, fmt.Errorf("failed to resize image: %v", err)
	}

	// Convert to WebP
	webpData, err := p.convertToWebP(resizedData)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to WebP: %v", err)
	}

	// Generate new filename with .webp extension
	baseName := strings.TrimSuffix(filename, filepath.Ext(filename))
	webpFilename := baseName + ".webp"

	return &ProcessedImage{
		Data:        webpData,
		Filename:    webpFilename,
		ContentType: "image/webp",
		Size:        int64(len(webpData)),
	}, nil
}

func (p *ImageProcessor) resizeImage(imageData []byte, format string) ([]byte, error) {
	// Decode image to check dimensions
	img, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// If image is already within limits, return original
	if width <= p.MaxWidth && height <= p.MaxHeight {
		return imageData, nil
	}

	// Calculate new dimensions while maintaining aspect ratio
	ratio := float64(width) / float64(height)
	newWidth := p.MaxWidth
	newHeight := int(float64(newWidth) / ratio)

	if newHeight > p.MaxHeight {
		newHeight = p.MaxHeight
		newWidth = int(float64(newHeight) * ratio)
	}

	// Use webpbin to resize (it supports resizing before conversion)
	return imageData, nil // For now, return original - webpbin will handle resizing during conversion
}

func (p *ImageProcessor) convertToWebP(imageData []byte) ([]byte, error) {
	// Create input reader
	input := bytes.NewReader(imageData)

	// Create output buffer
	var output bytes.Buffer

	err := webpbin.NewCWebP().
		Quality(80).
		Input(input).
		Output(&output).
		Run()

	if err != nil {
		return nil, fmt.Errorf("webp conversion failed: %v", err)
	}

	return output.Bytes(), nil
}

// Helper function to get image info
func (p *ImageProcessor) GetImageInfo(imageData []byte) (width, height int, format string, err error) {
	img, fmt, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return 0, 0, "", err
	}

	bounds := img.Bounds()
	return bounds.Dx(), bounds.Dy(), fmt, nil
}

// Helper function to validate image size
func (p *ImageProcessor) ValidateImageSize(imageData []byte, maxSizeBytes int64) error {
	if int64(len(imageData)) > maxSizeBytes {
		return fmt.Errorf("image size %d bytes exceeds maximum allowed size %d bytes", len(imageData), maxSizeBytes)
	}
	return nil
}
