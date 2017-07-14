// TODO: add number of videos and some sort of display for video information on auth template
// TODO: add email verification for new user accounts

package controller

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/nclandrei/YTSync/model"
	"github.com/nclandrei/YTSync/shared/file_manager"
	"github.com/nclandrei/YTSync/shared/session"
	"github.com/nclandrei/YTSync/shared/youtube/auth"
	"github.com/nclandrei/YTSync/shared/youtube/downloader"
	"google.golang.org/api/youtube/v3"
)

const (
	oauthStateString string = "random"
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

	// create this user's temporary folder where the zip will be created
	file_manager.CreateUserFolder(userID)

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

	// user created playlists - will retrieve all items in user created playlists
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

		var toAddVideos []model.Video

		if !isPlaylistNew {
			storedVideos, err := model.VideosByPlaylistID(item.Id)
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
			item := item
			go func() {
				err := model.VideoCreate(item.ID, item.Title, item.PlaylistID)
				if err != nil {
					log.Fatalf("Error adding the video to the database: %v", err.Error())
				}
				log.Printf("Added item with title '%v' to database", item.Title)
			}()
			go func() {
				err = downloader.DownloadYouTubeVideo(item.ID)
				if err != nil {
					log.Fatalf("Error downloading video with ID %v from YouTube - %v", item.ID, err.Error())
				}
			}()
		}
		file_manager.CreatePlaylistFolder(item.Snippet.Title)
	}

	// Finally, before redirecting to homepage, save the timestamp of the this sync
	err = model.UserUpdateLastSync(userID, time.Now())
	if err != nil {
		log.Fatalf("Error updating last sync timestamp for user: %v", err.Error())
	}

	http.Redirect(w, r, "/", http.StatusFound)
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
