package model

import (
	"time"

	"github.com/nclandrei/synctube/shared/database"

	"golang.org/x/oauth2"
	"gopkg.in/mgo.v2/bson"
)

// *****************************************************************************
// User
// *****************************************************************************

// User table contains the information for each user
type User struct {
	ObjectID  bson.ObjectId `bson:"_id"`
	ID        uint32        `db:"id" bson:"id,omitempty"`
	Email     string        `db:"email" bson:"email"`
	Password  string        `db:"password" bson:"password"`
	StatusID  uint8         `db:"status_id" bson:"status_id"`
	Token     oauth2.Token  `db:"token" bson:"token"`
	CreatedAt time.Time     `db:"created_at" bson:"created_at"`
	UpdatedAt time.Time     `db:"updated_at" bson:"updated_at"`
	Deleted   uint8         `db:"deleted" bson:"deleted"`
	LastSync  time.Time     `db:"last_sync" bson:"last_sync"`
}

// UserStatus table contains every possible user status (active/inactive)
type UserStatus struct {
	ID        uint8     `db:"id" bson:"id"`
	Status    string    `db:"status" bson:"status"`
	CreatedAt time.Time `db:"created_at" bson:"created_at"`
	UpdatedAt time.Time `db:"updated_at" bson:"updated_at"`
	Deleted   uint8     `db:"deleted" bson:"deleted"`
}

// UserID returns the user id
func (u *User) UserID() string {
	return u.ObjectID.Hex()
}

// UserByToken returns the user's token
func UserByToken(userID string) (User, error) {
	var err error
	result := User{}

	if database.CheckConnection() {
		session := database.Mongo.Copy()
		defer session.Close()
		c := session.DB(database.ReadConfig().MongoDB.Database).C("user")
		err = c.Find(bson.M{"_id": bson.ObjectIdHex(userID)}).One(&result)
	} else {
		err = ErrUnavailable
	}
	return result, standardizeError(err)
}

// UpdateUserToken updates the record with a new
func UpdateUserToken(userID string, token oauth2.Token) error {
	var err error
	if database.CheckConnection() {
		session := database.Mongo.Copy()
		defer session.Close()
		c := session.DB(database.ReadConfig().MongoDB.Database).C("user")
		err = c.Update(bson.M{"_id": bson.ObjectIdHex(userID)}, bson.M{"$set": bson.M{"token": token}})
	} else {
		err = ErrUnavailable
	}
	return err
}

// UserByEmail gets user information from email
func UserByEmail(email string) (User, error) {
	var err error

	result := User{}

	if database.CheckConnection() {
		session := database.Mongo.Copy()
		defer session.Close()
		c := session.DB(database.ReadConfig().MongoDB.Database).C("user")
		err = c.Find(bson.M{"email": email}).One(&result)
	} else {
		err = ErrUnavailable
	}
	return result, standardizeError(err)
}

// UserCreate creates user
func UserCreate(email, password string) error {
	var err error

	now := time.Now()

	if database.CheckConnection() {
		session := database.Mongo.Copy()
		defer session.Close()
		c := session.DB(database.ReadConfig().MongoDB.Database).C("user")

		user := &User{
			ObjectID:  bson.NewObjectId(),
			Email:     email,
			Password:  password,
			StatusID:  1,
			Token:     oauth2.Token{},
			CreatedAt: now,
			UpdatedAt: now,
			Deleted:   0,
		}
		err = c.Insert(user)
	} else {
		err = ErrUnavailable
	}
	return standardizeError(err)
}

// UserUpdateLastSync updates last synchronization timestamp for current user
func UserUpdateLastSync(userID string, timestamp time.Time) error {
	var err error
	if database.CheckConnection() {
		session := database.Mongo.Copy()
		defer session.Close()
		c := session.DB(database.ReadConfig().MongoDB.Database).C("user")
		err = c.Update(bson.M{"_id": bson.ObjectIdHex(userID)}, bson.M{"$set": bson.M{"last_sync": timestamp}})
	} else {
		err = ErrUnavailable
	}
	return err
}
