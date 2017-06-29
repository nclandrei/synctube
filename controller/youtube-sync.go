package controller

import (
	"net/http"
	"golang.org/x/oauth2"
	"github.com/nclandrei/YTSync/shared/youtube-sync"
	"fmt"
	"context"
	"google.golang.org/api/youtube/v3"
)

const (
	oauthStateString string = "random_token"
)

func YouTubeGET(w http.ResponseWriter, r *http.Request) {
	authURL := youtube_sync.GetAuthorizationURL()
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

func YouTubePOST(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")

	if state != oauthStateString {
		fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	token, err := youtube_sync.GetTokenFromWeb(code)

	if err != nil {
		fmt.Printf("oauthConf.Exchange() failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	client := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(token))

	driveService, err := youtube.New(client)

	myR, err := driveService.Files.List().MaxResults(10).Do()
	if err != nil {
		fmt.Fprintf(w, "Couldn't retrieve files ", err)
	}
	if len(myR.Items) > 0 {
		for _, i := range myR.Items {
			fmt.Fprintf(w, i.Title, " ", i.Id)
		}
	} else {
		fmt.Fprintf(w, "No files found.")
	}
}
