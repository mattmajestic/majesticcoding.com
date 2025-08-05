# Dockerfile
FROM golang:1.24-alpine

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod tidy
RUN go mod download

COPY . .

RUN go build -o main ./main.go

EXPOSE 8080 1935

CMD ["./main"]
