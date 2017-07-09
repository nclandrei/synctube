package model

import (
	"github.com/nclandrei/YTSync/shared/database"

	"gopkg.in/mgo.v2/bson"
)

// *****************************************************************************
// Video
// *****************************************************************************

// Video table contains the information for each Video per user
type Video struct {
	ObjectID   bson.ObjectId `bson:"_id"`
	ID         string        `db:"id" bson:"id,omitempty"`
	Title      string        `db:"title" bson:"title"`
	PlaylistID string        `bson:"playlist_id"`
}

// VideoByID gets video in a given playlist
func VideoByID(videoID string, playlistID string) (Video, error) {
	var err error

	result := Video{}

	if database.CheckConnection() {
		// Create a copy of mongo
		session := database.Mongo.Copy()
		defer session.Close()
		c := session.DB(database.ReadConfig().MongoDB.Database).C("video")
		err = c.Find(bson.M{"$and": []bson.M{bson.M{"id": videoID}, bson.M{"playlist_id": playlistID}}}).One(&result)

		if err != nil {
			result = Video{}
			err = ErrUnauthorized
		}
	} else {
		err = ErrUnavailable
	}

	return result, standardizeError(err)
}

// VideoByPlaylistID gets all Videos for a user
func VideosByPlaylistID(playlistID string) ([]Video, error) {
	var err error

	var result []Video

	if database.CheckConnection() {
		// Create a copy of mongo
		session := database.Mongo.Copy()
		defer session.Close()
		c := session.DB(database.ReadConfig().MongoDB.Database).C("video")
		err = c.Find(bson.M{"playlist_id": playlistID}).All(&result)
	} else {
		err = ErrUnavailable
	}
	return result, standardizeError(err)
}

// VideoCreate creates a video
func VideoCreate(id string, title string, playlistID string) error {
	var err error

	if database.CheckConnection() {
		// Create a copy of mongo
		session := database.Mongo.Copy()
		defer session.Close()
		c := session.DB(database.ReadConfig().MongoDB.Database).C("video")

		Video := &Video{
			ObjectID:   bson.NewObjectId(),
			ID:         id,
			Title:      title,
			PlaylistID: playlistID,
		}
		err = c.Insert(Video)
	} else {
		err = ErrUnavailable
	}

	return standardizeError(err)
}

// VideoDelete deletes a video
func VideoDelete(videoID string, playlistID string) error {
	var err error

	if database.CheckConnection() {
		// Create a copy of mongo
		session := database.Mongo.Copy()
		defer session.Close()
		c := session.DB(database.ReadConfig().MongoDB.Database).C("video")

		var video Video
		video, err = VideoByID(videoID, playlistID)

		if err == nil {
			// Confirm the owner is attempting to modify the note
			if video.ID == videoID && video.PlaylistID == playlistID {
				err = c.Remove(bson.M{"id": videoID})
			} else {
				err = ErrUnauthorized
			}
		}
	} else {
		err = ErrUnavailable
	}

	return standardizeError(err)
}
