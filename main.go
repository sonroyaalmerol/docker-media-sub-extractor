package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var wg sync.WaitGroup

func main() {
	mediaPath := os.Getenv("MEDIA_PATH")
	if mediaPath == "" {
		fmt.Println("MEDIA_PATH environment variable not set.")
		os.Exit(1)
	}

	// Create the processed files file if it doesn't exist
	processedFilesPath := initializeProcessedFiles()

	for {
		filepath.Walk(mediaPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() && hasAllowedExtension(info.Name()) {
				wg.Add(1)
				go extractSubtitles(path, processedFilesPath, &wg)
			}

			return nil
		})

		wg.Wait() // Wait for all goroutines to finish before the next iteration

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
