# Use a Golang alpine base image
FROM golang:1.21.5-alpine

# Set the working directory inside the builder stage
WORKDIR /app

# Copy the Go application source code into the builder stage
COPY . .

# Install FFmpeg in the builder stage
RUN apk --no-cache add ffmpeg

# Build the Go application
RUN go build -o extract_subtitles .

# Set the entry point for the container
CMD ["./extract_subtitles"]
