package models

import (
	"time"

	"github.com/Danik14/library/internal/validator"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	ID             primitive.ObjectID `json:"_id" bson:"_id"`
	CreatedAt      time.Time          `json:"-"`
	FirstName      string             `json:"firstName"`
	LastName       string             `json:"lastName"`
	Email          string             `json:"email"`
	HashedPassword string             `json:"-"`
	DOB            time.Time          `json:"dob"` // date of birth
}

type UserModel struct {
	DB *mongo.Collection
}

// func (u UserModel) Insert(user *User) error {
// 	// return &User{CreatedAt: time.Now(), FirstName: firstName, LastName: lastName, Email: email, HashedPassword: password, DOB: dob, Version: version}, nil
// 	// Define the SQL query for inserting a new record in the movies table and returning
// 	// the system-generated data.
// 	query := `INSERT INTO users (firstname, lastname, email, hashedpassword, dob) VALUES ($1, $2, $3, $4, $5) RETURNING id, createdAt, version;`
// 	// Create an args slice containing the values for the placeholder parameters from
// 	// the movie struct. Declaring this slice immediately next to our SQL query helps to
// 	// make it nice and clear *what values are being used where* in the query.

// 	// hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.HashedPassword), 12)
// 	// if err != nil {
// 	// 	return nil, err
// 	// }
// 	args := []any{user.FirstName, user.LastName, user.Email, user.HashedPassword, pq.FormatTimestamp(user.DOB)}

// 	// Create a context with a 3-second timeout
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	// Use QueryRowContext() and pass the context as the first argument.
// 	return u.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt, &user.Version)
// }

// func (u UserModel) GetAll(firstName string, lastName string, email string, filters data.Filters) ([]*User, data.Metadata, error) {
// 	// Construct the SQL query to retrieve all movie records.
// 	query := fmt.Sprintf(`SELECT count(*) OVER(), id, createdAt, firstName, lastName, email, hashedPassword, dob, version FROM users
// 	WHERE (to_tsvector('simple', firstName) @@ plainto_tsquery('simple', $1) OR $1 = '')
// 	AND (to_tsvector('simple', lastName) @@ plainto_tsquery('simple', $2) OR $2 = '')
// 	AND (to_tsvector('simple', email) @@ plainto_tsquery('simple', $3) OR $3 = '')
// 	ORDER BY %s %s, createdAt DESC
// 	LIMIT $4 OFFSET $5`, filters.SortColumn(), filters.SortDirection())
// 	// Create a context with a 3-second timeout.
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	// As our SQL query now has quite a few placeholder parameters, let's collect the
// 	// values for the placeholders in a slice. Notice here how we call the limit() and
// 	// offset() methods on the Filters struct to get the appropriate values for the
// 	// LIMIT and OFFSET clauses.
// 	args := []any{firstName, lastName, email, filters.Limit(), filters.Offset()}

// 	// Use QueryContext() to execute the query. This returns a sql.Rows resultset
// 	// containing the result.
// 	rows, err := u.DB.QueryContext(ctx, query, args...)
// 	if err != nil {
// 		return nil, data.Metadata{}, err
// 	}
// 	// Importantly, defer a call to rows.Close() to ensure that the resultset is closed
// 	// before GetAll() returns.
// 	defer rows.Close()
// 	// Initialize an empty slice to hold the movie data.
// 	totalRecords := 0
// 	users := []*User{}
// 	// Use rows.Next to iterate through the rows in the resultset.
// 	for rows.Next() {
// 		// Initialize an empty Movie struct to hold the data for an individual movie.
// 		var user User
// 		// Scan the values from the row into the Movie struct. Again, note that we're
// 		// using the pq.Array() adapter on the genres field here.
// 		err := rows.Scan(
// 			&totalRecords,
// 			&user.ID,
// 			&user.CreatedAt,
// 			&user.FirstName,
// 			&user.LastName,
// 			&user.Email,
// 			&user.HashedPassword,
// 			&user.DOB,
// 			&user.Version,
// 		)
// 		if err != nil {
// 			return nil, data.Metadata{}, err
// 		}
// 		// Add the Movie struct to the slice.
// 		users = append(users, &user)
// 	}
// 	// When the rows.Next() loop has finished, call rows.Err() to retrieve any error
// 	// that was encountered during the iteration.
// 	if err = rows.Err(); err != nil {
// 		return nil, data.Metadata{}, err
// 	}

