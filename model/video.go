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
    ObjectID    bson.ObjectId   `bson:"_id"`
    ID          string          `db:"id" bson:"id,omitempty"`
    Title       string          `db:"content" bson:"content"`
    URL         string          `db:"url" bson:"url"`
    PlaylistID  bson.ObjectId   `bson:"playlist_id"`
}

// VideoID returns the video id
func (u *Video) VideoID() string {
    r := ""
    r = u.ObjectID.Hex()
    return r
}

// VideoByID gets note by ID
func VideoByID(videoID string, playlistID string) (Video, error) {
	var err error

	result := Video{}

	if database.CheckConnection() {
		// Create a copy of mongo
		session := database.Mongo.Copy()
		defer session.Close()
		c := session.DB(database.ReadConfig().MongoDB.Database).C("video")

		// Validate the object id
		if bson.IsObjectIdHex(videoID) {
			err = c.Find(bson.M{"$and" : [bson.M{"playlist_id" : playlistID}, bson.M{"id" : videoID}]}).One(&result)
			if err != nil {
				result = Note{}
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

// VideoByUserID gets all Videos for a user
func VideoByPlaylistID(playlistID string) ([]Video, error) {
    var err error

    var result []Video

    if database.CheckConnection() {
        // Create a copy of mongo
        session := database.Mongo.Copy()
        defer session.Close()
        c := session.DB(database.ReadConfig().MongoDB.Database).C("video")

        // Validate the object id
        if bson.IsObjectIdHex(playlistID) {
            err = c.Find(bson.M{"playlist_id": bson.ObjectIdHex(playlistID)}).All(&result)
        } else {
            err = ErrNoResult
        }
    } else {
        err = ErrUnavailable
    }

    return result, standardizeError(err)
}

// VideoCreate creates a note
func VideoCreate(id string, title string, url string, playlistID string) error {
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
            URL:        url,
            PlaylistID: bson.ObjectIdHex(playlistID),
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
            if video.UserID.Hex() == userID {
                err = c.RemoveId(bson.ObjectIdHex(videoID))
            } else {
                err = ErrUnauthorized
            }
        }
    } else {
        err = ErrUnavailable
    }

    return standardizeError(err)
}
