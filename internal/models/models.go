package models

import (
	"database/sql"
	"errors"
	"time"

	"github.com/Danik14/library/internal/data"
	uuid "github.com/satori/go.uuid"
)

// Define a custom ErrRecordNotFound error. We'll return this from our Get() method when
// looking up a book that doesn't exist in our database.
var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

// Create a Models struct which wraps the BookModel. We'll add other models to this,
// like a UserModel and PermissionModel, as our build progresses.
// type Models struct {
// 	Users       UserModel
// 	Permissions PermissionModel
// 	Books       BookModel
// 	Tokens      TokenModel
// }

type Models struct {
	Books interface {
		Insert(book *Book) error
		Get(id uuid.UUID) (*Book, error)
		Update(book *Book) error
		Delete(id uuid.UUID) error
		GetAll(title string, author string, genres []string, filters data.Filters) ([]*Book, data.Metadata, error)
	}
	Users interface {
		Insert(user *User) error
		GetAll(firstName string, lastName string, email string, filters data.Filters) ([]*User, data.Metadata, error)
		Get(id uuid.UUID) (*User, error)
		GetByEmail(email string) (*User, error)
		Update(user *User) error
		Delete(id uuid.UUID) error
		GetForToken(tokenScope, tokenPlaintext string) (*User, error)
	}
	Tokens interface {
		DeleteAllForUser(scope string, userID uuid.UUID) error
		Insert(token *Token) error
		New(userID uuid.UUID, ttl time.Duration, scope string) (*Token, error)
	}
	Permissions interface {
		GetAllForUser(userID uuid.UUID) (Permissions, error)
		AddForUser(userID uuid.UUID, codes ...string) error
	}
}

// For ease of use, we also add a New() method which returns a Models struct containing
// the initialized BookModel.
func NewModels(db *sql.DB) Models {
	return Models{
		Users:       UserModel{DB: db},
		Permissions: PermissionModel{DB: db},
		Books:       BookModel{DB: db},
		Tokens:      TokenModel{DB: db},
	}
}

func NewMockModels() Models {
	return Models{
		Books: MockBookModel{},
		// Users:       MockUserModel{},
		// Tokens:      MockTokenModel{},
		// Permissions: MockPermissionModel{},
	}
}
