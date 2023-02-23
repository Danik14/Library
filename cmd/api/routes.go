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
	router.HandlerFunc(http.MethodGet, "/v1/users", app.listUsersHandler)
	router.HandlerFunc(http.MethodPost, "/v1/users", app.createUserHandler)
	router.HandlerFunc(http.MethodGet, "/v1/users/:id", app.showUserHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/users/:id", app.updateUserHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/users/:id", app.deleteUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)

	router.HandlerFunc(http.MethodPatch, "/v1/books/:id", app.requirePermission("books:write", app.updateBookHandler))
	router.HandlerFunc(http.MethodGet, "/v1/books/:id", app.requirePermission("books:read", app.showBookHandler))
	router.HandlerFunc(http.MethodGet, "/v1/books", app.requirePermission("books:read", app.listBooksHandler))
	router.HandlerFunc(http.MethodPost, "/v1/books", app.requirePermission("books:write", app.createBookHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/books/:id", app.requirePermission("books:write", app.deleteBookHandler))

	router.HandlerFunc(http.MethodGet, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	// Wrap the router with the panic recovery middleware.
	return app.recoverPanic(app.rateLimit(app.authenticate(router)))
}
