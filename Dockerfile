# 1. Start from the official Golang 1.20 image
FROM golang:1.20
# 2. Set the working directory inside the container
WORKDIR /app
# 3. Copy go.mod and go.sum files
COPY go.mod go.sum ./
# 4. Download dependencies
RUN go mod download
# 5. Copy the rest of the application's source code
COPY . .
# 6. Build the application
RUN go build -o main .
# 7. Expose port 8080
EXPOSE 8080
# 8. Set the startup command for the container
CMD ["./main"]