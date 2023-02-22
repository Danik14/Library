package models

import (
	"time"

	"github.com/Danik14/library/internal/validator"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	ID             primitive.ObjectID            `json:"_id" bson:"_id"`
	CreatedAt      time.Time                     `json:"-"`
	FirstName      string                        `json:"firstName"`
	LastName       string                        `json:"lastName"`
	Email          string                        `json:"email"`
	HashedPassword string                        `json:"-"`
	Marks          map[primitive.ObjectID]string `json:"marks"`
	DOB            time.Time                     `json:"dob"` // date of birth
}

type UserModel struct {
	DB *mongo.Collection
}

func ValidateUser(v *validator.Validator, user *User) {

	if user.FirstName != "" {
		v.Check(len(user.FirstName) <= 500, "firstName", "must not be more than 500 bytes long")
	}
	if user.LastName != "" {
		v.Check(len(user.LastName) <= 500, "lastName", "must not be more than 500 bytes long")
	}
	if user.Email != "" {
		v.Check(len(user.Email) <= 500, "email", "must not be more than 500 bytes long")
		v.Check(validator.Matches(user.Email, validator.EmailRX), "email", "email format is not correct")
	}
	if user.HashedPassword != "" {
		// v.Check(user.HashedPassword != "", "password", "must be provided")
		v.Check(len(user.HashedPassword) <= 500, "password", "must not be more than 500 bytes long")
	}
	if (user.DOB != time.Time{}) {
		// v.Check(user.DOB != time.Time{}, "dob", "must be provided")
		year := int32(user.DOB.Year())
		month := int32(user.DOB.Month())
		day := int32(user.DOB.Day())
		v.Check(year >= 1900, "dob", "must be greater than 1900")
		v.Check(year <= int32(time.Now().Year()), "dob", "year must not be in the future")
		if year == int32(time.Now().Year()) {
			v.Check(month <= int32(time.Now().Month()), "dob", "month must not be in the future")
			if month == int32(time.Now().Month()) {
				v.Check(day <= int32(time.Now().Day()), "dob", "day must not be in the future")
			}
		}
	}
}
