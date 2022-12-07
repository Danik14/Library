package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (app *application) routes() *mux.Router {
	router := mux.NewRouter()

	router.MethodNotAllowedHandler = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandleFunc("/book/getAll", app.listAllBooks).Methods(http.MethodGet)
	router.HandleFunc("/book/getOne", app.listOneBook).Methods(http.MethodGet)
	router.HandleFunc("/book/create", app.createBook).Methods(http.MethodPost)

	return router
}
