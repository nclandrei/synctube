package controller

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/nclandrei/synctube/model"
	"github.com/nclandrei/synctube/shared/session"
	"github.com/nclandrei/synctube/shared/youtube/auth"
	"github.com/nclandrei/synctube/shared/youtube/downloader"
	"github.com/nclandrei/synctube/shared/youtube/fetcher"
	"github.com/nclandrei/synctube/shared/youtube/file_manager"
	"github.com/nclandrei/synctube/shared/youtube/synchronizer"
	"google.golang.org/api/youtube/v3"
)

const (
	oauthStateString string = "random"
)

// YouTubeGET - handles getting the authorization code from Google
func YouTubeGET(w http.ResponseWriter, r *http.Request) {
	authURL := auth.GetAuthorizationURL()
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

// YouTubePOST - handles all the processing of the user's YouTube account and then returns
// the zip containing all synced files
func YouTubePOST(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	sess := session.Instance(r)

	if state != oauthStateString {
		fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// get the Google auth code from the session values
	code := r.FormValue("code")
	userID := fmt.Sprintf("%s", sess.Values["id"])

	// create a service that will perform the API calls to YouTube data
	client := auth.GetClient(context.Background(), code, userID)
	service, err := youtube.New(client)
	if err != nil {
		log.Fatalf("Could not retrieve client - %v", err.Error())
	}

	// fetch all videos from all channels for the currently logged in user
	videosMap := fetcher.FetchVideos(userID, service)

	// synchronize and return a playlist-videos_list map containing all videos that
	// have been fetched and are not in the database
	toDownloadVideosMap := synchronizer.Synchronize(videosMap)

	// download all videos previously returned by the synchronizer
	err = downloader.DownloadYouTubeVideos(toDownloadVideosMap)

	// create temporary user and playlist folders, create zip, return it to user
	// and, in the end, clean up everything
	err = file_manager.ManageFiles(userID, toDownloadVideosMap)

	// finally, before redirecting to homepage, save the timestamp of the this sync
	err = model.UserUpdateLastSync(userID, time.Now())
	if err != nil {
		log.Fatalf("Error updating last sync timestamp for user: %v", err.Error())
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
