FROM golang:1.17-alpine
WORKDIR /app
COPY . .
COPY . /app
RUN go build -o main main.go
EXPOSE 8000
CMD ["./main"]