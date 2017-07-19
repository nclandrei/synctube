package main

import (
	"log"

	"github.com/nclandrei/SyncTube/model"
)

// DiffVideos - takes in two lists of videos, the ones fetched from YouTube and the ones that
// reside in the database and returns two lists containing videos to be added in the DB and
// videos that need to be deleted from storage
func DiffVideos(dbVideos []model.Video, fetchedVideos []model.Video) ([]model.Video, []model.Video) {
	var toAddVideos, toDeleteVideos []model.Video

	toAddVideos = diffPlaylistVideos(fetchedVideos, dbVideos)
	log.Printf("Number of videos to add: %v", len(toAddVideos))
	toDeleteVideos = diffPlaylistVideos(dbVideos, fetchedVideos)
	return toAddVideos, toDeleteVideos
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
