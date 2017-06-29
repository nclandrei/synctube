package controller

import (
	"net/http"
	"github.com/nclandrei/YTSync/shared/view"
	"github.com/revel/config"
	"golang.org/x/oauth2"
	"github.com/nclandrei/YTSync/shared/youtube-sync"
)

func YouTubeGET(w http.ResponseWriter, r *http.Request) {
	authURL := youtube_sync.GetAuthorizationURL()
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

func YouTubePOST(w http.ResponseWriter, r *http.Request) {
	v := view.New(r)
	v.Name = "about/about"
	v.Render(w)
}
