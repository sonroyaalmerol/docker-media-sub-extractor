# Use a Golang and FFmpeg base image
FROM jrottenberg/ffmpeg:latest as ffmpeg

# Create a builder stage
FROM golang:latest as builder

# Set the working directory inside the builder stage
WORKDIR /app

# Copy the Go application source code into the builder stage
COPY . .

# Build the Go application
RUN go build -o extract_subtitles .

# Create a smaller final image
FROM alpine:latest

# Copy only the necessary files from the builder and FFmpeg images
COPY --from=ffmpeg /usr/local /usr/local
COPY --from=builder /app/extract_subtitles /app/extract_subtitles

# Set the working directory inside the final image
WORKDIR /app

# Set the entry point for the container
CMD ["./extract_subtitles"]
