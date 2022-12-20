package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()

	// router.MethodNotAllowedHandler = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/", app.listAllBooks)

	// router.HandleFunc("/user", app.listAllUsers).Methods(http.MethodGet)
	router.HandlerFunc(http.MethodGet, "/user", app.listUsersHandler)
	router.HandlerFunc(http.MethodPost, "/user", app.createUserHandler)
	router.HandlerFunc(http.MethodGet, "/user/:id", app.showUserHandler)
	router.HandlerFunc(http.MethodPatch, "/user/:id", app.updateUserHandler)
	router.HandlerFunc(http.MethodDelete, "/user/:id", app.deleteUserHandler)

	router.HandlerFunc(http.MethodGet, "/books", app.listAllBooks)
	router.HandlerFunc(http.MethodGet, "/books/:id", app.listOneBook)
	router.HandlerFunc(http.MethodPost, "/v1/books", app.createBook)

	return router
}
