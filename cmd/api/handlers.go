package main

import (
	"fmt"
	"net/http"

	"github.com/Danik14/library/internal/models"
)

func (app *application) listAllBooks(w http.ResponseWriter, r *http.Request) {
	book := models.NewBook(1, "AlibaSlave", "Danik", "1964", 300, []string{"gachi"}, 1)
	env := envelope{"book": book}

	err := app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, "JSON error")
		return
	}
}

func (app *application) listOneBook(w http.ResponseWriter, r *http.Request) {
	book := models.NewBook(1, "AlibaSlave", "Danik", "1964", 300, []string{"gachi"}, 1)

	env := envelope{"book": book}

	err := app.writeJSON(w, http.StatusOK, env, nil)
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
