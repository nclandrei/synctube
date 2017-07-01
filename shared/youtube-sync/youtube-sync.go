package youtube_sync

import (
	"log"
	"golang.org/x/oauth2"
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

// getTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
//noinspection ALL
func GetAuthorizationURL() string {
	authURL := ytConfig.AuthCodeURL("random_token")
	return authURL
}

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
