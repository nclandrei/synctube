package main

import (
	"github.com/rylio/ytdl"
	"log"
	"os"
	"strings"
)

func main() {
	vi, err := ytdl.GetVideoInfo("https://www.youtube.com/watch?v=AHNtfiovZhc")
	if err != nil {
		log.Fatalf("could not get video info: %v\n", err)
	}

	var format ytdl.Format
	for _, f := range vi.Formats {
		if f.Extension != "mp4" {
			continue
		}
		if f.AudioBitrate > format.AudioBitrate {
			format = f
		}
	}

	filename := strings.Replace(vi.Title, " ", "_", -1)
	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("could not create file: %v\n", err)
	}
	defer file.Close()


	err = vi.Download(format, file)
	if err != nil {
		log.Fatalf("could not download file: %v\n", err)
	}
}
