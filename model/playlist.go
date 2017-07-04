package model

import (
	"time"

	"github.com/nclandrei/YTSync/shared/database"

	"gopkg.in/mgo.v2/bson"
)

// *****************************************************************************
// Playlist
// *****************************************************************************

// Playlist table contains the information for each playlist per user
type Playlist struct {
	ObjectID  bson.ObjectId `bson:"_id"`
	ID        uint32        `db:"id" bson:"id,omitempty"`
	Content   string        `db:"content" bson:"content"`
	UserID    bson.ObjectId `bson:"user_id"`
	UID       uint32        `db:"user_id" bson:"userid,omitempty"`
	CreatedAt time.Time     `db:"created_at" bson:"created_at"`
	UpdatedAt time.Time     `db:"updated_at" bson:"updated_at"`
	Deleted   uint8         `db:"deleted" bson:"deleted"`
}

// PlaylistID returns the note id
func (u *Playlist) PlaylistID() string {
	r := ""
	r = u.ObjectID.Hex()
	return r
}

// PlaylistByID gets note by ID
func PlaylistByID(userID string, playlistID string) (Playlist, error) {
	var err error

	result := Playlist{}

	if database.CheckConnection() {
		// Create a copy of mongo
		session := database.Mongo.Copy()
		defer session.Close()
		c := session.DB(database.ReadConfig().MongoDB.Database).C("playlist")

		// Validate the object id
		if bson.IsObjectIdHex(playlistID) {
			err = c.FindId(bson.ObjectIdHex(playlistID)).One(&result)
			if result.UserID != bson.ObjectIdHex(userID) {
				result = Playlist{}
				err = ErrUnauthorized
			}
		} else {
			err = ErrNoResult
		}
	} else {
		err = ErrUnavailable
	}

	return result, standardizeError(err)
}

// PlaylistByUserID gets all playlists for a user
func PlaylistByUserID(userID string) ([]Playlist, error) {
	var err error

	var result []Playlist

	if database.CheckConnection() {
		// Create a copy of mongo
		session := database.Mongo.Copy()
		defer session.Close()
		c := session.DB(database.ReadConfig().MongoDB.Database).C("playlist")

		// Validate the object id
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

// NoteCreate creates a note
func PlaylistCreate(content string, userID string) error {
	var err error

	now := time.Now()

	if database.CheckConnection() {
		// Create a copy of mongo
		session := database.Mongo.Copy()
		defer session.Close()
		c := session.DB(database.ReadConfig().MongoDB.Database).C("playlist")

		playlist := &Playlist{
			ObjectID:  bson.NewObjectId(),
			Content:   content,
			UserID:    bson.ObjectIdHex(userID),
			CreatedAt: now,
			UpdatedAt: now,
			Deleted:   0,
		}
		err = c.Insert(playlist)
	} else {
		err = ErrUnavailable
	}

	return standardizeError(err)
}

// NoteUpdate updates a note
func PlaylistUpdate(content string, userID string, playlistID string) error {
	var err error

	now := time.Now()

	if database.CheckConnection() {
		// Create a copy of mongo
		session := database.Mongo.Copy()
		defer session.Close()
		c := session.DB(database.ReadConfig().MongoDB.Database).C("playlist")
		var playlist Playlist
		playlist, err = PlaylistByID(userID, playlistID)
		if err == nil {
			// Confirm the owner is attempting to modify the note
			if playlist.UserID.Hex() == userID {
				playlist.UpdatedAt = now
				playlist.Content = content
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

// PlaylistDelete deletes a note
func PlaylistDelete(userID string, playlistID string) error {
	var err error

	if database.CheckConnection() {
		// Create a copy of mongo
		session := database.Mongo.Copy()
		defer session.Close()
		c := session.DB(database.ReadConfig().MongoDB.Database).C("playlist")

		var playlist Playlist
		playlist, err = PlaylistByID(userID, playlistID)
		if err == nil {
			// Confirm the owner is attempting to modify the note
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
