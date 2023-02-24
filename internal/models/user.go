package models

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Danik14/library/internal/data"
	"github.com/Danik14/library/internal/validator"
	"github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail = errors.New("duplicate email")
)

type password struct {
	plaintext *string
	hash      []byte
}

var AnonymousUser = &User{}

type User struct {
	ID             uuid.UUID `json:"id"`
	CreatedAt      time.Time `json:"-"`
	FirstName      string    `json:"firstName"`
	LastName       string    `json:"lastName"`
	Email          string    `json:"email"`
	HashedPassword password  `json:"-"`
	DOB            CivilTime `json:"dob"` // date of birth
	Activated      bool      `json:"activated"`
	Version        int32     `json:"version"`
}

func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

// The Set() method calculates the bcrypt hash of a plaintext password, and stores both
// the hash and the plaintext versions in the struct.
func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}
	p.plaintext = &plaintextPassword
	p.hash = hash
	return nil
}

// The Matches() method checks whether the provided plaintext password matches the
// hashed password stored in the struct, returning true if it matches and false
// otherwise.
func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

type UserModel struct {
	DB *sql.DB
}

func (u UserModel) Insert(user *User) error {
	// return &User{CreatedAt: time.Now(), FirstName: firstName, LastName: lastName, Email: email, HashedPassword: password, DOB: dob, Version: version}, nil
	// Define the SQL query for inserting a new record in the movies table and returning
	// the system-generated data.
	query := `INSERT INTO users (firstname, lastname, email, hashedpassword, dob, activated) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, createdAt, version;`
	// Create an args slice containing the values for the placeholder parameters from
	// the movie struct. Declaring this slice immediately next to our SQL query helps to
	// make it nice and clear *what values are being used where* in the query.

	// hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.HashedPassword.plaintext), 12)
	// if err != nil {
	// 	return nil, err
	// }
	args := []any{user.FirstName, user.LastName, user.Email, user.HashedPassword.hash, pq.FormatTimestamp(time.Time(user.DOB)), user.Activated}

	// Create a context with a 3-second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Use QueryRowContext() and pass the context as the first argument.
	err := u.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt, &user.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}
	return nil
}

func (m UserModel) GetByEmail(email string) (*User, error) {
	query := `
	SELECT id, createdAt, firstName, lastName, email, hashedPassword, dob, version, activated FROM users
	WHERE email = $1`
	var user User
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.HashedPassword.hash,
		&user.DOB,
		&user.Version,
		&user.Activated,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}

func (u UserModel) GetAll(firstName string, lastName string, email string, filters data.Filters) ([]*User, data.Metadata, error) {
	// Construct the SQL query to retrieve all movie records.
	query := fmt.Sprintf(`SELECT count(*) OVER(), id, createdAt, firstName, lastName, email, hashedPassword, dob, version, activated FROM users
	WHERE (to_tsvector('simple', firstName) @@ plainto_tsquery('simple', $1) OR $1 = '')
	AND (to_tsvector('simple', lastName) @@ plainto_tsquery('simple', $2) OR $2 = '')
	AND (to_tsvector('simple', email) @@ plainto_tsquery('simple', $3) OR $3 = '')
	ORDER BY %s %s, createdAt DESC
	LIMIT $4 OFFSET $5`, filters.SortColumn(), filters.SortDirection())
	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// As our SQL query now has quite a few placeholder parameters, let's collect the
	// values for the placeholders in a slice. Notice here how we call the limit() and
	// offset() methods on the Filters struct to get the appropriate values for the
	// LIMIT and OFFSET clauses.
	args := []any{firstName, lastName, email, filters.Limit(), filters.Offset()}

	// Use QueryContext() to execute the query. This returns a sql.Rows resultset
	// containing the result.
	rows, err := u.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, data.Metadata{}, err
	}
	// Importantly, defer a call to rows.Close() to ensure that the resultset is closed
	// before GetAll() returns.
	defer rows.Close()
	// Initialize an empty slice to hold the movie data.
	totalRecords := 0
	users := []*User{}
	// Use rows.Next to iterate through the rows in the resultset.
	for rows.Next() {
		// Initialize an empty Movie struct to hold the data for an individual movie.
		var user User
		// Scan the values from the row into the Movie struct. Again, note that we're
		// using the pq.Array() adapter on the genres field here.
		err := rows.Scan(
			&totalRecords,
			&user.ID,
			&user.CreatedAt,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.HashedPassword.hash,
			&user.DOB,
			&user.Version,
			&user.Activated,
		)
		if err != nil {
			return nil, data.Metadata{}, err
		}
		// Add the Movie struct to the slice.
		users = append(users, &user)
	}
	// When the rows.Next() loop has finished, call rows.Err() to retrieve any error
	// that was encountered during the iteration.
	if err = rows.Err(); err != nil {
		return nil, data.Metadata{}, err
	}

	// Generate a Metadata struct, passing in the total record count and pagination
	// parameters from the client.
	metadata := data.CalculateMetadata(totalRecords, filters.Page, filters.PageSize)

	// If everything went OK, then return the slice of movies.
	return users, metadata, nil
}

