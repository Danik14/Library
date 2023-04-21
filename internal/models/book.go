package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Danik14/library/internal/data"
	"github.com/Danik14/library/internal/validator"
	"github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
)

type Book struct {
	ID        uuid.UUID `json:"id"`
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
	// Define the SQL query for inserting a new record in the books table and returning
	// the system-generated data.
	query := `INSERT INTO books (title, author, year, pages, genres) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, version;`
	// Create an args slice containing the values for the placeholder parameters from
	// the book struct. Declaring this slice immediately next to our SQL query helps to
	// make it nice and clear *what values are being used where* in the query.

	args := []any{book.Title, book.Author, book.Year, book.Pages, pq.Array(book.Genres)}

	// Create a context with a 3-second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Use QueryRowContext() and pass the context as the first argument.
	return b.DB.QueryRowContext(ctx, query, args...).Scan(&book.ID, &book.CreatedAt, &book.Version)
}

func (b BookModel) Get(id uuid.UUID) (*Book, error) {
	// Define the SQL query for retrieving the book data.
	query := `
SELECT id, created_at, title, year, author, pages, genres, version FROM books
WHERE id = $1`
	// Declare a Book struct to hold the data returned by the query.
	var book Book

	// Use the context.WithTimeout() function to create a context.Context which carries a
	// 3-second timeout deadline. Note that we're using the empty context.Background()
	// as the 'parent' context.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	// Importantly, use defer to make sure that we cancel the context before the Get()
	// method returns.
	// The defer cancel() line is necessary because it ensures that the resources
	// associated with our context will always be released before the Get() method returns,
	// thereby preventing a memory leak. Without it, the resources won’t be released until
	// either the 3- second timeout is hit or the parent context
	// (which in this specific example is context.Background()) is canceled.
	defer cancel()

	// Execute the query using the QueryRow() method, passing in the provided id value
	// as a placeholder parameter, and scan the response data into the fields of the
	// Book struct. Importantly, notice that we need to convert the scan target for the
	// genres column using the pq.Array() adapter function again.
	err := b.DB.QueryRowContext(ctx, query, id).Scan(&book.ID,
		&book.CreatedAt, &book.Title, &book.Year, &book.Author, &book.Pages, pq.Array(&book.Genres), &book.Version,
	)
	// Handle any errors. If there was no matching book found, Scan() will return
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

// Create a new GetAll() method which returns a slice of books. Although we're not
// using them right now, we've set this up to accept the various filter parameters as
// arguments.
func (m BookModel) GetAll(title string, author string, genres []string, filters data.Filters) ([]*Book, data.Metadata, error) {
	// Construct the SQL query to retrieve all book records.
	query := fmt.Sprintf(`
SELECT count(*) OVER(), id, created_at, title, author, year, pages, genres, version FROM books
WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
AND (to_tsvector('simple', author) @@ plainto_tsquery('simple', $2) OR $2 = '')
AND (genres @> $3 OR $3 = '{}')
ORDER BY %s %s, id ASC
LIMIT $4 OFFSET $5`, filters.SortColumn(), filters.SortDirection())
	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{title, author, pq.Array(genres), filters.Limit(), filters.Offset()}

	// Pass the title and genres as the placeholder parameter values.
	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		fmt.Println(query)
		return nil, data.Metadata{}, err
	}

	// Importantly, defer a call to rows.Close() to ensure that the resultset is closed
	// before GetAll() returns.
	defer rows.Close()

	totalRecords := 0
	books := []*Book{}
	// Use rows.Next to iterate through the rows in the resultset.
	for rows.Next() {
		// Initialize an empty Book struct to hold the data for an individual book.
		var book Book
		// Scan the values from the row into the Book struct. Again, note that we're
		// using the pq.Array() adapter on the genres field here.
		err := rows.Scan(
			&totalRecords,
			&book.ID,
			&book.CreatedAt,
			&book.Title,
			&book.Author,
			&book.Year,
			&book.Pages,
			pq.Array(&book.Genres),
			&book.Version,
		)
		if err != nil {
			return nil, data.Metadata{}, err
		}
		// Add the Book struct to the slice.
		books = append(books, &book)
	}
	// When the rows.Next() loop has finished, call rows.Err() to retrieve any error
	// that was encountered during the iteration.
	if err = rows.Err(); err != nil {
		return nil, data.Metadata{}, err
	}

	// Generate a Metadata struct, passing in the total record count and pagination
	// parameters from the client.
	metadata := data.CalculateMetadata(totalRecords, filters.Page, filters.PageSize)

	// If everything went OK, then return the slice of books.
	return books, metadata, nil
}

func (b BookModel) Update(book *Book) error {
	// Declare the SQL query for updating the record and returning the new version
	// number.
	query := `
UPDATE books
SET title = $1, author = $2, year = $3, pages = $4, genres = $5, version = version + 1
WHERE id = $6 AND version = $7
RETURNING version`
	// Create an args slice containing the values for the placeholder parameters.
	args := []any{book.Title, book.Author,
		book.Year, book.Pages, pq.Array(book.Genres), book.ID, book.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute the SQL query. If no matching row could be found, we know the book
	// version has changed (or the record has been deleted) and we return our custom
	// ErrEditConflict error.
	err := b.DB.QueryRowContext(ctx, query, args...).Scan(&book.Version)
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

func (b BookModel) Delete(id uuid.UUID) error {
	// Return an ErrRecordNotFound error if the book ID is less than 1.
	// Construct the SQL query to delete the record.
	query := `
DELETE FROM books WHERE id = $1`
	// Execute the SQL query using the Exec() method, passing in the id variable as
	// the value for the placeholder parameter. The Exec() method returns a sql.Result
	// object.

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := b.DB.ExecContext(ctx, query, id)
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

type MockBookModel struct{}

func (b MockBookModel) Insert(book *Book) error {
	return nil
}

func (b MockBookModel) Get(id uuid.UUID) (*Book, error) {
	switch id {
	case uuid.UUID{}:
		return &Book{
			ID:        uuid.UUID{},
			CreatedAt: time.Now(),
			Year:      2023,
			Pages:     105,
			Title:     "Test Mock",
			Genres:    []string{""},
		}, nil
	default:
		return nil, ErrRecordNotFound
	}
}
func (b MockBookModel) Update(book *Book) error {
	return nil
}

func (b MockBookModel) Delete(id uuid.UUID) error {
	switch id {
	case uuid.UUID{}:
		return nil
	default:
		return ErrRecordNotFound
	}
}

func (b MockBookModel) GetAll(title string, author string, genres []string, filters data.Filters) ([]*Book, data.Metadata, error) {
	return nil, data.Metadata{}, nil
}
