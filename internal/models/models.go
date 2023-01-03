package models

import (
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
)

// Define a custom ErrRecordNotFound error. We'll return this from our Get() method when
// looking up a movie that doesn't exist in our database.
var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

// Create a Models struct which wraps the MovieModel. We'll add other models to this,
// like a UserModel and PermissionModel, as our build progresses.
type Models struct {
	Users UserModel
	Books BookModel
}

// For ease of use, we also add a New() method which returns a Models struct containing
// the initialized MovieModel.
func NewModels(users *mongo.Collection, books *mongo.Collection) Models {
	return Models{
		Users: UserModel{DB: users},
		Books: BookModel{DB: books},
	}
}
