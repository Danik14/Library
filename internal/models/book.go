package models

import (
	"time"

	"github.com/Danik14/library/internal/validator"
)

type Book struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	Author    string    `json:"author"`
	Year      uint32    `json:"year,omitempty"`
	Pages     uint32    `json:"pages,omitempty"`
	Genres    []string  `json:"genres,omitempty"`
	Version   int32     `json:"version"`
}

func NewBook(title string, author string, year uint32, pages uint32, genres []string) (*Book, error) {
	return &Book{CreatedAt: time.Now(), Title: title, Author: author, Year: year, Pages: pages, Genres: genres}, nil
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
