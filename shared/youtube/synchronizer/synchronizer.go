package synchronizer

import (
	"log"

	"github.com/nclandrei/synctube/model"
)

// Synchronize takes a map of playlist-video arrays and returns a new map
// containing, per playlist, all videos that need to be downloaded
func Synchronize(videosMap map[model.Playlist][]model.Video) map[model.Playlist][]model.Video {
	toAddVideosMap := make(map[model.Playlist][]model.Video)
	for playlist, videos := range videosMap {
		storedVideos, err := model.VideosByPlaylistID(playlist.ID)
		if err != nil {
			log.Fatalf("Error when retrieving all videos in playlist: %v", err.Error())
		}
		toAddVideos := diffPlaylistVideos(videos, storedVideos)
		toDeleteVideos := diffPlaylistVideos(storedVideos, videos)
		for _, item := range toAddVideos {
			err := model.VideoCreate(item.ID, item.Title, item.PlaylistID)
			if err != nil {
				log.Fatalf("Error in deleting video with (ID: %v, PlaylistID: %v)", item.ID, item.PlaylistID)
			}
		}
		for _, item := range toDeleteVideos {
			err := model.VideoDelete(item.ID, item.PlaylistID)
			if err != nil {
				log.Fatalf("Error in deleting video with (ID: %v, PlaylistID: %v)", item.ID, item.PlaylistID)
			}
		}
		toAddVideosMap[playlist] = toAddVideos
	}
	return toAddVideosMap
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
