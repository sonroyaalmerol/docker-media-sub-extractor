package main

import (
	"os"
	"path/filepath"
	"time"
)

func main() {
	initializeConfig()

	// Create the processed files file if it doesn't exist
	processedFilesPath := initializeProcessedFiles()

	for {
		filepath.Walk(mediaPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() && hasAllowedExtension(info.Name()) {
				extractSubtitles(path, processedFilesPath)
			}

			return nil
		})

		// Sleep for 5 minutes before scanning again
		time.Sleep(5 * time.Minute)
	}
}

func hasAllowedExtension(filename string) bool {
	ext := filepath.Ext(filename)
	for _, allowedExt := range allowedExtensions {
		if ext == allowedExt {
			return true
		}
	}
	return false
}
