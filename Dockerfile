FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

ARG SERVICE_NAME
COPY ./${SERVICE_NAME} ./${SERVICE_NAME}

RUN go build -o /main ./${SERVICE_NAME}/main.go

FROM alpine:latest
WORKDIR /root
COPY --from=builder /main .

CMD ["./main"]
