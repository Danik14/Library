package main

import (
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

	app.writeJSON(w, http.StatusOK, env, nil)
}

func (app *application) listAllUsers(w http.ResponseWriter, r *http.Request) {
	user := models.NewUser("Danik", "Slave", "danik_slave@gmail.com", "123", time.Now(), 1)
	env := envelope{"user": user}

	app.writeJSON(w, http.StatusOK, env, nil)
}
