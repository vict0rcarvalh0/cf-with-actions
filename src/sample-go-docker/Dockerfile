# Dockerfile
FROM golang:alpine

WORKDIR /app

COPY main.go main_test.go ./

RUN go mod init myapp
RUN go mod tidy
RUN go build -o main main.go

EXPOSE 8080

CMD ["go", "run", "main.go"]
