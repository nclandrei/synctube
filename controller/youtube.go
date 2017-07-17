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
	err := file_manager.CreateUserFolder(userID)

	if err != nil {
		log.Fatalf("Error in creating the user's temporary folder: %v", err.Error())
	}

	client := auth.GetClient(context.Background(), code, userID)

	service, err := youtube.New(client)
	if err != nil {
		log.Fatalf("Could not retrieve client - %v", err.Error())
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

	// user created playlists - will retrieve all items in user created playlists

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

	// Finally, before redirecting to homepage, save the timestamp of the this sync
	err = model.UserUpdateLastSync(userID, time.Now())
	if err != nil {
		log.Fatalf("Error updating last sync timestamp for user: %v", err.Error())
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
