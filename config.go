package main

import (
	"os"
	"strconv"
	"strings"
)

var allowedExtensions []string
var processedFilesPath string
var concurrency int

func initializeConfig() {
	// Read configuration from environment variables or use default values
	allowedExtensions = strings.Split(os.Getenv("ALLOWED_EXTENSIONS"), ",")
	if len(allowedExtensions) == 0 {
		allowedExtensions = []string{".mp4", ".mkv", ".avi", ".wmv"}
	}

	processedFilesPath = os.Getenv("PROCESSED_FILES_PATH")
	if processedFilesPath == "" {
		processedFilesPath = "/app/processed_files.txt"
	}

	concurrencyStr := os.Getenv("CONCURRENCY")
	if concurrencyStr == "" {
		concurrency = 4
	} else {
		concurrency = parseConcurrency(concurrencyStr)
	}
}

func parseConcurrency(str string) int {
	val, err := strconv.Atoi(str)
	if err != nil {
		return 4
	}
	return val
}

func initializeProcessedFiles() string {
	// Create the processed files file if it doesn't exist
	if _, err := os.Stat(processedFilesPath); os.IsNotExist(err) {
		file, err := os.Create(processedFilesPath)
		if err != nil {
			os.Exit(1)
		}
		defer file.Close()
	}
	return processedFilesPath
}
