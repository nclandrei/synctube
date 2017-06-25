package controller

import (
	"net/http"
	"app/shared/session"
	"app/shared/view"
	"app/shared/youtube-sync"
	"github.com/josephspurrier/csrfbanana"
	"io/ioutil"
	"log"
	"golang.org/x/oauth2/google"
	"golang.org/x/net/context"
	"google.golang.org/api/youtube/v3"
)

// RegisterGET displays the register page
func YouTubeGET(w http.ResponseWriter, r *http.Request) {

		ctx := context.Background()

	b, err := ioutil.ReadFile("client_secret.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved credentials
	// at ~/.credentials/youtube-sync-go-quickstart.json
	config, err := google.ConfigFromJSON(b, youtube.YoutubeReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := youtube_sync.getClient(ctx, config)
	service, err := youtube.New(client)

	youtube_sync.handleError(err, "Error creating YouTube client")

	channelsListByUsername(service, "snippet,contentDetails,statistics", "GoogleDevelopers")



	// Get session
	sess := session.Instance(r)

	// Display the view
	v := view.New(r)
	v.Name = "youtube-sync/youtube-sync"
	v.Vars["token"] = csrfbanana.Token(w, r, sess)
	// Refill any form fields
	view.Repopulate([]string{"email"}, r.Form, v.Vars)
	v.Render(w)
}

