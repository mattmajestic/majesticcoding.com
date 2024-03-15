FROM golang:1.18-alpine
RUN apk add --no-cache git gcc musl-dev
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod tidy
RUN go mod download
RUN go get gorm.io/driver/sqlite
COPY . .
RUN go build -o main main.go
EXPOSE 8080
CMD ["./main"]