FROM golang:1.17-alpine
RUN apk add --no-cache git
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod tidy
RUN go mod download
COPY . .
RUN go build -o main main.go
EXPOSE 8080
CMD ["./main"]