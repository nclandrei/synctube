package synctube

import (
	"context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
	"io/ioutil"
	"log"
)

func main() {
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
