package models

import (
	"database/sql"
	"errors"
	"time"

	"github.com/Danik14/library/internal/validator"
	"github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
)

type User struct {
	ID             uuid.UUID `json:"id"`
	CreatedAt      time.Time `json:"-"`
	FirstName      string    `json:"firstName"`
	LastName       string    `json:"lastName"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"-"`
	DOB            time.Time `json:"dob"` // date of birth
	Version        int32     `json:"version"`
}

type UserModel struct {
	DB *sql.DB
}

func (u UserModel) Insert(user *User) error {
	// return &User{CreatedAt: time.Now(), FirstName: firstName, LastName: lastName, Email: email, HashedPassword: password, DOB: dob, Version: version}, nil
	// Define the SQL query for inserting a new record in the movies table and returning
	// the system-generated data.
	query := `INSERT INTO users (firstname, lastname, email, hashedpassword, dob) VALUES ($1, $2, $3, $4, $5) RETURNING id, createdAt, version;`
	// Create an args slice containing the values for the placeholder parameters from
	// the movie struct. Declaring this slice immediately next to our SQL query helps to
	// make it nice and clear *what values are being used where* in the query.

	// hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.HashedPassword), 12)
	// if err != nil {
	// 	return nil, err
	// }
	args := []any{user.FirstName, user.LastName, user.Email, user.HashedPassword, pq.FormatTimestamp(user.DOB)}
	// Use the QueryRow() method to execute the SQL query on our connection pool,
	// passing in the args slice as a variadic parameter and scanning the system-
	// generated id, created_at and version values into the movie struct.
	return u.DB.QueryRow(query, args...).Scan(&user.ID, &user.CreatedAt, &user.Version)
}

func (u UserModel) Get(id uuid.UUID) (*User, error) {
	// if id < 1 {
	// 	return nil, ErrRecordNotFound
	// }
	// Define the SQL query for retrieving the movie data.
	query := `SELECT id, createdAt, firstName, lastName, email, hashedPassword, dob, version FROM users WHERE id = $1`
	// Declare a Movie struct to hold the data returned by the query.
	var user User // Execute the query using the QueryRow() method, passing in the provided id value
	// as a placeholder parameter, and scan the response data into the fields of the
	// Movie struct. Importantly, notice that we need to convert the scan target for the
	// genres column using the pq.Array() adapter function again.
	err := u.DB.QueryRow(query, id).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.HashedPassword,
		&user.DOB,
		&user.Version,
	)
	// Handle any errors. If there was no matching movie found, Scan() will return
	// a sql.ErrNoRows error. We check for this and return our custom ErrRecordNotFound
	// error instead.
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	// Otherwise, return a pointer to the Movie struct.
	return &user, nil
}

func (u UserModel) Update(user *User) error {
	// Declare the SQL query for updating the record and returning the new version
	// number.
	query := `
	UPDATE users
	SET firstName = $1, lastName = $2, email = $3, hashedPassword = $4, dob = $5, version = version + 1
	WHERE id = $6
	RETURNING version`
	// Create an args slice containing the values for the placeholder parameters.
	args := []any{
		user.FirstName,
		user.LastName,
		user.Email,
		user.HashedPassword,
		user.DOB,
		user.ID,
	}
	// Use the QueryRow() method to execute the query, passing in the args slice as a
	// variadic parameter and scanning the new version value into the movie struct.
	return u.DB.QueryRow(query, args...).Scan(&user.Version)
}

func (u UserModel) Delete(id uuid.UUID) error {
	// Return an ErrRecordNotFound error if the movie ID is less than 1.
	// if id < 1 {
	// 	return ErrRecordNotFound
	// }
	// Construct the SQL query to delete the record.
	query := `DELETE FROM users WHERE id = $1`
	// Execute the SQL query using the Exec() method, passing in the id variable as
	// the value for the placeholder parameter. The Exec() method returns a sql.Result
	// object.
	result, err := u.DB.Exec(query, id)
	if err != nil {
		return err
	}
	// Call the RowsAffected() method on the sql.Result object to get the number of rows
	// affected by the query.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	// If no rows were affected, we know that the movies table didn't contain a record
	// with the provided ID at the moment we tried to delete it. In that case we
	// return an ErrRecordNotFound error.
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
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
