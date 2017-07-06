package controller

import (
	"fmt"
	"github.com/nclandrei/YTSync/shared/youtube/auth"
	"google.golang.org/api/youtube/v3"
	"net/http"
	"context"
	"github.com/nclandrei/YTSync/shared/session"
	//"log"
	"log"
    "github.com/nclandrei/YTSync/model"
	"time"
)

const (
	oauthStateString string = "random"
    youtubeVideoURLPrefix string = "https://www.youtube.com/watch?v="
)

func YouTubeGET(w http.ResponseWriter, r *http.Request) {
	authURL := auth.GetAuthorizationURL()
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

func YouTubePOST(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	sess := session.Instance(r)

	if state != oauthStateString {
		fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	userID := fmt.Sprintf("%s", sess.Values["id"])

	client := auth.GetClient(context.Background(), code, userID)

	service, err := youtube.New(client)
	if err != nil {
		fmt.Errorf("Could not retrieve client - %v", err.Error())
	}

	// First call - will retrieve all items in Likes playlist;
	// needs special call as it is a different kind of playlist
	call := service.Channels.List("contentDetails").Mine(true)

	response, err := call.Do()
	if err != nil {
		// The channels.list method call returned an error.
		log.Fatalf("Error making API call to list channels: %v", err.Error())
	}

	for _, channel := range response.Items {
		playlistId := channel.ContentDetails.RelatedPlaylists.Likes
		// Print the playlist ID for the list of uploaded videos.
		fmt.Printf("Videos in list %s\r\n", playlistId)

        model.PlaylistCreate(playlistId, "likes", userID)

		nextPageToken := ""
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
                 videoURL := youtubeVideoURLPrefix + playlistItem.Snippet.ResourceId.VideoId
                 if err != nil {
                     log.Fatalf("Error while trying to build video URL: %v", err.Error())
                 }
                 model.VideoCreate(videoId, playlistItem.Snippet.Title, videoURL, playlistItem.Snippet.PlaylistId)
				fmt.Printf("%v, (%v)\r\n", title, videoId)
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

	// Second call - will retrieve all items in user created playlists
	userCreatedPlaylistService := service.Playlists.List("snippet,contentDetails").Mine(true).MaxResults(25)
	userCreatedPlaylists, err := userCreatedPlaylistService.Do()

	if err != nil {
		log.Fatalf("Error making API call to list channels: %v", err.Error())
	}

	for _, item := range userCreatedPlaylists.Items {

		fmt.Printf("Videos in playlsit --- %s, %s\r\n", item.Id, item.Snippet.Title)

		playlist, _ := model.PlaylistByID(userID, item.Id)

		if playlist == (model.Playlist{}) {
			model.PlaylistCreate(item.Id, item.Snippet.Title, userID)
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
				videoURL := youtubeVideoURLPrefix + playlistItem.Snippet.ResourceId.VideoId

				currentVideo := model.Video{
					ID: videoId,
					Title: title,
					URL: videoURL,
					PlaylistID: playlistItem.Snippet.PlaylistId,
				}

				videos = append(videos, currentVideo)

				fmt.Printf("%v, (%v)\r\n", title, videoId)
			}

			// Set the token to retrieve the next page of results
			// or exit the loop if all results have been retrieved.
			nextPageToken = playlistResponse.NextPageToken
			if nextPageToken == "" {
				break
			}
			fmt.Println()
		}
		storedVideos, _ := model.VideoByPlaylistID(item.Id)
		toAddVideos := diffPlaylistVideos(videos, storedVideos)
		toDeleteVideos := diffPlaylistVideos(storedVideos, videos)

		for _, item := range toAddVideos {
			model.VideoCreate(item.ID, item.Title, item.URL, item.PlaylistID)
		}

		for _, item := range toDeleteVideos {
			model.VideoDelete(item.ID, item.PlaylistID)
		}
	}

	// Finally, before redirecting to homepage, save the timestamp of the this sync
	err = model.UserUpdateLastSync(userID, time.Now())
	if err != nil {
		log.Fatalf("Error updating last sync timestamp for user: %v", err.Error())
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

// Function that returns the videos that are in first slice but not in the second one
func diffPlaylistVideos(X, Y []model.Video) ([]model.Video) {
	counts := make(map[model.Video]int)
	var total int
	for _, val := range X {
		counts[val] += 1
		total += 1
	}

	for _, val := range Y {
		if count := counts[val]; count > 0 {
			counts[val] -= 1
			total -= 1
		}

	}

	diff := make([]model.Video, total)
	i := 0

	for val, count := range counts {
		for j := 0; j < count; j++ {
			diff[i] = val
			i++
		}
	}
	return diff
}
