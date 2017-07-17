package main

import "github.com/nclandrei/SyncTube/model"

func DiffVideos() {
	var toAddVideos []model.Video

	if !isPlaylistNew {
		storedVideos, err := model.VideosByPlaylistID(item.Id)
		if err != nil {
			log.Fatalf("Error when retrieving all videos in playlist: %v", err.Error())
		}
		toAddVideos = diffPlaylistVideos(videos, storedVideos)
		log.Printf("Number of videos to add: %v", len(toAddVideos))
		toDeleteVideos := diffPlaylistVideos(storedVideos, videos)
		for _, item := range toDeleteVideos {
			model.VideoDelete(item.ID, item.PlaylistID)
		}
	} else {
		toAddVideos = videos
	}

	for _, item := range toAddVideos {
		err := model.VideoCreate(item.ID, item.Title, item.PlaylistID)
		if err != nil {
			log.Fatalf("Error adding the video to the database: %v", err.Error())
		}
		log.Printf("Added video - (title: %v, ID: %v) to database", item.Title, item.ID)
		err = downloader.DownloadYouTubeVideo(item.ID)
		if err != nil {
			log.Fatalf("Error downloading video (title: %v, ID: %v) from YouTube - %v", item.Title, item.ID, err.Error())
		}
		log.Printf("Downloaded video - (title: %v, ID: %v)", item.Title, item.ID)
	}
	file_manager.CreatePlaylistFolder(item.Snippet.Title)
	}

}

// Function that returns the videos that are in first slice but not in the second one
func diffPlaylistVideos(X, Y []model.Video) []model.Video {
	var resultSlice []model.Video
	m := make(map[string]int)

	for _, y := range Y {
		m[y.ID]++
	}

	for _, x := range X {
		if m[x.ID] > 0 {
			m[x.ID]--
			continue
		}
		video := getVideoByIdFromSlice(X, x.ID)
		if video != (model.Video{}) {
			resultSlice = append(resultSlice, x)
		}
	}

	return resultSlice
}

// Returns a video from the array when given its id
func getVideoByIdFromSlice(x []model.Video, id string) model.Video {
	for _, item := range x {
		if item.ID == id {
			return item
		}
	}
	return model.Video{}
}
