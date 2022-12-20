package models

import (
	"context"
	"database/sql"
	"time"

	"github.com/Danik14/library/internal/validator"
	"github.com/lib/pq"
)

type Book struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	Author    string    `json:"author"`
	Year      uint32    `json:"year,omitempty"`
	Pages     Pages     `json:"pages,omitempty"`
	Genres    []string  `json:"genres,omitempty"`
	Version   int32     `json:"version"`
}

func NewBook(title string, author string, year uint32, pages Pages, genres []string) (*Book, error) {
	return &Book{CreatedAt: time.Now(), Title: title, Author: author, Year: year, Pages: pages, Genres: genres}, nil
}

type BookModel struct {
	DB *sql.DB
}

func (u BookModel) Insert(book *Book) error {
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
	return u.DB.QueryRowContext(ctx, query, args...).Scan(&book.ID, &book.CreatedAt, &book.Version)
}

func (m BookModel) Get(id int64) (*Book, error) {
	return nil, nil
}

func (m BookModel) Update(movie *Book) error {
	return nil
}

func (m BookModel) Delete(id int64) error {
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
