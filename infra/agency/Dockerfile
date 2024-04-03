# Start from the official Golang image
FROM golang:1.22.1-alpine3.19 as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the source code into the container
COPY ../../golang/agency ./agency
COPY ../../golang/shared ./shared

# Build the Go app
RUN cd ./agency && GOARCH=amd64 go build -o myapp ./cmd

# Start a new stage from scratch
FROM alpine:latest

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/agency/myapp .
COPY ../../golang/agency/example.jpg .


# Command to run the executable
CMD ["./myapp"]