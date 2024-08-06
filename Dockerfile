# Use an official Golang runtime as the base image for the build stage
FROM golang:1.21 AS build

# Set the working directory in the container
WORKDIR /go/src/app

# Copy the Go application source code into the container
COPY . .

# Build the Go application
RUN go build -ldflags="-w -s" -o /golang-kanister-app

# Use a minimal base image for the final Docker image
FROM debian:latest

# Expose the port that your Golang application listens on
EXPOSE 8080

# Copy the Go application binary from the build stage into the final image
COPY --from=build /golang-kanister-app /golang-kanister-app

# Start your Golang application
CMD ["/golang-kanister-app"]
