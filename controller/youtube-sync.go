package controller

import (
	"net/http"
	"github.com/nclandrei/YTSync/shared/view"
)

func YouTubeGET(w http.ResponseWriter, r *http.Request) {
	v := view.New(r)
	v.Name = "about/about"
	v.Render(w)
}

func YouTubePOST(w http.ResponseWriter, r *http.Request) {
	v := view.New(r)
	v.Name = "about/about"
	v.Render(w)
}
