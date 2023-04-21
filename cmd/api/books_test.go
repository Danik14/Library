package main

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/Danik14/library/internal/assert"
)

func TestShowBook(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routesTest())
	defer ts.Close()

	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody string
	}{
		{
			name:     "Valid ID",
			urlPath:  "/v1/books/1",
			wantCode: http.StatusOK,
		},
		{
			name:     "Non-existent ID",
			urlPath:  "/v1/books/2",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Negative ID",
			urlPath:  "/v1/books/-1",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Decimal ID",
			urlPath:  "/v1/books/1.23",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "String ID",
			urlPath:  "/v1/books/foo",
			wantCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			code, _, body := ts.get(t, tt.urlPath)

			assert.Equal(t, code, tt.wantCode)

			if tt.wantBody != "" {
				assert.StringContains(t, body, tt.wantBody)
			}

		})
	}

}

func TestCreateBook(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routesTest())
	defer ts.Close()

	const (
		validTitle = "Test Title"
		validYear  = 2021
		validPages = "105"
	)

	validGenres := []string{"comedy", "drama"}

	tests := []struct {
		name     string
		Title    string
		Year     int32
		Pages    string
		Genres   []string
		wantCode int
	}{
		{
			name:     "Valid submission",
			Title:    validTitle,
			Year:     validYear,
			Pages:    validPages,
			Genres:   validGenres,
			wantCode: http.StatusCreated,
		},
		{
			name:     "Empty Title",
			Title:    "",
			Year:     validYear,
			Pages:    validPages,
			Genres:   validGenres,
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name:     "year < 1888",
			Title:    validTitle,
			Year:     1500,
			Pages:    validPages,
			Genres:   validGenres,
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name:     "test for wrong input",
			Title:    validTitle,
			Year:     validYear,
			Pages:    validPages,
			Genres:   validGenres,
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputData := struct {
				Title  string   `json:"title"`
				Year   int32    `json:"year"`
				Pages  string   `json:"pages"`
				Genres []string `json:"genres"`
			}{
				Title:  tt.Title,
				Year:   tt.Year,
				Pages:  tt.Pages,
				Genres: tt.Genres,
			}

			b, err := json.Marshal(&inputData)
			if err != nil {
				t.Fatal("wrong input data")
			}
			if tt.name == "test for wrong input" {
				b = append(b, 'a')
			}

			code, _, _ := ts.postForm(t, "/v1/books", b)

			assert.Equal(t, code, tt.wantCode)

		})
	}
}

func TestDeleteBook(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routesTest())
	defer ts.Close()

	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody string
	}{
		{
			name:     "deleting existing movie",
			urlPath:  "/v1/books/1",
			wantCode: http.StatusOK,
		},
		{
			name:     "Non-existent ID",
			urlPath:  "/v1/books/2",
			wantCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			code, _, body := ts.deleteReq(t, tt.urlPath)

			assert.Equal(t, code, tt.wantCode)

			if tt.wantBody != "" {
				assert.StringContains(t, body, tt.wantBody)
			}

		})
	}

}
