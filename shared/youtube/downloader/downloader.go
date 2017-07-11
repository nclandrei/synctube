package downloader

import (
	"fmt"
	"os/exec"
)

const (
	youtubePrefix      string = "https://www.youtube.com/watch?v="
	youtubeDownloadCmd string = "youtube-dl"
	extractAudio       string = "--extract-audio"
	audioFormat        string = "--audio-format mp3"
	outputFormat       string = "--output '%(title)s.%(ext)s'"
)

func DownloadYouTubeVideo(url string) error {
	fullURL := fmt.Sprintf("%v%v", youtubePrefix, url)
	args := []string{extractAudio, audioFormat, outputFormat, fullURL}
	command := youtubeDownloadCmd
	for _, arg := range args {
		command += " " + arg
	}
	err := exec.Command("bash", "-c", command).Run()
	return err
}
