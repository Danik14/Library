package models

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Danik14/library/internal/validator"
	"github.com/lib/pq"
)

type Book struct {
	ID        int32     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	Author    string    `json:"author"`
	Year      int32     `json:"year,omitempty"`
	Pages     Pages     `json:"pages,omitempty"`
	Genres    []string  `json:"genres,omitempty"`
	Version   int32     `json:"version"`
}

type BookModel struct {
	DB *sql.DB
}

func (b BookModel) Insert(book *Book) error {
	// return &User{CreatedAt: time.Now(), FirstName: firstName, LastName: lastName, Email: email, HashedPassword: password, DOB: dob, Version: version}, nil
	// Define the SQL query for inserting a new record in the movies table and returning
	// the system-generated data.
	query := `INSERT INTO books (title, author, year, pages, genres) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, version;`
	// Create an args slice containing the values for the placeholder parameters from
	// the movie struct. Declaring this slice immediately next to our SQL query helps to
	// make it nice and clear *what values are being used where* in the query.

	args := []any{book.Title, book.Author, book.Year, book.Pages, pq.Array(book.Genres)}

	// Create a context with a 3-second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Use QueryRowContext() and pass the context as the first argument.
	return b.DB.QueryRowContext(ctx, query, args...).Scan(&book.ID, &book.CreatedAt, &book.Version)
}

func (b BookModel) Get(id int64) (*Book, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	// Define the SQL query for retrieving the book data.
	query := `
SELECT id, created_at, title, year, author, pages, genres, version FROM books
WHERE id = $1`
	// Declare a Movie struct to hold the data returned by the query.
	var book Book
	// Execute the query using the QueryRow() method, passing in the provided id value
	// as a placeholder parameter, and scan the response data into the fields of the
	// Book struct. Importantly, notice that we need to convert the scan target for the
	// genres column using the pq.Array() adapter function again.
	err := b.DB.QueryRow(query, id).Scan(&book.ID,
		&book.CreatedAt, &book.Title, &book.Year, &book.Author, &book.Pages, pq.Array(&book.Genres), &book.Version,
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
	// Otherwise, return a pointer to the Book struct.
	return &book, nil
}

func (b BookModel) Update(book *Book) error {
	// Declare the SQL query for updating the record and returning the new version
	// number.
	query := `
UPDATE books
SET title = $1, author = $2, year = $3, pages = $4, genres = $5, version = version + 1
WHERE id = $6 AND version = $6
RETURNING version`
	// Create an args slice containing the values for the placeholder parameters.
	args := []any{book.Title, book.Author,
		book.Year, book.Pages, pq.Array(book.Genres), book.ID, book.Version,
	}

	// Execute the SQL query. If no matching row could be found, we know the movie
	// version has changed (or the record has been deleted) and we return our custom
	// ErrEditConflict error.
	err := b.DB.QueryRow(query, args...).Scan(&book.Version)
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

func (b BookModel) Delete(id int64) error {
	// Return an ErrRecordNotFound error if the movie ID is less than 1.
	if id < 1 {
		return ErrRecordNotFound
	}
	// Construct the SQL query to delete the record.
	query := `
DELETE FROM books WHERE id = $1`
	// Execute the SQL query using the Exec() method, passing in the id variable as
	// the value for the placeholder parameter. The Exec() method returns a sql.Result
	// object.
	result, err := b.DB.Exec(query, id)
	if err != nil {
		return err
	}
	// Call the RowsAffected() method on the sql.Result object to get the number of rows
	// affected by the query.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	// If no rows were affected, we know that the books table didn't contain a record
	// with the provided ID at the moment we tried to delete it. In that case we
	// return an ErrRecordNotFound error.
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

func ValidateBook(v *validator.Validator, book *Book) {
	// Use the Check() method to execute our validation checks. This will add the
	// provided key and error message to the errors map if the check does not evaluate
	// to true. For example, in the first line here we "check that the title is not
	// equal to the empty string". In the second, we "check that the length of the title
	// is less than or equal to 500 bytes" and so on.
	v.Check(book.Title != "", "title", "must be provided")
	v.Check(len(book.Title) <= 500, "title", "must not be more than 500 bytes long")

	v.Check(book.Author != "", "author", "must be provided")
	v.Check(len(book.Author) <= 500, "author", "must not be more than 500 bytes long")

	v.Check(book.Year > 0, "year", "must be more than 0")
	v.Check(book.Pages > 0, "pages", "must be more than 0")
}
