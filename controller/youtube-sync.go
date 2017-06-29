package controller

import (
	"net/http"
	"github.com/nclandrei/YTSync/shared/view"
	"github.com/revel/config"
	"golang.org/x/oauth2"
)

func YouTubeGET(w http.ResponseWriter, r *http.Request) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

func YouTubePOST(w http.ResponseWriter, r *http.Request) {
	v := view.New(r)
	v.Name = "about/about"
	v.Render(w)
}
