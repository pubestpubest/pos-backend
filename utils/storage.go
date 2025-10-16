package utils

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/pkg/errors"
)

// UploadImageToMinio uploads an image to MinIO and returns the public URL
func UploadImageToMinio(client *minio.Client, file *multipart.FileHeader) (string, error) {
	// Validate file type
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}

	if !allowedExtensions[ext] {
		return "", errors.New("Invalid file type. Only image files are allowed (jpg, jpeg, png, gif, webp)")
	}

	// Validate file size (max 5MB)
	maxSize := int64(5 * 1024 * 1024) // 5MB
	if file.Size > maxSize {
		return "", errors.New("File size exceeds maximum allowed size of 5MB")
	}

	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		return "", errors.Wrap(err, "[UploadImageToMinio]: Failed to open uploaded file")
	}
	defer src.Close()

	// Generate unique filename
	filename := fmt.Sprintf("%s-%d%s", uuid.New().String(), time.Now().Unix(), ext)

	// Get bucket name from environment
	bucketName := os.Getenv("MINIO_BUCKET_NAME")
	if bucketName == "" {
		bucketName = "menu-images"
	}

	// Determine content type
	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// Upload to MinIO
	ctx := context.Background()
	_, err = client.PutObject(ctx, bucketName, filename, src, file.Size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", errors.Wrap(err, "[UploadImageToMinio]: Failed to upload file to MinIO")
	}

	// Construct the public URL
	minioEndpoint := os.Getenv("MINIO_ENDPOINT")
	useSSL := os.Getenv("MINIO_USE_SSL") == "true"
	protocol := "http"
	if useSSL {
		protocol = "https"
	}

	publicURL := fmt.Sprintf("%s://%s/%s/%s", protocol, minioEndpoint, bucketName, filename)

	return publicURL, nil
}

// DeleteImageFromMinio deletes an image from MinIO
func DeleteImageFromMinio(client *minio.Client, imageURL string) error {
	if imageURL == "" {
		return nil // Nothing to delete
	}

	// Extract filename from URL
	parts := strings.Split(imageURL, "/")
	if len(parts) < 2 {
		return errors.New("Invalid image URL")
	}
	filename := parts[len(parts)-1]

	// Get bucket name from environment
	bucketName := os.Getenv("MINIO_BUCKET_NAME")
	if bucketName == "" {
		bucketName = "menu-images"
	}

	// Delete from MinIO
	ctx := context.Background()
	err := client.RemoveObject(ctx, bucketName, filename, minio.RemoveObjectOptions{})
	if err != nil {
		return errors.Wrap(err, "[DeleteImageFromMinio]: Failed to delete file from MinIO")
	}

	return nil
}
