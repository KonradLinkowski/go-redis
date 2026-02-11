# Go Photo Processor (Grayscale)

A high-performance, asynchronous image processing service built with Go. This application allows users to upload images, which are then processed into black-and-white versions using a distributed worker pattern.

## Features

* Asynchronous Processing: Uses Redis Streams for reliable task queuing.
* Object Storage: Leverages MinIO (S3-compatible) for storing raw and processed images.
* Real-time Updates: Server-Sent Events (SSE) notify the frontend as soon as the image is ready.
* Scalable Architecture: Decoupled API and Worker services.

## Architecture

1. API: Receives the image, uploads the original to MinIO, and pushes a message to a Redis Stream.
2. Worker: Listens to the Redis Stream, downloads the image, converts it to grayscale, uploads the result back to MinIO, and signals completion via another stream.
3. SSE: The frontend receives a real-time event when the processing is finished to display the result.

## Tech Stack

* Backend: Go (Golang)
* Message Broker: Redis (Streams)
* Storage: MinIO
* Updates: Server-Sent Events (SSE)
* Infrastructure: Docker Compose

## Getting Started

### Development
To run the project in development mode:

```bash
docker compose -f compose.dev.yml up
```

Access the application at: http://localhost:8000

### Production
For a stable, optimized environment:

```bash
docker compose up -d
```

### Infrastructure Access
* Web App: http://localhost:8000
* MinIO Console: http://localhost:9001

## API Endpoints

* POST /upload - Upload an image to start processing
* GET /events - SSE endpoint for real-time status updates
