package main

import (
	"encoding/json"
	"log"
	"os"
	"runtime"

	"github.com/nclandrei/synctube/route"
	"github.com/nclandrei/synctube/shared/database"
	"github.com/nclandrei/synctube/shared/email"
	"github.com/nclandrei/synctube/shared/jsonconfig"
	"github.com/nclandrei/synctube/shared/recaptcha"
	"github.com/nclandrei/synctube/shared/server"
	"github.com/nclandrei/synctube/shared/session"
	"github.com/nclandrei/synctube/shared/view"
	"github.com/nclandrei/synctube/shared/view/plugin"
	"github.com/nclandrei/synctube/shared/youtube/auth"
	"github.com/nclandrei/synctube/shared/youtube/file_manager"
	"golang.org/x/oauth2"
	"google.golang.org/api/youtube/v3"
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

	// Configure YouTube settings
	auth.Configure(loadYTConfig(localConfig.YouTube))

	// Initialize downloads folder for all users
	file_manager.Init()

	// Setup the views
	view.Configure(localConfig.View)
	view.LoadTemplates(localConfig.Template.Root, localConfig.Template.Children)
	view.LoadPlugins(
		plugin.TagHelper(localConfig.View),
		plugin.NoEscape(),
		plugin.PrettyTime(),
		recaptcha.Plugin())

	// Start the listener
	server.Run(route.LoadHTTP(), route.LoadHTTPS(), localConfig.Server)
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
	YouTube   auth.YT         `json:"YouTube"`
}

// ParseJSON unmarshals bytes to structs
func (c *configuration) ParseJSON(b []byte) error {
	return json.Unmarshal(b, &c)
}

// loadYTConfig parses the information needed for setting up YouTube
func loadYTConfig(conf auth.YT) oauth2.Config {
	return oauth2.Config{
		ClientID:     conf.ClientID,
		ClientSecret: conf.ClientSecret,
		Scopes:       []string{youtube.YoutubeReadonlyScope},
		RedirectURL:  conf.RedirectURI[0],
		Endpoint: oauth2.Endpoint{
			AuthURL:  conf.AuthURI,
			TokenURL: conf.TokenURI,
		},
	}
}
