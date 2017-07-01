package model

import (
	"time"

	"github.com/nclandrei/YTSync/shared/database"

	"gopkg.in/mgo.v2/bson"
)

// *****************************************************************************
// Note
// *****************************************************************************

// Note table contains the information for each note
type Note struct {
	ObjectID  bson.ObjectId `bson:"_id"`
	ID        uint32        `db:"id" bson:"id,omitempty"` // Don't use Id, use NoteID() instead for consistency with MongoDB
	Content   string        `db:"content" bson:"content"`
	UserID    bson.ObjectId `bson:"user_id"`
	UID       uint32        `db:"user_id" bson:"userid,omitempty"`
	CreatedAt time.Time     `db:"created_at" bson:"created_at"`
	UpdatedAt time.Time     `db:"updated_at" bson:"updated_at"`
	Deleted   uint8         `db:"deleted" bson:"deleted"`
}

// NoteID returns the note id
func (u *Note) NoteID() string {
	r := ""
	r = u.ObjectID.Hex()
	return r
}

// NoteByID gets note by ID
func NoteByID(userID string, noteID string) (Note, error) {
	var err error

	result := Note{}

	if database.CheckConnection() {
		// Create a copy of mongo
		session := database.Mongo.Copy()
		defer session.Close()
		c := session.DB(database.ReadConfig().MongoDB.Database).C("note")

		// Validate the object id
		if bson.IsObjectIdHex(noteID) {
			err = c.FindId(bson.ObjectIdHex(noteID)).One(&result)
			if result.UserID != bson.ObjectIdHex(userID) {
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

// NotesByUserID gets all notes for a user
func NotesByUserID(userID string) ([]Note, error) {
	var err error

	var result []Note

	if database.CheckConnection() {
		// Create a copy of mongo
		session := database.Mongo.Copy()
		defer session.Close()
		c := session.DB(database.ReadConfig().MongoDB.Database).C("note")

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
func NoteCreate(content string, userID string) error {
	var err error

	now := time.Now()

	if database.CheckConnection() {
		// Create a copy of mongo
		session := database.Mongo.Copy()
		defer session.Close()
		c := session.DB(database.ReadConfig().MongoDB.Database).C("note")

		note := &Note{
			ObjectID:  bson.NewObjectId(),
			Content:   content,
			UserID:    bson.ObjectIdHex(userID),
			CreatedAt: now,
			UpdatedAt: now,
			Deleted:   0,
		}
		err = c.Insert(note)
	} else {
		err = ErrUnavailable
	}

	return standardizeError(err)
}

// NoteUpdate updates a note
func NoteUpdate(content string, userID string, noteID string) error {
	var err error

	now := time.Now()

	if database.CheckConnection() {
		// Create a copy of mongo
		session := database.Mongo.Copy()
		defer session.Close()
		c := session.DB(database.ReadConfig().MongoDB.Database).C("note")
		var note Note
		note, err = NoteByID(userID, noteID)
		if err == nil {
			// Confirm the owner is attempting to modify the note
			if note.UserID.Hex() == userID {
				note.UpdatedAt = now
				note.Content = content
				err = c.UpdateId(bson.ObjectIdHex(noteID), &note)
			} else {
				err = ErrUnauthorized
			}
		}
	} else {
		err = ErrUnavailable
	}

	return standardizeError(err)
}

// NoteDelete deletes a note
func NoteDelete(userID string, noteID string) error {
	var err error

	if database.CheckConnection() {
		// Create a copy of mongo
		session := database.Mongo.Copy()
		defer session.Close()
		c := session.DB(database.ReadConfig().MongoDB.Database).C("note")

		var note Note
		note, err = NoteByID(userID, noteID)
		if err == nil {
			// Confirm the owner is attempting to modify the note
			if note.UserID.Hex() == userID {
				err = c.RemoveId(bson.ObjectIdHex(noteID))
			} else {
				err = ErrUnauthorized
			}
		}
	} else {
		err = ErrUnavailable
	}

	return standardizeError(err)
}
