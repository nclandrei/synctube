package downloader

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/nclandrei/synctube/model"
)

const (
	youtubePrefix      string = "https://www.youtube.com/watch?v="
	youtubeDownloadCmd string = "youtube-dl"
	extractAudio       string = "--extract-audio"
	audioFormat        string = "--audio-format mp3"
	outputFormat       string = "--output '%(title)s.%(ext)s'"
	tmpFolderPath      string = "tmp/"
)

func DownloadYouTubeVideos(videosMap map[model.Playlist][]model.Video) error {
	var err error
	for _, videos := range videosMap {
		for _, video := range videos {
			fullURL := fmt.Sprintf("%v%v", youtubePrefix, video.ID)
			args := []string{extractAudio, audioFormat, outputFormat, fullURL}
			command := youtubeDownloadCmd
			for _, arg := range args {
				command += " " + arg
			}
			err = exec.Command("bash", "-c", command).Run()
			log.Printf("Downloading YouTube video - (title: %v, ID: %v)", video.Title, video.ID)
		}
	}
	return err
}
