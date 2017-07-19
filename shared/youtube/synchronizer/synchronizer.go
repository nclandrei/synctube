package synchronizer

import (
	"log"

	"github.com/nclandrei/synctube/model"
)

func Synchronize(videosMap map[string][]model.Video) []model.Video {
	var toAddVideos []model.Video
	for playlistID, videos := range videosMap {
		storedVideos, err := model.VideosByPlaylistID(playlistID)
		if err != nil {
			log.Fatalf("Error when retrieving all videos in playlist: %v", err.Error())
		}
		toAddVideos = diffPlaylistVideos(videos, storedVideos)
		toDeleteVideos := diffPlaylistVideos(storedVideos, videos)
		for _, item := range toAddVideos {
			model.VideoCreate(item.ID, item.Title, item.PlaylistID)
		}
		for _, item := range toDeleteVideos {
			model.VideoDelete(item.ID, item.PlaylistID)
		}
	}
	return toAddVideos
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
