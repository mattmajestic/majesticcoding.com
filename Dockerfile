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
HEALTHCHECK --interval=30s --timeout=30s --start-period=5s --retries=3 CMD [ "wget", "-q", "-O", "-", "http://localhost:8000/health" ] || exit 1
CMD ["./main"]