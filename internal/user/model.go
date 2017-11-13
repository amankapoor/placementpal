// All material is licensed under the Apache License Version 2.0, January 2004
// http://www.apache.org/licenses/LICENSE-2.0

package user

import "time"

const (
	Student = iota
)

// CreateUser contains information needed to create or update a user.
type CreateUser struct {
	UserType        int        `bson:"type" validate:"required"`
	FirstName       string     `bson:"first_name" validate:"required"`
	LastName        string     `bson:"last_name" validate:"required"`
	Email           string     `bson:"email" validate:"required"`
	GoogleID        string     `bson:"google_id" validate:"required"`
	AvatarURL       string     `bson:"avatar" validate:"required"`
	Batch           string     `bson:"batch" validate:"required"`
	Degree          string     `bson:"degree" validate:"required"`
	College         string     `bson:"college" validate:"required"`
	CollegeLocation string     `bson:"college_location" validate:"required"`
	DateModified    *time.Time `bson:"date_modified" validate:"required"`
}

// User contains information about a user.
type User struct {
	ID              string     `bson:"_id" json:"_id"`
	UserType        int        `bson:"type" json:"type"`
	FirstName       string     `bson:"first_name" json:"first_name"`
	LastName        string     `bson:"last_name" json:"last_name"`
	Email           string     `bson:"email" json:"email"`
	GoogleID        string     `bson:"google_id" json:"google_id"`
	AvatarURL       string     `bson:"avatar" json:"avatar"`
	Batch           string     `bson:"batch" json:"batch"`
	Degree          string     `bson:"degree" json:"degree"`
	College         string     `bson:"college" json:"college"`
	CollegeLocation string     `bson:"college_location" json:"college_location"`
	DateModified    *time.Time `bson:"date_modified" json:"date_modified"`
	DateCreated     *time.Time `bson:"date_created" json:"date_created"`
}
