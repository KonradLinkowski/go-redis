package main

import (
	"bytes"
	"context"

	"github.com/minio/minio-go/v7"

	"github.com/KonradLinkowski/go-redis/shared"
)

func UploadToMinio(ctx context.Context, minioClient *minio.Client, objectName string, data []byte) (minio.UploadInfo, error) {
	bucketName := shared.MinioUploadedBucketName
	return minioClient.PutObject(ctx, bucketName, objectName, bytes.NewReader(data), int64(len(data)), minio.PutObjectOptions{})
}

func DownloadFromMinio(ctx context.Context, minioClient *minio.Client, objectName string) ([]byte, error) {
	object, err := minioClient.GetObject(ctx, shared.MinioProcessedBucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer object.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(object)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
