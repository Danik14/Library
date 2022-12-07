package models

import "time"

type Book struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	Author    string    `json:"author"`
	Year      string    `json:"year,omitempty"`
	Pages     int       `json:"pages,omitempty"`
	Genres    []string  `json:"genres,omitempty"`
	Version   int32     `json:"version"`
}

func NewBook(id int64, title string, author string, year string, pages int, genres []string, version int32) (*Book, error) {
	return &Book{ID: id, CreatedAt: time.Now(), Title: title, Author: author, Year: year, Pages: pages, Genres: genres, Version: version}, nil
}
