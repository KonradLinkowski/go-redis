package main

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"log"

	"github.com/disintegration/imaging"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"

	"github.com/KonradLinkowski/go-redis/shared"
)

func main() {
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr: shared.RedisAddr,
	})

	minioClient, err := shared.InitMinioClient(ctx)
	if err != nil {
		log.Fatalln("Failed to create MinIO client:", err)
	}

	lastID := "0"

	for {
		streams, err := rdb.XRead(ctx, &redis.XReadArgs{
			Streams: []string{shared.RedisStreamName, lastID},
			Block:   0,
		}).Result()

		if err != nil {
			log.Println(err)
			continue
		}

		for _, stream := range streams {
			for _, msg := range stream.Messages {
				err := handleMessage(ctx, rdb, minioClient, msg)
				if err != nil {
					log.Println("Error processing message:", err)
				}
				lastID = msg.ID

			}
		}
	}
}

func handleMessage(ctx context.Context, rdb *redis.Client, minioClient *minio.Client, msg redis.XMessage) error {
	log.Println("Processing message with id:", msg.ID)

	imageName := msg.Values[shared.RedisStreamValueKey].(string)
	object, err := minioClient.GetObject(ctx, shared.MinioUploadedBucketName, imageName, minio.GetObjectOptions{})
	if err != nil {
		return err
	}
	imgBytes, err := io.ReadAll(object)
	if err != nil {
		return err
	}

	out, err := processImage(imgBytes)
	if err != nil {
		return err
	}

	objectName := shared.CreateFileName(imageName)
	_, err = minioClient.PutObject(ctx, shared.MinioProcessedBucketName, objectName, bytes.NewReader(out), int64(len(out)), minio.PutObjectOptions{})
	if err != nil {
		return err
	}

	_, err = rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: shared.RedisProcessedStreamName,
		Values: map[string]interface{}{
			shared.RedisStreamValueKey: objectName,
		},
		MaxLen: 100,
		Approx: true,
	}).Result()

	return err
}

func processImage(data []byte) ([]byte, error) {
	img, format, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	gray := imaging.Grayscale(img)

	var out bytes.Buffer
	switch format {
	case "png":
		err = png.Encode(&out, gray)
	case "jpeg", "jpg":
		err = jpeg.Encode(&out, gray, nil)
	case "gif":
		err = gif.Encode(&out, gray, nil)
	default:
		return nil, fmt.Errorf("unsupported image format: %s", format)
	}

	if err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}
