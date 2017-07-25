package controller

import (
	"log"
	"net/http"

	"github.com/nclandrei/synctube/shared/session"
	"github.com/nclandrei/synctube/shared/view"
)

// IndexGET displays the home page
func IndexGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	session := session.Instance(r)

	// Display the view
	v := view.New(r)

	videosToDownload := session.Values["videos_map"]
	v.Vars["videos_map"] = videosToDownload

	log.Printf("Map received: %v", videosToDownload)

	if session.Values["id"] != nil {
		v.Name = "index/auth"
		v.Vars["first_name"] = session.Values["first_name"]
	} else {
		v.Name = "index/anon"
	}
	v.Render(w)
}
