package user

import (
	"context"
	"fmt"
	"time"

	"github.com/amankapoor/goth"
	"github.com/amankapoor/placementpal/internal/platform/db"
	"github.com/amankapoor/placementpal/internal/platform/web"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const usersCollection = "users_students"

// List retrieves a list of existing users from the database.
func List(ctx context.Context, dbConn *db.DB) ([]User, error) {
	u := []User{}

	f := func(collection *mgo.Collection) error {
		return collection.Find(nil).All(&u)
	}
	if err := dbConn.MGOExecute(ctx, usersCollection, f); err != nil {
		return nil, errors.Wrap(err, "db.users.find()")
	}

	return u, nil
}

// Retrieve gets the specified user from the database.
func Retrieve(ctx context.Context, dbConn *db.DB, googleID string) (*User, error) {

	q := bson.M{"google_id": googleID}

	var u *User
	f := func(collection *mgo.Collection) error {
		return collection.Find(q).One(&u)
	}
	if err := dbConn.MGOExecute(ctx, usersCollection, f); err != nil {
		if err == mgo.ErrNotFound {
			return nil, err
		}
		return nil, errors.Wrap(err, fmt.Sprintf("db.users.find(%s)", db.Query(q)))
	}

	return u, nil
}

// Create inserts a new user into the database.
func Create(ctx context.Context, dbConn *db.DB, cu CreateUser) (*User, error) {
	now := time.Now()

	u := User{
		ID:              bson.NewObjectId().Hex(),
		UserType:        Student,
		FirstName:       cu.FirstName,
		LastName:        cu.LastName,
		Email:           cu.Email,
		GoogleID:        cu.GoogleID,
		AvatarURL:       cu.AvatarURL,
		Batch:           cu.Batch,
		Degree:          cu.Degree,
		College:         cu.College,
		CollegeLocation: cu.CollegeLocation,
		DateCreated:     &now,
		DateModified:    &now,
	}

	f := func(collection *mgo.Collection) error {
		return collection.Insert(u)
	}
	if err := dbConn.MGOExecute(ctx, usersCollection, f); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("db.users.insert(%s)", db.Query(u)))
	}

	return &u, nil
}

// Update replaces a user document in the database.
func Update(ctx context.Context, dbConn *db.DB, userID string, gu goth.User) error {
	if !bson.IsObjectIdHex(userID) {
		return errors.Wrap(web.ErrInvalidID, "check objectid")
	}

	//now := time.Now()
	//gu.DateModified = &now
	// for _, cua := range cu.Addresses {
	// 	cua.DateModified = &now
	// }

	q := bson.M{"user_id": userID}
	m := bson.M{"$set": gu}

	f := func(collection *mgo.Collection) error {
		return collection.Update(q, m)
	}
	if err := dbConn.MGOExecute(ctx, usersCollection, f); err != nil {
		if err == mgo.ErrNotFound {
			return web.ErrNotFound
		}
		return errors.Wrap(err, fmt.Sprintf("db.users.update(%s, %s)", db.Query(q), db.Query(m)))
	}

	return nil
}

// Delete removes a user from the database.
func Delete(ctx context.Context, dbConn *db.DB, userID string) error {
	if !bson.IsObjectIdHex(userID) {
		return errors.Wrapf(web.ErrInvalidID, "bson.IsObjectIdHex: %s", userID)
	}

	q := bson.M{"user_id": userID}

	f := func(collection *mgo.Collection) error {
		return collection.Remove(q)
	}
	if err := dbConn.MGOExecute(ctx, usersCollection, f); err != nil {
		if err == mgo.ErrNotFound {
			return web.ErrNotFound
		}
		return errors.Wrap(err, fmt.Sprintf("db.users.remove(%s)", db.Query(q)))
	}

	return nil
}
