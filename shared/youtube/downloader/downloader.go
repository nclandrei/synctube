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
	tmpFolderPath      string = "tmp"
)

// DownloadYouTubeVideos - downloads all videos in the map for a specific user
func DownloadYouTubeVideos(userID string, videosMap map[model.Playlist][]model.Video) error {
	var err error
	for playlist, videos := range videosMap {
		for _, video := range videos {
			log.Printf("Downloading YouTube video for user %v - (title: %v, ID: %v)", userID, video.Title, video.ID)
			fullURL := fmt.Sprintf("%v%v", youtubePrefix, video.ID)
			output := fmt.Sprintf("--output '%v/%v/%v%%(title)s.%%(ext)s'", tmpFolderPath, userID, playlist.Title)
			args := []string{extractAudio, audioFormat, output, fullURL}
			command := youtubeDownloadCmd
			for _, arg := range args {
				command += " " + arg
			}
			err = exec.Command("bash", "-c", command).Run()
			log.Printf("Downloaded YouTube video for user %v - (title: %v, ID: %v)", userID, video.Title, video.ID)
		}
	}
	return err
}
