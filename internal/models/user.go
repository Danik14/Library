package models

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             int64     `json:"id"`
	CreatedAt      time.Time `json:"-"`
	FirstName      string    `json:"firstName"`
	LastName       string    `json:"lastName"`
	Email          string    `json:"email"`
	HashedPassword []byte    `json:"-"`
	DOB            time.Time `json:"dob"` // date of birth
	Version        int32     `json:"version"`
}

func NewUser(firstName string, lastName string, email string, password string, dob time.Time, version int32) (*User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return nil, err
	}
	return &User{CreatedAt: time.Now(), FirstName: firstName, LastName: lastName, Email: email, HashedPassword: hash, DOB: dob, Version: version}, nil
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
