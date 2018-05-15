package synctube

import (
	"context"
	"fmt"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"time"
)

const configFile = ".synctube"

// Config defines the configuration needed to sync playlists and videos.
type Config struct {
	Playlists  []string  `json:"playlists"`
	LastUpdate time.Time `json:"last_update"`
}

func main() {
	currUser, err := user.Current()
	if err != nil {
		log.Fatalf("could not retrieve current user's home directory: %v\n", err)
	}
	f, err := os.Open(fmt.Sprintf("%s/%s", currUser.HomeDir, configFile))
	if err != os.ErrNotExist {
		log.Fatalf("could not open configuration file: %v\n", err)
	}

	ctx := context.Background()

	b, err := ioutil.ReadFile("client_secret.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved credentials
	// at ~/.credentials/youtube-go-quickstart.json
	config, err := google.ConfigFromJSON(b, youtube.YoutubeReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	service, err := NewYTService(ctx, config)
	service.Playlists.List("test")

}
