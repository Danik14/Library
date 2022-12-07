package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Danik14/library/internal/models"
)

func (app *application) listAllBooks(w http.ResponseWriter, r *http.Request) {
	book, err := models.NewBook(1, "AlibaSlave", "Danik", "1964", 300, []string{"gachi"}, 1)
	if err != nil {
		app.logger.Fatal("Book error")
	}
	env := envelope{"book": book}

	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, "JSON error")
		return
	}
}

func (app *application) listOneBook(w http.ResponseWriter, r *http.Request) {
	book, err := models.NewBook(1, "AlibaSlave", "Danik", "1964", 300, []string{"gachi"}, 1)
	if err != nil {
		app.logger.Fatal("User error")
	}
	env := envelope{"book": book}

	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, "JSON error")
		return
	}
}

func (app *application) createBook(w http.ResponseWriter, r *http.Request) {
	book := &models.Book{}

	err := app.readJSON(w, r, book)
	if err != nil {
		app.logError(r, err)
		app.errorResponse(w, r, http.StatusBadRequest, "JSON error")
		return
	}

	fmt.Println(book)
}

func (app *application) listAllUsers(w http.ResponseWriter, r *http.Request) {
	user, err := models.NewUser("Danik", "Slave", "danik_slave@gmail.com", "123", time.Now(), 1)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, "JSON Error")
		return
	}
	env := envelope{"user": user}

	app.writeJSON(w, http.StatusOK, env, nil)
}

func (app *application) createUser(w http.ResponseWriter, r *http.Request) {
	book := &models.Book{}

	err := app.readJSON(w, r, book)
	if err != nil {
		app.logError(r, err)
		app.errorResponse(w, r, http.StatusBadRequest, "JSON error")
		return
	}

	fmt.Println(book)
}
