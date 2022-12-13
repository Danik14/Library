package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Danik14/library/internal/models"
	"github.com/Danik14/library/internal/validator"
)

func (app *application) listAllBooks(w http.ResponseWriter, r *http.Request) {
	book, err := models.NewBook("AlibaSlave", "Danik", 1964, 300, []string{"gachi"})
	if err != nil {
		app.logger.Fatal("Book error")
	}
	env := envelope{"book": book}

	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}
}

func (app *application) listOneBook(w http.ResponseWriter, r *http.Request) {
	book, err := models.NewBook("AlibaSlave", "Danik", 1964, 300, []string{"gachi"})
	if err != nil {
		app.logger.Fatal("User error")
	}
	env := envelope{"book": book}

	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}
}

func (app *application) createBook(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Title  string   `json:"title"`
		Author string   `json:"author"`
		Year   uint32   `json:"year"`
		Pages  uint32   `json:"pages"`
		Genres []string `json:"genres"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.logError(r, err)
		app.badRequestResponse(w, r, err)
		return
	}

	book, err := models.NewBook(input.Title, input.Author, input.Year, input.Pages, input.Genres)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	fmt.Println(book)
}

func (app *application) listAllUsers(w http.ResponseWriter, r *http.Request) {
	user, err := models.NewUser("Danik", "Slave", "danik_slave@gmail.com", "123", time.Now(), 1)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	env := envelope{"user": user}

	app.writeJSON(w, http.StatusOK, env, nil)
}

func (app *application) createUser(w http.ResponseWriter, r *http.Request) {
	var input struct {
		FirstName string    `json:"firstName"`
		LastName  string    `json:"lastName"`
		Email     string    `json:"email"`
		Password  string    `json:"password"`
		DOB       time.Time `json:"dob"` // date of birth
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.logError(r, err)
		app.badRequestResponse(w, r, err)
		return
	}

	// Copy the values from the input struct to a new Movie struct.
	user := &models.User{
		FirstName:      input.FirstName,
		LastName:       input.LastName,
		Email:          input.Email,
		HashedPassword: input.Password,
		DOB:            input.DOB,
	}

	// Initialize a new Validator.
	v := validator.New()
	// Call the ValidateMovie() function and return a response containing the errors if
	// any of the checks fail.
	if models.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	fmt.Fprintf(w, "%+v\n", input)

	// fmt.Println(input)
}
