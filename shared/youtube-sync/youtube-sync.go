package youtube_sync

import (
	"log"
	"net/http"
	"fmt"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"github.com/revel/config"
)

type YT struct {
	ClientID string
	ProjectID string
	AuthURI string
	TokenURI string
	ClientSecret string
	RedirectURI []string
	JavaScriptOrigins []string
}

var (
	ytConfig oauth2.Config
)

func Configure (config oauth2.Config) {
	ytConfig = config
}

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func getClient(ctx context.Context) *http.Client {
	tok := getTokenFromWeb()
	return ytConfig.Client(ctx, tok)
}

// getTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
//noinspection ALL
func getTokenFromWeb() *oauth2.Token {
	authURL := ytConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var code string

	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

func handleError(err error, message string) {
	if message == "" {
		message = "Error making API call"
	}
	if err != nil {
		log.Fatalf(message + ": %v", err.Error())
	}
}
