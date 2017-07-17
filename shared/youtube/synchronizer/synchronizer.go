package synchronizer

import (
	"fmt"
	"log"
	"net/http"

	"github.com/nclandrei/ytsync/model"
	youtube "google.golang.org/api/youtube/v3"
)

// DownloadLikes returns all user's liked videos given a YouTube service and the user's ID
func DownloadLikes(userID string, client *Service) {
	// First call - will retrieve all items in Likes playlist;
	// needs special call as it is a different kind of playlist
	call := service.Channels.List("contentDetails").Mine(true)

	response, err := call.Do()
	if err != nil {
		// The channels.list method call returned an error.
		log.Fatalf("Error making API call to list channels: %v", err.Error())
	}

	for _, channel := range response.Items {
		var isPlaylistNew bool
		playlistId := channel.ContentDetails.RelatedPlaylists.Likes

		// Print the playlist ID for the list of uploaded videos.
		fmt.Printf("Videos in list %s\r\n", playlistId)

		_, err := model.PlaylistByID(playlistId, userID)

		if err == model.ErrNoResult {
			isPlaylistNew = true
			err := model.PlaylistCreate(playlistId, "Likes", userID)
			if err != nil {
				log.Fatalf("Error creating playlist: %v", err.Error())
			}
			log.Printf("Added Likes playlist for user ID %v", userID)
		} else if err != model.ErrNoResult && err != nil {
			log.Fatalf("Error fetching Likes playlist from the database: %v", err.Error())
		}

		nextPageToken := ""
		var videos []model.Video

		for {
			playlistCall := service.PlaylistItems.List("snippet").
				PlaylistId(playlistId).
				MaxResults(50).
				PageToken(nextPageToken)

			playlistResponse, err := playlistCall.Do()

			if err != nil {
				log.Fatalf("Error fetching playlist items: %v", err.Error())
			}

			for _, playlistItem := range playlistResponse.Items {
				title := playlistItem.Snippet.Title
				videoId := playlistItem.Snippet.ResourceId.VideoId

				video := model.Video{
					ID:         videoId,
					Title:      title,
					PlaylistID: playlistId,
				}

				videos = append(videos, video)
				log.Printf("YouTube video - (title: %v, ID: %v)", title, videoId)
			}

			// Set the token to retrieve the next page of results
			// or exit the loop if all results have been retrieved.
			nextPageToken = playlistResponse.NextPageToken

			if nextPageToken == "" {
				break
			}
			fmt.Println()
		}

		var toAddVideos []model.Video

		if !isPlaylistNew {
			storedVideos, err := model.VideosByPlaylistID(playlistId)
			if err != nil {
				log.Fatalf("Error when retrieving all videos in playlist: %v", err.Error())
			}
			toAddVideos = diffPlaylistVideos(videos, storedVideos)
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
			log.Printf("adding item with title '%v' to mongo", item.Title)
		}

	}
}

func DownloadUserPlaylistVideos(userID string, service *Service) {
	userCreatedPlaylistService := service.Playlists.List("snippet,contentDetails").Mine(true).MaxResults(25)

	userCreatedPlaylists, err := userCreatedPlaylistService.Do()

	if err != nil {
		log.Fatalf("Error making API call to list channels: %v", err.Error())
	}

	for _, item := range userCreatedPlaylists.Items {
		var isPlaylistNew bool

		_, err := model.PlaylistByID(item.Id, userID)

		if err == model.ErrNoResult {
			log.Printf("Could not find playlist in database - will create a new one.")
			isPlaylistNew = true
			err := model.PlaylistCreate(item.Id, item.Snippet.Title, userID)
			if err != nil {
				log.Fatalf("Error creating playlist: %v", err.Error())
			}
			log.Printf("created playlist - %v, %v", item.Snippet.Title, userID)
		} else if err != model.ErrNoResult && err != nil {
			log.Fatalf("Error fetching playlist from the database: %v", err.Error())
		}

		nextPageToken := ""
		var videos []model.Video

		for {
			playlistItems := service.PlaylistItems.List("snippet,contentDetails").
				PlaylistId(item.Id).MaxResults(50).PageToken(nextPageToken)

			playlistResponse, err := playlistItems.Do()

			if err != nil {
				log.Fatalf("Error fetching playlist items: %v", err.Error())
			}

			for _, playlistItem := range playlistResponse.Items {
				title := playlistItem.Snippet.Title
				videoId := playlistItem.Snippet.ResourceId.VideoId

				currentVideo := model.Video{
					ID:         videoId,
					Title:      title,
					PlaylistID: playlistItem.Snippet.PlaylistId,
				}

				videos = append(videos, currentVideo)

				log.Printf("New video with title: %v and id: %v", title, videoId)
			}

			// Set the token to retrieve the next page of results
			// or exit the loop if all results have been retrieved.
			nextPageToken = playlistResponse.NextPageToken

			if nextPageToken == "" {
				break
			}

			fmt.Println()
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

func getVideoByIdFromSlice(x []model.Video, id string) model.Video {
	for _, item := range x {
		if item.ID == id {
			return item
		}
	}
	return model.Video{}
}
