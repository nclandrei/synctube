package main

import (
	"encoding/json"
	"log"
	"os"
	"runtime"

	"github.com/nclandrei/YTSync/route"
	"github.com/nclandrei/YTSync/shared/database"
	"github.com/nclandrei/YTSync/shared/email"
	"github.com/nclandrei/YTSync/shared/jsonconfig"
	"github.com/nclandrei/YTSync/shared/recaptcha"
	"github.com/nclandrei/YTSync/shared/server"
	"github.com/nclandrei/YTSync/shared/session"
	"github.com/nclandrei/YTSync/shared/view"
	"github.com/nclandrei/YTSync/shared/view/plugin"
	"github.com/nclandrei/YTSync/shared/youtube-sync"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"google.golang.org/api/youtube/v3"
	"fmt"
)

// *****************************************************************************
// Application Logic
// *****************************************************************************

func init() {
	// Verbose logging with file name and line number
	log.SetFlags(log.Lshortfile)

	// Use all CPU cores
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	// Load the configuration file
	jsonconfig.Load("config"+string(os.PathSeparator)+"config.json", localConfig)

	// Configure the session cookie store
	session.Configure(localConfig.Session)

	// Connect to database
	database.Connect(localConfig.Database)

	// Configure the Google reCAPTCHA prior to loading view plugins
	recaptcha.Configure(localConfig.Recaptcha)

	// Configure YouTube specific settings
	youtube_sync.Configure(localConfig.YouTube)

	// Setup the views
	view.Configure(localConfig.View)
	view.LoadTemplates(localConfig.Template.Root, localConfig.Template.Children)
	view.LoadPlugins(
		plugin.TagHelper(localConfig.View),
		plugin.NoEscape(),
		plugin.PrettyTime(),
		recaptcha.Plugin())

	ctx := context.Background()

	b := localConfig.YouTube

	// If modifying these scopes, delete your previously saved credentials
	// at ~/.credentials/youtube-go-quickstart.json
	config := &oauth2.Config{
		ClientID:     b.ClientID,
		ClientSecret: b.ClientSecret,
		Scopes:       []string{youtube.YoutubeReadonlyScope},
		RedirectURL:  b.RedirectURI[0],
		Endpoint: oauth2.Endpoint{
			AuthURL:  b.AuthURI,
			TokenURL: b.TokenURI,
		},
	}

	//youtube_sync.HandleError(err, "Error creating YouTube client")

	// Start the listener
	go server.Run(route.LoadHTTP(), route.LoadHTTPS(), localConfig.Server)

	testYouTube(ctx, *config)
}

// *****************************************************************************
// Application Settings
// *****************************************************************************

// config the settings variable
var localConfig = &configuration{}

// configuration contains the application settings
type configuration struct {
	Database  database.Info   `json:"Database"`
	Email     email.SMTPInfo  `json:"Email"`
	Recaptcha recaptcha.Info  `json:"Recaptcha"`
	Server    server.Server   `json:"Server"`
	Session   session.Session `json:"Session"`
	Template  view.Template   `json:"Template"`
	View      view.View       `json:"View"`
	YouTube   youtube_sync.YT `json:"YouTube"`
}

// ParseJSON unmarshals bytes to structs
func (c *configuration) ParseJSON(b []byte) error {
	return json.Unmarshal(b, &c)
}

func testYouTube(ctx context.Context, config oauth2.Config) {
	client := youtube_sync.GetClient(ctx, &config)
	service, err := youtube.New(client)

	call := service.Channels.List("contentDetails").Mine(true)

	response, err := call.Do()
	if err != nil {
		// The channels.list method call returned an error.
		log.Fatalf("Error making API call to list channels: %v", err.Error())
	}

	for _, channel := range response.Items {
		playlistId := channel.ContentDetails.RelatedPlaylists.Uploads
		// Print the playlist ID for the list of uploaded videos.
		fmt.Printf("Videos in list %s\r\n", playlistId)

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
}
