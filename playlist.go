package synctube

import (
	"hash"
	"net/url"
)

// Playlist describes the information held by a YouTube playlist.
type Playlist struct {
	Name   string
	Videos []Video
	hash   hash.Hash
}

// Video describes the information held by a YouTube video.
type Video struct {
	Name   string
	URL    url.URL
	Length float64
}
