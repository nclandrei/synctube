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
	ObjectID	   bson.ObjectId   `bson:"_id"`
    ID          string          `db:"id" bson:"id,omitempty"`
    Title       string          `db:"title" bson:"title,omitempty"`
    URL         string          `db:"url" bson:"url,omitempty"`
    PlaylistID  string          `bson:"playlist_id,omitempty"`
}

// VideoID returns the video object id
func (u *Video) VideoID() string {
    return u.ObjectID.Hex()
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

		// TODO - complete this so that it selects properly
		//err = c.Find(bson.M{"$and" : [bson.M{"playlist_id" : playlistID}, bson.M{"id" : videoID}]}).One(&result)
		err = c.Find(bson.ObjectIdHex(videoID)).One(&result)

		// Validate the object id
		err = c.Find(bson.M{"user_id": videoID}).All(&result)

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

        // Validate the object id
        if bson.IsObjectIdHex(playlistID) {
            err = c.Find(bson.M{"playlist_id": playlistID}).All(&result)
        } else {
            err = ErrNoResult
        }
    } else {
        err = ErrUnavailable
    }

    return result, standardizeError(err)
}

// VideoCreate creates a video
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
            if video.VideoID() == videoID {
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
