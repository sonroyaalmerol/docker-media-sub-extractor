package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

type SubtitleStream struct {
	index    int
	language string
}

func extractSubtitles(videoPath string, processedFilesPath string, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("Extracting subtitles for: %s\n", videoPath)

	// Check if subtitles have already been extracted for this file
	if hasBeenProcessed(videoPath, processedFilesPath) {
		fmt.Printf("Subtitles already extracted for: %s\n", videoPath)
		return
	}

	// Run ffprobe to get information about subtitle streams
	ffprobeCmd := exec.Command("ffprobe", "-v", "error", "-select_streams", "s", "-show_entries", "stream=index,codec_name,language", "-of", "csv=p=0", videoPath)
	ffprobeCmd.Stderr = os.Stderr

	output, err := ffprobeCmd.Output()
	if err != nil {
		fmt.Printf("Error running ffprobe: %v\n", err)
		return
	}

	// Parse ffprobe output to get subtitle stream information
	subtitleStreams := parseSubtitleStreams(string(output))

	// Extract subtitles for each stream
	for _, stream := range subtitleStreams {
		languageCode := getLanguageCode(stream.language)
		outputPath := fmt.Sprintf("%s.%s.srt", videoPath[:len(videoPath)-len(filepath.Ext(videoPath))], languageCode)
		cmd := exec.Command("ffmpeg", "-i", videoPath, "-map", fmt.Sprintf("0:s:%d", stream.index), outputPath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err := cmd.Run()
		if err != nil {
			fmt.Printf("Error extracting subtitles: %v\n", err)
			return
		}

		fmt.Printf("Subtitles extracted successfully for stream %d (Language: %s)\n", stream.index, stream.language)

		// Mark the file as processed
		markAsProcessed(videoPath, processedFilesPath)
	}
}

func hasBeenProcessed(filename string, processedFilesPath string) bool {
	processedFiles, err := readProcessedFiles(processedFilesPath)
	if err != nil {
		fmt.Printf("Error reading processed files: %v\n", err)
		return false
	}

	for _, processedFile := range processedFiles {
		if processedFile == filename {
			return true
		}
	}

	return false
}

func markAsProcessed(filename string, processedFilesPath string) {
	file, err := os.OpenFile(processedFilesPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Printf("Error marking file as processed: %v\n", err)
		return
	}
	defer file.Close()

	if _, err := file.WriteString(filename + "\n"); err != nil {
		fmt.Printf("Error marking file as processed: %v\n", err)
	}
}

func readProcessedFiles(processedFilesPath string) ([]string, error) {
	content, err := os.ReadFile(processedFilesPath)
	if err != nil {
		return nil, err
	}

	return strings.Split(string(content), "\n"), nil
}

func parseSubtitleStreams(output string) []SubtitleStream {
	var streams []SubtitleStream

	lines := splitLines(output)
	for _, line := range lines {
		var index int
		var language string
		var codecName string

		_, err := fmt.Sscanf(line, "%d,%s,%s", &index, &codecName, &language)
		if err == nil && codecName == "subrip" && language != "" {
			streams = append(streams, SubtitleStream{index: index, language: language})
		}
	}

	return streams
}

func splitLines(input string) []string {
	return strings.Split(input, "\n")
}

func getLanguageCode(language string) string {
	languageCode := strings.ToLower(language)
	languageCode = strings.ReplaceAll(languageCode, " ", "_")
	return languageCode
}