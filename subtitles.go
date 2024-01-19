package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type SubtitleStream struct {
	index    int
	language string
	subtype  string
}

func extractSubtitles(videoPath string, processedFilesPath string) {
	log.Printf("Extracting subtitles for: %s\n", videoPath)

	// Check if subtitles have already been extracted for this file
	if hasBeenProcessed(videoPath, processedFilesPath) {
		log.Printf("Subtitles already extracted for: %s\n", videoPath)
		return
	}

	// Run ffprobe to get information about subtitle streams
	ffprobeCmd := exec.Command("ffprobe", "-v", "error", "-select_streams", "s", "-show_entries", "stream=index,codec_name,language", "-of", "csv=p=0", videoPath)
	ffprobeCmd.Stderr = os.Stderr

	output, err := ffprobeCmd.Output()
	if err != nil {
		log.Printf("Error running ffprobe: %v\n", err)
		return
	}

	// Parse ffprobe output to get subtitle stream information
	subtitleStreams := parseSubtitleStreams(string(output))

	// Extract subtitles for each stream
	for _, stream := range subtitleStreams {
		languageCode := getLanguageCode(stream.language)

		// Extract the base name of the file
		baseName := filepath.Base(videoPath)
		outputPath := fmt.Sprintf("%s.%s.%s", strings.TrimSuffix(baseName, filepath.Ext(baseName)), languageCode, stream.subtype)

		mapString := fmt.Sprintf("0:%d", stream.index)

		log.Printf("ffmpeg -i %s -map %s -c copy %s", videoPath, mapString, outputPath)
		cmd := exec.Command("ffmpeg", "-i", videoPath, "-map", mapString, "-c", "copy", outputPath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err := cmd.Run()
		if err != nil {
			log.Printf("Error extracting subtitles: %v\n", err)
			return
		}

		log.Printf("Subtitles extracted successfully for stream %d (Language: %s)\n", stream.index, stream.language)

		// Mark the file as processed
		markAsProcessed(videoPath, processedFilesPath)
	}
}

func hasBeenProcessed(filename string, processedFilesPath string) bool {
	processedFiles, err := readProcessedFiles(processedFilesPath)
	if err != nil {
		log.Printf("Error reading processed files: %v\n", err)
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
		log.Printf("Error marking file as processed: %v\n", err)
		return
	}
	defer file.Close()

	if _, err := file.WriteString(filename + "\n"); err != nil {
		log.Printf("Error marking file as processed: %v\n", err)
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
		if len(language) == 0 {
			language = "english"
		}
		log.Printf("err: %s, codecName: %s, language: %s\n", err.Error(), codecName, language)
		if err == nil && codecName == "subrip" {
			log.Printf("index %d: appending subrip\n", index)
			streams = append(streams, SubtitleStream{index: index, language: language, subtype: "srt"})
		} else if err == nil && codecName == "ass" {
			log.Printf("index %d: appending ass\n", index)
			streams = append(streams, SubtitleStream{index: index, language: language, subtype: "ass"})
		}
	}

	log.Println(streams)
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
