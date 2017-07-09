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
)

func DownloadYouTubeVideo(url string) error {
	fullUrl := fmt.Sprintf("%v%v", youtubePrefix, url)
	cmd := exec.Command(youtubeDownloadCmd, extractAudio, audioFormat, fullUrl)
	err := cmd.Run()
	return err
}