// 	// Generate a Metadata struct, passing in the total record count and pagination
// 	// parameters from the client.
// 	metadata := data.CalculateMetadata(totalRecords, filters.Page, filters.PageSize)

// 	// If everything went OK, then return the slice of movies.
// 	return users, metadata, nil
// }

// func (u UserModel) Get(id uuid.UUID) (*User, error) {
// 	// if id < 1 {
// 	// 	return nil, ErrRecordNotFound
// 	// }
// 	// Define the SQL query for retrieving the movie data.
// 	query := `SELECT id, createdAt, firstName, lastName, email, hashedPassword, dob, version FROM users WHERE id = $1`

// 	// Declare a Movie struct to hold the data returned by the query.
// 	var user User // Execute the query using the QueryRow() method, passing in the provided id value

// 	// Use the context.WithTimeout() function to create a context.Context which carries a
// 	// 3-second timeout deadline. Note that we're using the empty context.Background()
// 	// as the 'parent' context.
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	// Importantly, use defer to make sure that we cancel the context before the Get()
// 	// method returns.
// 	defer cancel()
// 	// Use the QueryRowContext() method to execute the query, passing in the context
// 	// with the deadline as the first argument.
// 	err := u.DB.QueryRowContext(ctx, query, id).Scan(
// 		&user.ID,
// 		&user.CreatedAt,
// 		&user.FirstName,
// 		&user.LastName,
// 		&user.Email,
// 		&user.HashedPassword,
// 		&user.DOB,
// 		&user.Version,
// 	)
// 	// Handle any errors. If there was no matching movie found, Scan() will return
// 	// a sql.ErrNoRows error. We check for this and return our custom ErrRecordNotFound
// 	// error instead.
// 	if err != nil {
// 		switch {
// 		case errors.Is(err, sql.ErrNoRows):
// 			return nil, ErrRecordNotFound
// 		default:
// 			return nil, err
// 		}
// 	}
// 	// Otherwise, return a pointer to the Movie struct.
// 	return &user, nil
// }

// func (u UserModel) Update(user *User) error {
// 	// Declare the SQL query for updating the record and returning the new version
// 	// number.
// 	query := `
// 	UPDATE users
// 	SET firstName = $1, lastName = $2, email = $3, hashedPassword = $4, dob = $5, version = version + 1
// 	WHERE id = $6 AND version = $7
// 	RETURNING version`
// 	// Create an args slice containing the values for the placeholder parameters.
// 	args := []any{
// 		user.FirstName,
// 		user.LastName,
// 		user.Email,
// 		user.HashedPassword,
// 		user.DOB,
// 		user.ID,
// 		user.Version,
// 	}

// 	// Create a context with a 3-second timeout.
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()
// 	// Use QueryRowContext() and pass the context as the first argument.// Use QueryRowContext() and pass the context as the first argument.
// 	err := u.DB.QueryRowContext(ctx, query, args...).Scan(&user.Version)
// 	if err != nil {
// 		switch {
// 		case errors.Is(err, sql.ErrNoRows):
// 			return ErrEditConflict
// 		default:
// 			return err
// 		}
// 	}

// 	return nil
// }

// func (u UserModel) Delete(id uuid.UUID) error {
// 	// Return an ErrRecordNotFound error if the movie ID is less than 1.
// 	// if id < 1 {
// 	// 	return ErrRecordNotFound
// 	// }
// 	// Construct the SQL query to delete the record.
// 	query := `DELETE FROM users WHERE id = $1`

// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()
// 	// Use ExecContext() and pass the context as the first argument.
// 	result, err := u.DB.ExecContext(ctx, query, id)
// 	if err != nil {
// 		return err
// 	}

// 	// Call the RowsAffected() method on the sql.Result object to get the number of rows
// 	// affected by the query.
// 	rowsAffected, err := result.RowsAffected()
// 	if err != nil {
// 		return err
// 	}
// 	// If no rows were affected, we know that the movies table didn't contain a record
// 	// with the provided ID at the moment we tried to delete it. In that case we
// 	// return an ErrRecordNotFound error.
// 	if rowsAffected == 0 {
// 		return ErrRecordNotFound
// 	}
// 	return nil
// }

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
