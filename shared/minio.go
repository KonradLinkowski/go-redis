package shared

import (
	"context"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func InitMinioClient(ctx context.Context) (*minio.Client, error) {
	minioClient, err := minio.New(MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(MinioAccessKey, MinioSecretKey, ""),
		Secure: false,
	})

	if err == nil {
		err = createBucketIfNotExists(ctx, minioClient, MinioUploadedBucketName)
		if err != nil {
			log.Fatal(err)
		}
		err = createBucketIfNotExists(ctx, minioClient, MinioProcessedBucketName)
		if err != nil {
			log.Fatal(err)
		}
	}

	return minioClient, err
}

func createBucketIfNotExists(ctx context.Context, minioClient *minio.Client, bucketName string) error {
	exists, err := minioClient.BucketExists(ctx, bucketName)
	if err != nil {
		return err
	}

	if !exists {
		return minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	}

	return nil
}
