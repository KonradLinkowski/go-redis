package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/KonradLinkowski/go-redis/shared"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
)

func main() {
	fmt.Println("Starting API server on :8000")
	ctx := context.Background()

	minioClient, err := shared.InitMinioClient(ctx)
	if err != nil {
		log.Fatal(err)
	}

	rdb := InitRedisClient()

	hub := NewHub()
	go hub.Run()
	go SubscribeToRedis(ctx, rdb, hub, func(imageName string) {
		log.Println("Received processed image name from Redis:", imageName)
		bytes, err := DownloadFromMinio(ctx, minioClient, imageName)
		if err != nil {
			log.Fatal("Error getting object from MinIO:", err)
		}
		encoded := base64.StdEncoding.EncodeToString(bytes)
		hub.broadcast <- encoded
	})

	http.HandleFunc("/events", createSSEHandler(hub))
	http.HandleFunc("/upload", createUploadHandler(ctx, rdb, minioClient))

	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/", fs)

	http.ListenAndServe(":8000", nil)
}

func createSSEHandler(hub *Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		clientChan := make(chan string)
		hub.register <- clientChan

		defer func() {
			hub.unregister <- clientChan
		}()

		notify := r.Context().Done()

		flusher, _ := w.(http.Flusher)

		log.Println("Client connected to SSE")

		for {
			select {
			case msg := <-clientChan:
				fmt.Fprintf(w, "data: %s\n\n", msg)
				flusher.Flush()
			case <-notify:
				return
			}
		}
	}
}

func createUploadHandler(ctx context.Context, rdb *redis.Client, minioClient *minio.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		file, header, err := r.FormFile("image")
		if err != nil {
			http.Error(w, "Failed to get file", http.StatusBadRequest)
			log.Fatalln("Failed to upload to Minio:", err.Error())
			return
		}
		defer file.Close()

		buf := new(bytes.Buffer)
		buf.ReadFrom(file)
		imageData := buf.Bytes()

		objectName := fmt.Sprintf("%d%s", time.Now().UnixNano(), filepath.Ext(header.Filename))
		_, err = UploadToMinio(ctx, minioClient, objectName, imageData)

		if err != nil {
			http.Error(w, "Failed to upload to Minio", http.StatusInternalServerError)
			log.Fatalln("Failed to upload to Minio:", err.Error())
			return
		}

		err = PushToRedis(ctx, rdb, objectName)
		if err != nil {
			http.Error(w, "Failed to add message to stream", http.StatusInternalServerError)
			log.Fatalln("Failed to add message to stream:", err.Error())
			return
		}

		w.Write([]byte("File uploaded successfully"))
	}
}
