package model

import (
	"github.com/nclandrei/YTSync/shared/database"

	"gopkg.in/mgo.v2/bson"
)

// *****************************************************************************
// Playlist
// *****************************************************************************

// Playlist table contains the information for each playlist per user
type Playlist struct {
	ObjectID bson.ObjectId `bson:"_id"`
	ID       string        `db:"id" bson:"id"`
	Title    string        `db:"title" bson:"title"`
	UserID   bson.ObjectId `bson:"user_id"`
}

// PlaylistID - returns the playlist id
func (u *Playlist) PlaylistID() string {
	r := ""
	r = u.ObjectID.Hex()
	return r
}

// PlaylistByID - gets playlist by its YouTube-specific ID and the user's ID
func PlaylistByID(playlistID string, userID string) (Playlist, error) {
	var err error

	result := Playlist{}

	if database.CheckConnection() {
		session := database.Mongo.Copy()
		defer session.Close()
		c := session.DB(database.ReadConfig().MongoDB.Database).C("playlist")
		err = c.Find(bson.M{"id": playlistID}).One(&result)
	} else {
		err = ErrUnavailable
	}

	return result, standardizeError(err)
}

// PlaylistByUserID - gets all playlists for a user
func PlaylistByUserID(userID string) ([]Playlist, error) {
	var err error

	var result []Playlist

	if database.CheckConnection() {
		session := database.Mongo.Copy()
		defer session.Close()
		c := session.DB(database.ReadConfig().MongoDB.Database).C("playlist")

		if bson.IsObjectIdHex(userID) {
			err = c.Find(bson.M{"user_id": bson.ObjectIdHex(userID)}).All(&result)
		} else {
			err = ErrNoResult
		}
	} else {
		err = ErrUnavailable
	}

	return result, standardizeError(err)
}

// PlaylistCreate - creates a playlist given the YT-specific ID, its title and the user who owns it
func PlaylistCreate(id string, title string, userID string) error {
	var err error

	if database.CheckConnection() {
		// Create a copy of mongo
		session := database.Mongo.Copy()
		defer session.Close()
		c := session.DB(database.ReadConfig().MongoDB.Database).C("playlist")

		playlist := &Playlist{
			ObjectID: bson.NewObjectId(),
			ID:       id,
			Title:    title,
			UserID:   bson.ObjectIdHex(userID),
		}
		err = c.Insert(playlist)
	} else {
		err = ErrUnavailable
	}

	return standardizeError(err)
}

// PlaylistUpdate - updates a playlist given its content, user ID and the playlist's ID
func PlaylistUpdate(content string, userID string, playlistID string) error {
	var err error

	if database.CheckConnection() {
		session := database.Mongo.Copy()
		defer session.Close()
		c := session.DB(database.ReadConfig().MongoDB.Database).C("playlist")
		var playlist Playlist
		playlist, err = PlaylistByID(userID, playlistID)
		if err == nil {
			if playlist.UserID.Hex() == userID {
				playlist.Title = content
				err = c.UpdateId(bson.ObjectIdHex(playlistID), &playlist)
			} else {
				err = ErrUnauthorized
			}
		}
	} else {
		err = ErrUnavailable
	}

	return standardizeError(err)
}

// PlaylistDelete - deletes a playlist given its user's ID and the playlist's ID
func PlaylistDelete(userID string, playlistID string) error {
	var err error

	if database.CheckConnection() {
		session := database.Mongo.Copy()
		defer session.Close()
		c := session.DB(database.ReadConfig().MongoDB.Database).C("playlist")

		var playlist Playlist
		playlist, err = PlaylistByID(userID, playlistID)
		if err == nil {
			if playlist.UserID.Hex() == userID {
				err = c.RemoveId(bson.ObjectIdHex(playlistID))
			} else {
				err = ErrUnauthorized
			}
		}
	} else {
		err = ErrUnavailable
	}

	return standardizeError(err)
}
