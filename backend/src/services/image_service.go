package services

import (
	"fmt"
	"io"
	"mime/multipart"

	"recycle-waste-management-backend/src/infrastructure/providers"
)

type imageService struct {
	ImageProcessor providers.IImageProcessor
	S3Provider     providers.IS3Provider
}

type IImageService interface {
	ProcessAndUploadProfileImage(file *multipart.FileHeader, userID string) (string, error)
	DeleteProfileImage(imageURL string) error
	ValidateImageFile(file *multipart.FileHeader) error
}

func NewImageService() IImageService {
	return &imageService{
		ImageProcessor: providers.NewImageProcessor(),
		S3Provider:     providers.NewS3Provider(),
	}
}

func (s *imageService) ProcessAndUploadProfileImage(file *multipart.FileHeader, userID string) (string, error) {
	// Validate file
	if err := s.ValidateImageFile(file); err != nil {
		return "", err
	}

	// Open and read file
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %v", err)
	}
	defer src.Close()

	// Read file data
	fileData, err := io.ReadAll(src)
	if err != nil {
		return "", fmt.Errorf("failed to read file data: %v", err)
	}

	// Process image (resize and convert to WebP)
	processedImage, err := s.ImageProcessor.ProcessImage(fileData, file.Filename)
	if err != nil {
		return "", fmt.Errorf("failed to process image: %v", err)
	}

	// Upload to S3
	imageURL, err := s.S3Provider.UploadImage(
		processedImage.Data,
		processedImage.Filename,
		processedImage.ContentType,
		userID,
	)
	if err != nil {
		return "", fmt.Errorf("failed to upload image to S3: %v", err)
	}

	return imageURL, nil
}

func (s *imageService) DeleteProfileImage(imageURL string) error {
	if imageURL == "" {
		return nil // Nothing to delete
	}

	err := s.S3Provider.DeleteImage(imageURL)
	if err != nil {
		return fmt.Errorf("failed to delete image from S3: %v", err)
	}

	return nil
}

func (s *imageService) ValidateImageFile(file *multipart.FileHeader) error {
	// Check if file exists
	if file == nil {
		return fmt.Errorf("no file uploaded")
	}

	// Validate file size (max 10MB)
	const maxSize = 10 * 1024 * 1024 // 10MB
	if file.Size > maxSize {
		return fmt.Errorf("file size %d bytes exceeds maximum allowed size %d bytes", file.Size, maxSize)
	}

	// Validate file format by extension
	if !s.ImageProcessor.IsValidImageFormat(file.Filename) {
		return fmt.Errorf("unsupported image format. Supported formats: JPG, JPEG, PNG, GIF")
	}

	return nil
}
