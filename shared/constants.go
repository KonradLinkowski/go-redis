package shared

import "os"

const (
	RedisStreamName          = "image_uploaded"
	RedisStreamValueKey      = "image_name"
	RedisProcessedStreamName = "image_processed"
)

var (
	RedisAddr                = os.Getenv("REDIS_ADDR")
	MinioEndpoint            = os.Getenv("MINIO_ENDPOINT")
	MinioUploadedBucketName  = os.Getenv("MINIO_UPLOADED_BUCKET_NAME")
	MinioProcessedBucketName = os.Getenv("MINIO_PROCESSED_BUCKET_NAME")
	MinioAccessKey           = os.Getenv("MINIO_ROOT_USER")
	MinioSecretKey           = os.Getenv("MINIO_ROOT_PASSWORD")
)
