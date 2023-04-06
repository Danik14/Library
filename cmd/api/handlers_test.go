package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/Danik14/library/internal/models"
	"github.com/joho/godotenv"
)

func newTestApp(t *testing.T) *application {

	err := godotenv.Load()
	if err != nil {
		t.Errorf("newTestApp() FAILED. .env file error: %s", err.Error())
	}

	dsn := os.Getenv("DB_DSN")
	cfg := config{db: struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}{dsn: dsn, maxOpenConns: 25, maxIdleConns: 25, maxIdleTime: "3m"}}

	db, err := openDB(cfg)
	if err != nil {
		t.Errorf("newTestApp() FAILED. db connction error: %s", err.Error())
	}
	return &application{models: models.NewModels(db)}
}

func TestCreateBookHandler(t *testing.T) {
	// Create a new instance of the application struct with any necessary dependencies.
	app := newTestApp(t)

	// Create a new HTTP POST request with a JSON request body containing book data.
	jsonBody := `{"title": "Test Book", "author": "Test Author", "year": 2022, "pages": 123, "genres": ["fiction", "mystery"]}`
	req := httptest.NewRequest("POST", "/v1/books", strings.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Create a new ResponseRecorder to record the response from the handler.
	rr := httptest.NewRecorder()

	// Call the createBookHandler method with the ResponseRecorder and HTTP request.
	app.createBookHandler(rr, req)

	// Check the HTTP status code in the response recorder. It should be 201.
	if rr.Code != http.StatusCreated {
		t.Errorf("expected HTTP status 201; got %d", rr.Code)
	}

	// Check the Location header in the response. It should match the expected value.
	location, _ := rr.HeaderMap["Location"]
	expectedLocation := "/v1/books/1"
	if len(location) == 0 || location[0] != expectedLocation {
		t.Errorf("expected Location header %q; got %q", expectedLocation, location)
	}

	// Check the JSON response body. It should contain the book data in the expected format.
	expectedBody := `{"book":{"id":1,"title":"Test Book","author":"Test Author","year":2022,"pages":123,"genres":["fiction","mystery"]}}`
	if rr.Body.String() != expectedBody {
		t.Errorf("expected response body %q; got %q", expectedBody, rr.Body.String())
	}
}
