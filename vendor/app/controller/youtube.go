package controller

import (
	"net/http"
	"app/shared/session"
	"app/shared/view"
	"github.com/josephspurrier/csrfbanana"
)

// RegisterGET displays the register page
func YouTubeGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	// Display the view
	v := view.New(r)
	v.Name = "youtube/youtube"
	v.Vars["token"] = csrfbanana.Token(w, r, sess)
	// Refill any form fields
	view.Repopulate([]string{"email"}, r.Form, v.Vars)
	v.Render(w)
}

