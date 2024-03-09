# Start from the official Golang image
FROM golang:1.22.1-alpine3.19 as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the source code into the container
COPY channel_owner ./channel_owner
COPY shared ./shared

# Build the Go app
RUN cd ./channel_owner && GOARCH=amd64 go build -o myapp ./cmd

# Start a new stage from scratch
FROM alpine:latest

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/channel_owner/myapp .
COPY ./channel_owner/.env .


# Command to run the executable
CMD ["./myapp"]


#
## Start from the official Golang image
#FROM golang:1.22.1-alpine3.19
#
## Set the Current Working Directory inside the container
#WORKDIR /app
#
## Copy the source code into the container
#COPY channel_owner ./channel_owner
#COPY shared ./shared
#
## Build the Go app
#RUN cd ./channel_owner && GOARCH=amd64 go build -o myapp ./cmd
#
#CMD ["/app/channel_owner/myapp"]
#
