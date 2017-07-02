package model

import (
	"time"

	"github.com/nclandrei/YTSync/shared/database"

	"gopkg.in/mgo.v2/bson"
)

// *****************************************************************************
// User
// *****************************************************************************

// User table contains the information for each user
type User struct {
	ObjectID  		bson.ObjectId `bson:"_id"`
	ID               uint32        `db:"id" bson:"id,omitempty"` // Don't use Id, use UserID() instead for consistency with MongoDB
	Email     		string        `db:"email" bson:"email"`
	Password  		string        `db:"password" bson:"password"`
	StatusID  		uint8         `db:"status_id" bson:"status_id"`
	RefreshToken		string        `db:"refresh_token" bson:"refresh_token"`
	CreatedAt 		time.Time     `db:"created_at" bson:"created_at"`
	UpdatedAt 		time.Time     `db:"updated_at" bson:"updated_at"`
	Deleted   		uint8         `db:"deleted" bson:"deleted"`
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

func (u *User) UserRefreshToken() string {
	return u.RefreshToken
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
