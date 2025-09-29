package providers

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
)

type S3Provider struct {
	session *session.Session
	service *s3.S3
	bucket  string
	region  string
}

type IS3Provider interface {
	UploadImage(imageData []byte, filename string, contentType string, userID string) (string, error)
	DeleteImage(imageURL string) error
	GetImageURL(key string) string
}

func NewS3Provider() IS3Provider {
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "ap-southeast-1" // Default region
	}

	bucket := os.Getenv("AWS_BUCKET_NAME")
	if bucket == "" {
		panic("AWS_BUCKET_NAME environment variable is required")
	}

	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	if accessKey == "" || secretKey == "" {
		panic("AWS credentials (AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY) are required")
	}

	// Create AWS session with credentials
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
	})

	if err != nil {
		panic(fmt.Sprintf("Failed to create AWS session: %v", err))
	}

	return &S3Provider{
		session: sess,
		service: s3.New(sess),
		bucket:  bucket,
		region:  region,
	}
}

func (s *S3Provider) UploadImage(imageData []byte, filename string, contentType string, userID string) (string, error) {
	// Generate unique key for the image
	key := s.generateImageKey(filename, userID)

	// Upload to S3
	_, err := s.service.PutObject(&s3.PutObjectInput{
		Bucket:        aws.String(s.bucket),
		Key:           aws.String(key),
		Body:          bytes.NewReader(imageData),
		ContentType:   aws.String(contentType),
		ContentLength: aws.Int64(int64(len(imageData))),
		ACL:           aws.String("public-read"), // Make image publicly accessible
		CacheControl:  aws.String("max-age=31536000"), // Cache for 1 year
	})

	if err != nil {
		return "", fmt.Errorf("failed to upload image to S3: %v", err)
	}

	// Return the public URL
	url := s.GetImageURL(key)
	return url, nil
}

func (s *S3Provider) DeleteImage(imageURL string) error {
	// Extract key from URL
	key := s.extractKeyFromURL(imageURL)
	if key == "" {
		return fmt.Errorf("invalid image URL: %s", imageURL)
	}

	// Delete from S3
	_, err := s.service.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		return fmt.Errorf("failed to delete image from S3: %v", err)
	}

	return nil
}

func (s *S3Provider) GetImageURL(key string) string {
	// Return public URL for S3 object
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.bucket, s.region, key)
}

func (s *S3Provider) generateImageKey(filename string, userID string) string {
	// Get file extension
	ext := filepath.Ext(filename)

	// Generate UUID for uniqueness
	uniqueID := uuid.New().String()

	// Create timestamp for organization
	timestamp := time.Now().Format("2006/01/02")

	// Generate key: profile-images/YYYY/MM/DD/userID/uuid.ext
	key := fmt.Sprintf("profile-images/%s/%s/%s%s", timestamp, userID, uniqueID, ext)

	return key
}

func (s *S3Provider) extractKeyFromURL(imageURL string) string {
	// Expected URL format: https://bucket.s3.region.amazonaws.com/key
	baseURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/", s.bucket, s.region)

	if strings.HasPrefix(imageURL, baseURL) {
		return strings.TrimPrefix(imageURL, baseURL)
	}

	// Alternative format: https://s3.region.amazonaws.com/bucket/key
	altBaseURL := fmt.Sprintf("https://s3.%s.amazonaws.com/%s/", s.region, s.bucket)
	if strings.HasPrefix(imageURL, altBaseURL) {
		return strings.TrimPrefix(imageURL, altBaseURL)
	}

	return ""
}

// Helper function to check if S3 connection is working
func (s *S3Provider) TestConnection() error {
	// Try to list bucket contents (limit to 1 object)
	_, err := s.service.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket:  aws.String(s.bucket),
		MaxKeys: aws.Int64(1),
	})

	if err != nil {
		return fmt.Errorf("S3 connection test failed: %v", err)
	}

	return nil
}