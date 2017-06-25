//package controller
//
//import (
//	"net/http"
//	"app/shared/youtube"
//	"fmt"
//)
//
//// Get user's playlists after getting access
//func YouTubeGET(w http.ResponseWriter, r *http.Request) {
//
//}
//
//func printPlaylistsListResults(response *youtube.PlaylistListResponse) {
//	for _, item := range response.Items {
//		fmt.Println(item.Id, ": ", item.Snippet.Title)
//	}
//}
//
//func playlistsListMine(service *youtube.Service, part string, mine bool, maxResults int64, onBehalfOfContentOwner string, onBehalfOfContentOwnerChannel string) {
//	call := service.Playlists.List(part)
//	if mine {
//		call = call.Mine(mine)
//	}
//	if maxResults != 0 {
//		call = call.MaxResults(maxResults)
//	}
//	if onBehalfOfContentOwner != "" {
//		call = call.OnBehalfOfContentOwner(onBehalfOfContentOwner)
//	}
//	if onBehalfOfContentOwnerChannel != "" {
//		call = call.OnBehalfOfContentOwnerChannel(onBehalfOfContentOwnerChannel)
//	}
//	response, err := call.Do()
//	handleError(err, "")
//	printPlaylistsListResults(response)
//}