func (u UserModel) Get(id uuid.UUID) (*User, error) {
	// if id < 1 {
	// 	return nil, ErrRecordNotFound
	// }
	// Define the SQL query for retrieving the movie data.
	query := `SELECT id, createdAt, firstName, lastName, email, hashedPassword, dob, version, activated FROM users WHERE id = $1`

	// Declare a Movie struct to hold the data returned by the query.
	var user User // Execute the query using the QueryRow() method, passing in the provided id value

	// Use the context.WithTimeout() function to create a context.Context which carries a
	// 3-second timeout deadline. Note that we're using the empty context.Background()
	// as the 'parent' context.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// Importantly, use defer to make sure that we cancel the context before the Get()
	// method returns.
	defer cancel()
	// Use the QueryRowContext() method to execute the query, passing in the context
	// with the deadline as the first argument.
	err := u.DB.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.HashedPassword.hash,
		&user.DOB,
		&user.Version,
		&user.Activated,
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
	SET firstName = $1, lastName = $2, email = $3, hashedPassword = $4, dob = $5, version = version + 1, activated=$6
	WHERE id = $7 AND version = $8
	RETURNING version`
	// Create an args slice containing the values for the placeholder parameters.
	args := []any{
		user.FirstName,
		user.LastName,
		user.Email,
		user.HashedPassword.hash,
		pq.FormatTimestamp(time.Time(user.DOB)),
		user.Activated,
		user.ID,
		user.Version,
	}

	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// Use QueryRowContext() and pass the context as the first argument.// Use QueryRowContext() and pass the context as the first argument.
	err := u.DB.QueryRowContext(ctx, query, args...).Scan(&user.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

func (u UserModel) Delete(id uuid.UUID) error {
	// Return an ErrRecordNotFound error if the movie ID is less than 1.
	// if id < 1 {
	// 	return ErrRecordNotFound
	// }
	// Construct the SQL query to delete the record.
	query := `DELETE FROM users WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// Use ExecContext() and pass the context as the first argument.
	result, err := u.DB.ExecContext(ctx, query, id)
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

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(len(email) <= 500, "email", "must not be more than 500 bytes long")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "email format is not correct")
}

func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) <= 500, "password", "must not be more than 500 bytes long")
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

	// v.Check(user.Email != "", "email", "must be provided")
	// v.Check(len(user.Email) <= 500, "email", "must not be more than 500 bytes long")
	// v.Check(validator.Matches(user.Email, validator.EmailRX), "email", "email format is not correct")
	ValidateEmail(v, user.Email)

	// v.Check(*user.HashedPassword.plaintext != "", "password", "must be provided")
	// v.Check(len(*user.HashedPassword.plaintext) <= 500, "password", "must not be more than 500 bytes long")
	fmt.Println(user.HashedPassword.plaintext)
	if user.HashedPassword.plaintext != nil {
		ValidatePasswordPlaintext(v, *user.HashedPassword.plaintext)
	}
	// If the password hash is ever nil, this will be due to a logic error in our
	// codebase (probably because we forgot to set a password for the user). It's a
	// useful sanity check to include here, but it's not a problem with the data
	// provided by the client. So rather than adding an error to the validation map we
	// raise a panic instead.
	if user.HashedPassword.hash == nil {
		panic("missing password hash for user")
	}
	v.Check(time.Time(user.DOB) != time.Time{}, "dob", "must be provided")
	year := int32(time.Time(user.DOB).Year())
	month := int32(time.Time(user.DOB).Month())
	day := int32(time.Time(user.DOB).Day())
	v.Check(year >= 1900, "dob", "must be greater than 1900")
	v.Check(year <= int32(time.Now().Year()), "dob", "year must not be in the future")
	if year == int32(time.Now().Year()) {
		v.Check(month <= int32(time.Now().Month()), "dob", "month must not be in the future")
		if month == int32(time.Now().Month()) {
			v.Check(day <= int32(time.Now().Day()), "dob", "day must not be in the future")
		}
	}
}

func (m UserModel) GetForToken(tokenScope, tokenPlaintext string) (*User, error) {
	// Calculate the SHA-256 hash of the plaintext token provided by the client.
	// Remember that this returns a byte *array* with length 32, not a slice.
	tokenHash := sha256.Sum256([]byte(tokenPlaintext))
	// Set up the SQL query.
	//createdAt, firstName, lastName, email, hashedPassword, dob, version, activated
	query := `
SELECT users.id, users.createdAt, users.firstName, users.lastName, users.email, users.hashedPassword, users.dob, users.version, users.activated
FROM users
INNER JOIN tokens
ON users.id = tokens.user_id
WHERE tokens.hash = $1
AND tokens.scope = $2
AND tokens.expiry > $3`
	// Create a slice containing the query arguments. Notice how we use the [:] operator
	// to get a slice containing the token hash, rather than passing in the array (which
	// is not supported by the pq driver), and that we pass the current time as the
	// value to check against the token expiry.
	args := []any{tokenHash[:], tokenScope, time.Now()}
	var user User
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// Execute the query, scanning the return values into a User struct. If no matching
	// record is found we return an ErrRecordNotFound error.
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.HashedPassword.hash,
		&user.DOB,
		&user.Version,
		&user.Activated,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	// Return the matching user.
	return &user, nil
}
