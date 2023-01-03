package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	// router.HandlerFunc(http.MethodGet, "/", app.listAllBooks)

	// router.HandleFunc("/user", app.listAllUsers).Methods(http.MethodGet)
	// router.HandlerFunc(http.MethodGet, "/user", app.listUsersHandler)
	// router.HandlerFunc(http.MethodPost, "/user", app.createUserHandler)
	// router.HandlerFunc(http.MethodGet, "/user/:id", app.showUserHandler)
	// router.HandlerFunc(http.MethodPatch, "/user/:id", app.updateUserHandler)
	// router.HandlerFunc(http.MethodDelete, "/user/:id", app.deleteUserHandler)

	// router.HandlerFunc(http.MethodPatch, "/v1/book/:id", app.updateBookHandler)
	// router.HandlerFunc(http.MethodGet, "/v1/book/:id", app.showBookHandler)
	// router.HandlerFunc(http.MethodGet, "/v1/book", app.listBooksHandler)
	router.HandlerFunc(http.MethodPost, "/v1/book", app.createBookHandler)
	// router.HandlerFunc(http.MethodDelete, "/v1/book/:id", app.deleteMovieHandler)

	// Wrap the router with the panic recovery middleware.
	return app.recoverPanic(app.rateLimit(router))
}
