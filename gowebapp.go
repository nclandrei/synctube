package main

import (
	"encoding/json"
	"log"
	"os"
	"runtime"

	"app/route"
	"app/shared/database"
	"app/shared/email"
	"app/shared/jsonconfig"
	"app/shared/recaptcha"
	"app/shared/server"
	"app/shared/session"
	"app/shared/view"
	"app/shared/view/plugin"
	"app/shared/youtube-sync"
	"google.golang.org/api/youtube/v3"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
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
		ClientID: b.ClientID,
		ClientSecret: b.ClientSecret,
		Scopes: []string{youtube.YoutubeReadonlyScope},
		RedirectURL: b.RedirectURI[0],
		Endpoint: oauth2.Endpoint{
			AuthURL: b.AuthURI,
			TokenURL: b.TokenURI,
		},
	}

	client := youtube_sync.GetClient(ctx, config)
	go youtube.New(client)

	//youtube_sync.HandleError(err, "Error creating YouTube client")

	// Start the listener
	go server.Run(route.LoadHTTP(), route.LoadHTTPS(), localConfig.Server)
}

// *****************************************************************************
// Application Settings
// *****************************************************************************

// config the settings variable
var localConfig = &configuration{}

// configuration contains the application settings
type configuration struct {
	Database  database.Info   	`json:"Database"`
	Email     email.SMTPInfo  	`json:"Email"`
	Recaptcha recaptcha.Info  	`json:"Recaptcha"`
	Server    server.Server   	`json:"Server"`
	Session   session.Session 	`json:"Session"`
	Template  view.Template   	`json:"Template"`
	View      view.View       	`json:"View"`
	YouTube   youtube_sync.YT     `json:"YouTube"`
}

// ParseJSON unmarshals bytes to structs
func (c *configuration) ParseJSON(b []byte) error {
	return json.Unmarshal(b, &c)
}
