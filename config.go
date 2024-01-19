package main

import (
	"log"
	"os"
	"strings"
)

var allowedExtensions []string
var processedFilesPath string

func initializeConfig() {
	// Read configuration from environment variables or use default values
	envAllowedExt := os.Getenv("ALLOWED_EXTENSIONS")
	allowedExtensions = strings.Split(envAllowedExt, ",")
	if envAllowedExt == "" || len(allowedExtensions) == 0 {
		allowedExtensions = []string{".mp4", ".mkv", ".avi", ".wmv"}
	}

	processedFilesPath = os.Getenv("PROCESSED_FILES_PATH")
	if processedFilesPath == "" {
		processedFilesPath = "/app/processed_files.txt"
	}
}

func initializeProcessedFiles() string {
	// Create the processed files file if it doesn't exist
	if _, err := os.Stat(processedFilesPath); os.IsNotExist(err) {
		file, err := os.Create(processedFilesPath)
		if err != nil {
			log.Print(err.Error())
			os.Exit(1)
		}
		defer file.Close()
	}
	return processedFilesPath
}
