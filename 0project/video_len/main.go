package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
)

// traverseDirectory 遍历目录，找出所有视频文件
func traverseDirectory(dir string) ([]string, error) {
	var videoFiles []string
	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		ext := filepath.Ext(path)
		switch ext {
		case ".mp4", ".mov", ".avi", ".mkv":
			videoFiles = append(videoFiles, path)
		}
		return nil
	})
	return videoFiles, err
}

// getVideoDuration 使用 FFmpeg 获取单个视频的时长
func getVideoDuration(filePath string) (float64, error) {
	cmd := exec.Command("D:\\my\\ffmpeg-7.1.1-essentials_build\\bin\\ffmpeg.exe", "-i", filePath)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return 0, fmt.Errorf("failed to run ffmpeg command: %v, output: %s", err, stderr.String())
	}

	output := stderr.String()
	re := regexp.MustCompile(`Duration: (\d{2}):(\d{2}):(\d{2})`)
	match := re.FindStringSubmatch(output)
	if match == nil {
		return 0, fmt.Errorf("failed to extract duration from ffmpeg output for %s", filePath)
	}

	hours, _ := strconv.ParseFloat(match[1], 64)
	minutes, _ := strconv.ParseFloat(match[2], 64)
	seconds, _ := strconv.ParseFloat(match[3], 64)

	totalSeconds := hours*3600 + minutes*60 + seconds
	return totalSeconds, nil
}

func main() {

	dir := "E:\\BaiduNetdiskDownload\\01567 - 马士兵云原生架构师2023"
	videoFiles, err := traverseDirectory(dir)
	if err != nil {
		fmt.Printf("Error traversing directory: %v\n", err)
		return
	}

	var totalDuration float64
	for _, file := range videoFiles {
		duration, err := getVideoDuration(file)
		if err != nil {
			fmt.Printf("Error getting duration for %s: %v\n", file, err)
			continue
		}
		totalDuration += duration
	}

	hours := int(totalDuration / 3600)
	minutes := int((totalDuration - float64(hours*3600)) / 60)

	fmt.Printf("Total video duration in the directory: %dh %dm\n", hours, minutes)
}
