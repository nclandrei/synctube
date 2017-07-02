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

func GetClient (code string) *http.Client {
	token, err := getTokenFromWeb(code)
	if err != nil {
		fmt.Printf("oauthConf.Exchange() failed with '%s'\n", err)
	}
	return ytConfig.Client(context.Background(), token)
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

// tokenCacheFile generates credential file path/filename.
// It returns the generated credential path/filename.
func tokenCacheFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	tokenCacheDir := filepath.Join(usr.HomeDir, ".credentials")
	os.MkdirAll(tokenCacheDir, 0700)
	return filepath.Join(tokenCacheDir,
		url.QueryEscape("youtube-go-quickstart.json")), err
}

// tokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encountered.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
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
