package controller

import (
	"log"
	"net/http"

	"github.com/nclandrei/synctube/model"
	"github.com/nclandrei/synctube/shared/session"
	"github.com/nclandrei/synctube/shared/view"
)

// IndexGET displays the home page
func IndexGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	session := session.Instance(r)

	// Display the view
	v := view.New(r)

	log.Printf("map: %v", session.Values["videosMap"])

	if videosMap, videosMapPresent := session.Values["videosMap"].(map[model.Playlist][]model.Video); videosMapPresent {
		v.Vars["videosMap"] = videosMap
	} else {
		log.Printf("Videos map not present yet.")
	}

	log.Printf("Vars received: %v", v.Vars)

	if session.Values["id"] != nil {
		v.Name = "index/auth"
		v.Vars["first_name"] = session.Values["first_name"]
	} else {
		v.Name = "index/anon"
	}
	v.Render(w)
}
