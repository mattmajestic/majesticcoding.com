# Dockerfile
FROM golang:1.18-alpine

RUN apk add --no-cache git gcc musl-dev ffmpeg

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod tidy
RUN go mod download

COPY . .

RUN go build -o main main.go

EXPOSE 8080 1935

CMD ["./main"]
