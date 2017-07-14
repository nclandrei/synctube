package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/nclandrei/YTSync/model"
	"golang.org/x/oauth2"
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

// Configure - used for setting up the YT instance
func Configure(config oauth2.Config) {
	ytConfig = config
}

func GetClient(ctx context.Context, code string, userId string) *http.Client {
	tok, err := getTokenFromDb(userId)
	fmt.Println(tok)
	if err != nil || !tok.Valid() {
		tok, _ = getTokenFromWeb(code)
		updateTokenInDb(userId, tok)
	}
	return ytConfig.Client(ctx, tok)
}

// GetAuthorizationURL - uses Config to request a Token.
func GetAuthorizationURL() string {
	url := ytConfig.AuthCodeURL("random", oauth2.AccessTypeOffline)
	return url
}

// given an authorization code it returns the token from a page
func getTokenFromWeb(code string) (*oauth2.Token, error) {
	tok, err := ytConfig.Exchange(context.Background(), code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok, err
}

// returns the token for a specific user from the database
func getTokenFromDb(userId string) (*oauth2.Token, error) {
	currentUser, err := model.UserByToken(userId)
	if err != nil {
		log.Fatalf("Unable to retrieve user by id: %v", err)
	}
	return &currentUser.Token, err
}

// given a new token it replaces the already existing one in DB
func updateTokenInDb(userID string, token *oauth2.Token) {
	err := model.UpdateUserToken(userID, *token)
	if err != nil {
		log.Fatalf("Unable to update user's token: %v", err)
	}
}

func handleError(err error, message string) {
	if message == "" {
		message = "Error making API call"
	}
	if err != nil {
		log.Fatalf(message+": %v", err.Error())
	}
}
