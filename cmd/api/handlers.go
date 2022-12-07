package main

import (
	"net/http"

	"github.com/Danik14/library/internal/models"
)

func (app *application) listAllBooks(w http.ResponseWriter, r *http.Request) {
	book := models.NewBook(1, "AlibaSlave", "Danik", "1964", 300, []string{"gachi"}, 1)
	env := envelope{"book": book}

	app.writeJSON(w, http.StatusOK, env, nil)
}
