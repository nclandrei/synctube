package fetcher

import (
	"fmt"
	"log"

	youtube "google.golang.org/api/youtube/v3"

	"github.com/nclandrei/synctube/model"
)

// DownloadLikes returns all user's liked videos given a YouTube service and the user's ID
func FetchVideos(userID string, service *youtube.Service) map[model.Playlist][]model.Video {
	playlistVideosMap := make(map[model.Playlist][]model.Video)

	getLikedVideos(userID, service, playlistVideosMap)
	getUserPlaylistVideos(userID, service, playlistVideosMap)

	return playlistVideosMap
}

func getLikedVideos(userID string, service *youtube.Service, playlistVideosMap map[model.Playlist][]model.Video) {
	// First call - will retrieve all items in Likes playlist;
	// needs special call as it is a different kind of playlist
	call := service.Channels.List("contentDetails").Mine(true)

	likesPlaylist, err := call.Do()
	if err != nil {
		// The channels.list method call returned an error.
		log.Fatalf("Error making API call to list channels: %v", err.Error())
	}

	for _, channel := range likesPlaylist.Items {
		playlistId := channel.ContentDetails.RelatedPlaylists.Likes

		// Print the playlist ID for the list of uploaded videos.
		fmt.Printf("Videos in list %s\r\n", playlistId)

		playlist, err := model.PlaylistByID(playlistId, userID)

		if err == model.ErrNoResult {
			err := model.PlaylistCreate(playlistId, "Likes", userID)
			if err != nil {
				log.Fatalf("Error creating playlist: %v", err.Error())
			}
			log.Printf("Created Likes playlist for user ID %v", userID)
		} else if err != model.ErrNoResult && err != nil {
			log.Fatalf("Error fetching Likes playlist from the database: %v", err.Error())
		}
		getVideosFromPlaylist(service, playlist, playlistVideosMap)
	}
}

// FetchUserPlaylistVideos - returns all playlists created by a specific
// user along with all their videos
func getUserPlaylistVideos(userID string, service *youtube.Service, playlistVideosMap map[model.Playlist][]model.Video) {
	userCreatedPlaylistService := service.Playlists.List("snippet,contentDetails").Mine(true).MaxResults(25)

	userCreatedPlaylists, err := userCreatedPlaylistService.Do()

	if err != nil {
		log.Fatalf("Error making API call to list channels: %v", err.Error())
	}

	for _, item := range userCreatedPlaylists.Items {

		playlist, err := model.PlaylistByID(item.Id, userID)

		if err == model.ErrNoResult {
			log.Printf("Could not find playlist in database - will create a new one.")
			err := model.PlaylistCreate(item.Id, item.Snippet.Title, userID)
			if err != nil {
				log.Fatalf("Error creating playlist: %v", err.Error())
			}
			log.Printf("created playlist - %v, %v", item.Snippet.Title, userID)
		} else if err != model.ErrNoResult && err != nil {
			log.Fatalf("Error fetching playlist from the database: %v", err.Error())
		}
		getVideosFromPlaylist(service, playlist, playlistVideosMap)
	}
}

func getVideosFromPlaylist(service *youtube.Service, playlist model.Playlist, playlistVideosMap map[model.Playlist][]model.Video) {
	var videos []model.Video
	nextPageToken := ""

	for {
		playlistCall := service.PlaylistItems.List("snippet").
			PlaylistId(playlist.ID).
			MaxResults(50).
			PageToken(nextPageToken)

		playlistResponse, err := playlistCall.Do()

		if err != nil {
			log.Fatalf("Error fetching playlist items: %v", err.Error())
		}

		for _, playlistItem := range playlistResponse.Items {
			title := playlistItem.Snippet.Title
			videoId := playlistItem.Snippet.ResourceId.VideoId
			thumbnailURL := playlistItem.Snippet.Thumbnails.Default.Url

			video := model.Video{
				ID:           videoId,
				Title:        title,
				PlaylistID:   playlist.ID,
				ThumbnailURL: thumbnailURL,
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

	playlistVideosMap[playlist] = videos
}
