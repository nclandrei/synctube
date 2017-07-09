package downloader

import (
	"fmt"
	"log"
	"os/exec"
)

const (
	youtubePrefix      string = "https://www.youtube.com/watch?v="
	youtubeDownloadCmd string = "youtube-dl"
	extractAudio       string = "--extract-audio"
	audioFormat        string = "--audio-format mp3"
)

func DownloadYouTubeVideo(url string) error {
	fullURL := fmt.Sprintf("%v%v", youtubePrefix, url)
	cmd := youtubeDownloadCmd
	args := []string{extractAudio, audioFormat, fullURL}
	log.Printf("THIS IS THE COMMAND: %v", exec.Command(cmd, args...))
	err := exec.Command(cmd, args...).Run()
	return err
}
