package ytsync

import (
	"golang.org/x/oauth2"
	"log"
)

/**
YT - struct that holds all necessary information
regarding the project; gets populated through config
files
*/
type YT struct {
	ClientID          string
	ProjectID         string
	AuthURI           string
	TokenURI          string
	ClientSecret      string
	RedirectURI       []string
	JavaScriptOrigins []string
}

var (
	ytConfig oauth2.Config
)

// function that sets up the YT instance
func Configure(config oauth2.Config) {
	ytConfig = config
}

// GetAuthorizationURL - uses Config to request a Token.
func GetAuthorizationURL() string {
	authURL := ytConfig.AuthCodeURL("random_token")
	return authURL
}

// GetTokenFromWeb - given an authorization code it returns the token from a page
func GetTokenFromWeb(code string) (*oauth2.Token, error) {
	tok, err := ytConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok, err
}

/*func handleError(err error, message string) {
	if message == "" {
		message = "Error making API call"
	}
	if err != nil {
		log.Fatalf(message + ": %v", err.Error())
	}
}*/
