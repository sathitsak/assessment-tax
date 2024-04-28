# First stage: build the executable.
FROM golang:1.22-alpine as builder

# Set the Current Working Directory inside the container.
WORKDIR /app

# Copy go mod and sum files.
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files do not change.
RUN go mod download

# Copy the source code into the container.
COPY . .

# Build the Go app. Disable CGO and compile for Linux.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o main .

# Second stage: use the alpine base image.
FROM alpine:latest  

# Add ca-certificates in case you need HTTPS.
RUN apk --no-cache add ca-certificates

# Set the working directory in the container.
WORKDIR /root/

# Copy the pre-built binary file from the previous stage.
COPY --from=builder /app/main .

# Expose port 8080.
EXPOSE 8080

# Command to run the executable.
CMD ["./main"]
