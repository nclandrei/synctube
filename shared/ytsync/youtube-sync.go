package ytsync

import (
	"golang.org/x/oauth2"
	"log"
	"context"
	"net/http"
	"fmt"
	"net/url"
	"os"
	"encoding/json"
	"os/user"
	"path/filepath"
	"github.com/nclandrei/YTSync/model"
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

func GetClient (ctx context.Context, code string, userId string) *http.Client {
	if err != nil {
		log.Fatalf("Unable to get path to cached credential file. %v", err)
	}
	tok, err := getTokenFromDb(userId)
	if err != nil {
		tok, _ = getTokenFromWeb(code)
		updateTokenInDb(tok)
	}
	return ytConfig.Client(ctx, tok)
}

// GetAuthorizationURL - uses Config to request a Token.
func GetAuthorizationURL() string {
	authURL := ytConfig.AuthCodeURL("random", oauth2.AccessTypeOffline)
	return authURL
}

// GetTokenFromWeb - given an authorization code it returns the token from a page
func getTokenFromWeb(code string) (*oauth2.Token, error) {
	tok, err := ytConfig.Exchange(context.Background(), code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok, err
}

func handleError(err error, message string) {
	if message == "" {
		message = "Error making API call"
	}
	if err != nil {
		log.Fatalf(message + ": %v", err.Error())
	}
}

// tokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encountered.
func getTokenFromDb(userId string) (*oauth2.Token, error) {
	return model.UserRefreshToken(userId)
}

// saveToken uses a file path to create a file and store the
// token in it.
func saveToken(file string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", file)
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
