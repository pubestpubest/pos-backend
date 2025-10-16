package database

import (
	"context"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var MinioClient *minio.Client

func ConnectMinio() error {
	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKeyID := os.Getenv("MINIO_ROOT_USER")
	secretAccessKey := os.Getenv("MINIO_ROOT_PASSWORD")
	useSSL := os.Getenv("MINIO_USE_SSL") == "true"

	if endpoint == "" || accessKeyID == "" || secretAccessKey == "" {
		return errors.New("Missing MinIO configuration")
	}

	// Initialize MinIO client
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return errors.Wrap(err, "[ConnectMinio]: Failed to initialize MinIO client")
	}

	MinioClient = minioClient
	log.Info("[ConnectMinio]: Successfully connected to MinIO")

	// Create default bucket if it doesn't exist
	bucketName := os.Getenv("MINIO_BUCKET_NAME")
	if bucketName == "" {
		bucketName = "menu-images"
	}

	ctx := context.Background()
	exists, err := minioClient.BucketExists(ctx, bucketName)
	if err != nil {
		return errors.Wrap(err, "[ConnectMinio]: Failed to check bucket existence")
	}

	if !exists {
		err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return errors.Wrap(err, "[ConnectMinio]: Failed to create bucket")
		}
		log.Infof("[ConnectMinio]: Successfully created bucket '%s'", bucketName)

		// Set bucket policy to public read for images
		policy := `{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Effect": "Allow",
					"Principal": {"AWS": ["*"]},
					"Action": ["s3:GetObject"],
					"Resource": ["arn:aws:s3:::` + bucketName + `/*"]
				}
			]
		}`
		err = minioClient.SetBucketPolicy(ctx, bucketName, policy)
		if err != nil {
			log.Warnf("[ConnectMinio]: Failed to set bucket policy: %v", err)
		}
	} else {
		log.Infof("[ConnectMinio]: Bucket '%s' already exists", bucketName)
	}

	return nil
}
