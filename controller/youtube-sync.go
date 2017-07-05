package controller

import (
	"fmt"
	"github.com/nclandrei/YTSync/shared/ytsync"
	"google.golang.org/api/youtube/v3"
	"net/http"
	"context"
	"github.com/nclandrei/YTSync/shared/session"
	//"log"
	"log"
    "github.com/nclandrei/YTSync/model"
)

const (
	oauthStateString string = "random"
    youtubeVideoURLPrefix string = "https://www.youtube.com/watch?v="
)

func YouTubeGET(w http.ResponseWriter, r *http.Request) {
	authURL := ytsync.GetAuthorizationURL()
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

	client := ytsync.GetClient(context.Background(), code, userID)

	service, err := youtube.New(client)
	if err != nil {
		fmt.Errorf("Could not retrieve client - %v", err.Error())
	}

	// Start making YouTube API calls.
	// Call the channels.list method. Set the mine parameter to true to
	// retrieve the playlist ID for uploads to the authenticated user's
	// channel.
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
			// Call the playlistItems.list method to retrieve the
			// list of uploaded videos. Each request retrieves 50
			// videos until all videos have been retrieved.
			playlistCall := service.PlaylistItems.List("snippet").
				PlaylistId(playlistId).
				MaxResults(50).
				PageToken(nextPageToken)

			playlistResponse, err := playlistCall.Do()

			if err != nil {
				// The playlistItems.list method call returned an error.
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

	//callTwo := service.Playlists.List("snippet,contentDetails").Mine(true).MaxResults(25)
	//responseTwo, err := callTwo.Do()
	//if err != nil {
	//	// The channels.list method call returned an error.
	//	log.Fatalf("Error making API call to list channels: %v", err.Error())
	//}
	//
	//for _, item := range responseTwo.Items {
	//
	//	fmt.Printf("Videos in playlsit --- %s, %s\r\n", item.Id, item.Snippet.Title)
	//
	//	nextPageToken := ""
	//	for {
	//		playlistItems := service.PlaylistItems.List("snippet,contentDetails").
	//			PlaylistId(item.Id).MaxResults(50).PageToken(nextPageToken)
	//
	//		playlistResponse, err := playlistItems.Do()
	//
	//		if err != nil {
	//			// The playlistItems.list method call returned an error.
	//			log.Fatalf("Error fetching playlist items: %v", err.Error())
	//		}
	//
	//		for _, playlistItem := range playlistResponse.Items {
	//			title := playlistItem.Snippet.Title
	//			videoId := playlistItem.Snippet.ResourceId.VideoId
	//			fmt.Printf("%v, (%v)\r\n", title, videoId)
	//		}
	//
	//		// Set the token to retrieve the next page of results
	//		// or exit the loop if all results have been retrieved.
	//		nextPageToken = playlistResponse.NextPageToken
	//		if nextPageToken == "" {
	//			break
	//		}
	//		fmt.Println()
	//	}
	//}
	http.Redirect(w, r, "/", http.StatusFound)
}
