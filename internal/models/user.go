package models

import (
	"fmt"
	"time"

	"github.com/Danik14/library/internal/validator"
)

type User struct {
	ID             int64     `json:"id"`
	CreatedAt      time.Time `json:"-"`
	FirstName      string    `json:"firstName"`
	LastName       string    `json:"lastName"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"-"`
	DOB            time.Time `json:"dob"` // date of birth
	Version        int32     `json:"version"`
}

func NewUser(firstName string, lastName string, email string, password string, dob time.Time, version int32) (*User, error) {
	// hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	// if err != nil {
	// 	return nil, err
	// }
	return &User{CreatedAt: time.Now(), FirstName: firstName, LastName: lastName, Email: email, HashedPassword: password, DOB: dob, Version: version}, nil
}

func GetUser(id int64) (*User, error) {
	//stmt := `SELECT * FROM users WHERE userId = $1`;
	//returnedUser := &User{}
	//err := m.DB.QueryRow(context.Background(),stmt, id).Scan(&returnedUser.id, &returnedUser.FirstName, &returnedUser.LastName,&returnedUser.Email, &returnedUser.DOB,&returnedUSer.Version)
	// if err != nil {
	// 	return nil, err
	// }
	//return returnedUser, nil
	fmt.Println("Import flask")
	return nil, nil
}

func ValidateUser(v *validator.Validator, user *User) {
	// Use the Check() method to execute our validation checks. This will add the
	// provided key and error message to the errors map if the check does not evaluate
	// to true. For example, in the first line here we "check that the title is not
	// equal to the empty string". In the second, we "check that the length of the title
	// is less than or equal to 500 bytes" and so on.
	v.Check(user.FirstName != "", "firstName", "must be provided")
	v.Check(len(user.FirstName) <= 500, "firstName", "must not be more than 500 bytes long")

	v.Check(user.LastName != "", "lastName", "must be provided")
	v.Check(len(user.LastName) <= 500, "lastName", "must not be more than 500 bytes long")

	v.Check(user.Email != "", "email", "must be provided")
	v.Check(len(user.Email) <= 500, "email", "must not be more than 500 bytes long")
	v.Check(validator.Matches(user.Email, validator.EmailRX), "email", "email format is not correct")

	v.Check(user.HashedPassword != "", "password", "must be provided")
	v.Check(len(user.HashedPassword) <= 500, "password", "must not be more than 500 bytes long")

	v.Check(user.DOB != time.Time{}, "dob", "must be provided")
	v.Check(user.DOB.Year() >= 1900, "dob", "must be greater than 1900")
	v.Check(int32(user.DOB.Year()) <= int32(time.Now().Year()), "dob", "year must not be in the future")
	v.Check(int32(user.DOB.Month()) <= int32(time.Now().Month()), "dob", "month must not be in the future")
	v.Check(int32(user.DOB.Day()) <= int32(time.Now().Day()), "dob", "day must not be in the future")
}
