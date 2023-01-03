package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Danik14/library/internal/models"
	"github.com/Danik14/library/internal/validator"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// func (app *application) listAllBooks(w http.ResponseWriter, r *http.Request) {
// 	book, err := models.NewBook("AlibaSlave", "Danik", 1964, 300, []string{"gachi"})
// 	if err != nil {
// 		app.logger.PrintFatal("Book error")
// 	}
// 	env := envelope{"book": book}

// 	err = app.writeJSON(w, http.StatusOK, env, nil)
// 	if err != nil {
// 		app.errorResponse(w, r, http.StatusBadRequest, err.Error())
// 		return
// 	}
// }

// func (app *application) listOneBook(w http.ResponseWriter, r *http.Request) {
// 	book, err := models.NewBook("AlibaSlave", "Danik", 1964, 300, []string{"gachi"})
// 	if err != nil {
// 		app.logger.Fatal("User error")
// 	}
// 	env := envelope{"book": book}

// 	err = app.writeJSON(w, http.StatusOK, env, nil)
// 	if err != nil {
// 		app.errorResponse(w, r, http.StatusBadRequest, err.Error())
// 		return
// 	}
// }

func (app *application) createBookHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Title  string       `json:"title"`
		Author string       `json:"author"`
		Year   int32        `json:"year"`
		Pages  models.Pages `json:"pages"`
		Genres []string     `json:"genres"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.logError(r, err)
		app.badRequestResponse(w, r, err)
		return
	}

	book := &models.Book{
		Title:  input.Title,
		Author: input.Author,
		Year:   input.Year,
		Genres: input.Genres,
		Pages:  input.Pages,
	}

	v := validator.New()

	if models.ValidateBook(v, book); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	book.ID = primitive.NewObjectID()
	_, anyerr := app.models.Books.DB.InsertOne(ctx, book)
	if anyerr != nil {
		app.serverErrorResponse(w, r, anyerr)
	}

	// // Call the Insert() method on our books model, passing in a pointer to the
	// // validated book struct. This will create a record in the database and update the
	// // book struct with the system-generated information.
	// err = app.models.Books.Insert(book)
	// if err != nil {
	// 	app.serverErrorResponse(w, r, err)
	// 	return
	// }

	// // When sending a HTTP response, we want to include a Location header to let the
	// // client know which URL they can find the newly-created resource at. We make an
	// // empty http.Header map and then use the Set() method to add a new Location header,
	// // interpolating the system-generated ID for our new movie in the URL.
	// headers := make(http.Header)
	// headers.Set("Location", fmt.Sprintf("/v1/books/%d", book.ID))

	// // Write a JSON response with a 201 Created status code, the movie data in the
	// // response body, and the Location header.
	// err = app.writeJSON(w, http.StatusCreated, envelope{"book": book}, headers)
	// if err != nil {
	// 	app.serverErrorResponse(w, r, err)
	// }

	fmt.Println(book)
}

// func (app *application) listAllUsers(w http.ResponseWriter, r *http.Request) {
// 	user, err := models.NewUser("Danik", "Slave", "danik_slave@gmail.com", "123", time.Now(), 1)
// 	if err != nil {
// 		app.badRequestResponse(w, r, err)
// 		return
// 	}
// 	env := envelope{"user": user}

// 	app.writeJSON(w, http.StatusOK, env, nil)
// }

// func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {
// 	var input struct {
// 		FirstName string    `json:"firstName"`
// 		LastName  string    `json:"lastName"`
// 		Email     string    `json:"email"`
// 		Password  string    `json:"password"`
// 		DOB       time.Time `json:"dob"` // date of birth
// 	}
// 	err := app.readJSON(w, r, &input)
// 	if err != nil {
// 		app.logError(r, err)
// 		app.badRequestResponse(w, r, err)
// 		return
// 	}

// 	// Copy the values from the input struct to a new Movie struct.
// 	user := &models.User{
// 		FirstName:      input.FirstName,
// 		LastName:       input.LastName,
// 		Email:          input.Email,
// 		HashedPassword: input.Password,
// 		DOB:            input.DOB,
// 	}

// 	// Initialize a new Validator.
// 	v := validator.New()
// 	// Call the ValidateMovie() function and return a response containing the errors if
// 	// any of the checks fail.
// 	if models.ValidateUser(v, user); !v.Valid() {
// 		app.failedValidationResponse(w, r, v.Errors)
// 		return
// 	}

// 	err = app.models.Users.Insert(user)
// 	if err != nil {
// 		app.serverErrorResponse(w, r, err)
// 		return
// 	}
// 	// When sending a HTTP response, we want to include a Location header to let the
// 	// client know which URL they can find the newly-created resource at. We make an
// 	// empty http.Header map and then use the Set() method to add a new Location header,
// 	// interpolating the system-generated ID for our new movie in the URL.
// 	headers := make(http.Header)
// 	headers.Set("Location", fmt.Sprintf("/user/%d", user.ID))
// 	// Write a JSON response with a 201 Created status code, the movie data in the
// 	// response body, and the Location header.
// 	err = app.writeJSON(w, http.StatusCreated, envelope{"user": user}, headers)
// 	if err != nil {
// 		app.serverErrorResponse(w, r, err)
// 	}
// }

// func (app *application) showBookHandler(w http.ResponseWriter, r *http.Request) {
// 	id, err := app.readUUIDParam(r)
// 	if err != nil {
// 		app.notFoundResponse(w, r)
// 		return
// 	}
// 	// Call the Get() method to fetch the data for a specific book. We also need to
// 	// use the errors.Is() function to check if it returns a data.ErrRecordNotFound
// 	// error, in which case we send a 404 Not Found response to the client.
// 	book, err := app.models.Books.Get(id)
// 	if err != nil {
// 		switch {
// 		case errors.Is(err, models.ErrRecordNotFound):
// 			app.notFoundResponse(w, r)
// 		default:
// 			app.serverErrorResponse(w, r, err)
// 		}
// 		return
// 	}
// 	err = app.writeJSON(w, http.StatusOK, envelope{"book": book}, nil)
// 	if err != nil {
// 		app.serverErrorResponse(w, r, err)
// 	}
// }

// func (app *application) listBooksHandler(w http.ResponseWriter, r *http.Request) {
// 	// To keep things consistent with our other handlers, we'll define an input struct
// 	// to hold the expected values from the request query string.
// 	var input struct {
// 		Title  string
// 		Author string
// 		// Year     int
// 		Genres []string
// 		data.Filters
// 	}
// 	// Initialize a new Validator instance.
// 	v := validator.New()
// 	// Call r.URL.Query() to get the url.Values map containing the query string data.
// 	qs := r.URL.Query()
// 	// Use our helpers to extract the title and genres query string values, falling back
// 	// to defaults of an empty string and an empty slice respectively if they are not
// 	// provided by the client.
// 	input.Title = app.readString(qs, "title", "")
// 	input.Author = app.readString(qs, "author", "")
// 	// input.Year = app.readInt(qs, "year", 0, v)
// 	input.Genres = app.readCSV(qs, "genres", []string{})
// 	// Get the page and page_size query string values as integers. Notice that we set
// 	// the default page value to 1 and default page_size to 20, and that we pass the
// 	// validator instance as the final argument here.
// 	input.Filters.Page = app.readInt(qs, "page", 1, v)
// 	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
// 	// Extract the sort query string value, falling back to "id" if it is not provided
// 	// by the client (which will imply a ascending sort on movie ID).
// 	input.Filters.Sort = app.readString(qs, "sort", "id")
// 	input.Filters.SortSafelist = []string{"id", "title", "author", "year", "runtime", "-id", "-title", "-author", "-year", "-runtime"}

// 	// Check the Validator instance for any errors and use the failedValidationResponse()
// 	// helper to send the client a response if necessary.
// 	if data.ValidateFilters(v, input.Filters); !v.Valid() {
// 		app.failedValidationResponse(w, r, v.Errors)
// 		return
// 	}
// 	// Call the GetAll() method to retrieve the books, passing in the various filter
// 	// parameters.
// 	books, err := app.models.Books.GetAll(input.Title, input.Author, input.Genres, input.Filters)
// 	if err != nil {
// 		app.serverErrorResponse(w, r, err)
// 		return
// 	}
// 	// Send a JSON response containing the movie data.
// 	err = app.writeJSON(w, http.StatusOK, envelope{"books": books}, nil)
// 	if err != nil {
// 		app.serverErrorResponse(w, r, err)
// 	}
// }

// func (app *application) updateBookHandler(w http.ResponseWriter, r *http.Request) {
// 	// Extract the movie ID from the URL.
// 	id, err := app.readUUIDParam(r)
// 	if err != nil {
// 		app.notFoundResponse(w, r)
// 		return
// 	}
// 	// Fetch the existing movie record from the database, sending a 404 Not Found
// 	// response to the client if we couldn't find a matching record.
// 	book, err := app.models.Books.Get(id)
// 	if err != nil {
// 		switch {
// 		case errors.Is(err, models.ErrRecordNotFound):
// 			app.notFoundResponse(w, r)
// 		default:
// 			app.serverErrorResponse(w, r, err)
// 		}
// 		return
// 	}

// 	// Declare an input struct to hold the expected data from the client.
// 	var input struct {
// 		Title  *string       `json:"title"`
// 		Author *string       `json:"author"`
// 		Year   *int32        `json:"year"`
// 		Pages  *models.Pages `json:"pages"`
// 		Genres []string      `json:"genres"`
// 	}

// 	err = app.readJSON(w, r, &input)
// 	if err != nil {
// 		app.badRequestResponse(w, r, err)
// 		return
// 	}

// 	// If the input.Title value is nil then we know that no corresponding "title" key/
// 	// value pair was provided in the JSON request body. So we move on and leave the
// 	// movie record unchanged. Otherwise, we update the movie record with the new title
// 	// value. Importantly, because input.Title is a now a pointer to a string, we need
// 	// to dereference the pointer using the * operator to get the underlying value
// 	// before assigning it to our movie record.
// 	if input.Title != nil {
// 		book.Title = *input.Title
// 	}
// 	if input.Author != nil {
// 		book.Author = *input.Author
// 	}
// 	if input.Year != nil {
// 		book.Year = *input.Year
// 	}
// 	if input.Pages != nil {
// 		book.Pages = *input.Pages
// 	}
// 	if input.Genres != nil {
// 		book.Genres = input.Genres
// 	}

// 	// Validate the updated movie record, sending the client a 422 Unprocessable Entity
// 	// response if any checks fail.
// 	v := validator.New()
// 	if models.ValidateBook(v, book); !v.Valid() {
// 		app.failedValidationResponse(w, r, v.Errors)
// 		return
// 	}
// 	// Intercept any ErrEditConflict error and call the new editConflictResponse()
// 	// helper.
// 	err = app.models.Books.Update(book)
// 	if err != nil {
// 		switch {
// 		case errors.Is(err, models.ErrEditConflict):
// 			app.editConflictResponse(w, r)
// 		default:
// 			app.serverErrorResponse(w, r, err)
// 		}
// 		return
// 	}
// 	// Write the updated movie record in a JSON response.
// 	err = app.writeJSON(w, http.StatusOK, envelope{"book": book}, nil)
// 	if err != nil {
// 		app.serverErrorResponse(w, r, err)
// 	}
// }

// func (app *application) deleteMovieHandler(w http.ResponseWriter, r *http.Request) {
// 	// Extract the movie ID from the URL.
// 	id, err := app.readUUIDParam(r)
// 	if err != nil {
// 		app.notFoundResponse(w, r)
// 		return
// 	}
// 	// Delete the movie from the database, sending a 404 Not Found response to the
// 	// client if there isn't a matching record.
// 	err = app.models.Books.Delete(id)
// 	if err != nil {
// 		switch {
// 		case errors.Is(err, models.ErrRecordNotFound):
// 			app.notFoundResponse(w, r)
// 		default:
// 			app.serverErrorResponse(w, r, err)
// 		}
// 		return
// 	}
// 	// Return a 200 OK status code along with a success message.
// 	err = app.writeJSON(w, http.StatusOK, envelope{"message": "book successfully deleted"}, nil)
// 	if err != nil {
// 		app.serverErrorResponse(w, r, err)
// 	}
// }

// func (app *application) showUserHandler(w http.ResponseWriter, r *http.Request) {
// 	id, err := app.readUUIDParam(r)
// 	if err != nil {
// 		app.notFoundResponse(w, r)
// 		return
// 	}

// 	user, err := app.models.Users.Get(id)
// 	if err != nil {
// 		switch {
// 		case errors.Is(err, models.ErrRecordNotFound):
// 			app.notFoundResponse(w, r)
// 		default:
// 			app.serverErrorResponse(w, r, err)
// 		}
// 		return
// 	}

// 	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
// 	if err != nil {
// 		app.serverErrorResponse(w, r, err)
// 	}
// }

// func (app *application) updateUserHandler(w http.ResponseWriter, r *http.Request) {
// 	id, err := app.readUUIDParam(r)
// 	if err != nil {
// 		app.notFoundResponse(w, r)
// 		return
// 	}
// 	// Fetch the existing movie record from the database, sending a 404 Not Found
// 	// response to the client if we couldn't find a matching record.
// 	user, err := app.models.Users.Get(id)
// 	if err != nil {
// 		switch {
// 		case errors.Is(err, models.ErrRecordNotFound):
// 			app.notFoundResponse(w, r)
// 		default:
// 			app.serverErrorResponse(w, r, err)
// 		}
// 		return
// 	}
// 	// Declare an input struct to hold the expected data from the client.
// 	var input struct {
// 		FirstName *string    `json:"firstName"`
// 		LastName  *string    `json:"lastName"`
// 		Email     *string    `json:"email"`
// 		Password  *string    `json:"password"`
// 		DOB       *time.Time `json:"dob"` // date of birth
// 	}
// 	// Read the JSON request body data into the input struct.
// 	err = app.readJSON(w, r, &input)
// 	if err != nil {
// 		app.badRequestResponse(w, r, err)
// 		return
// 	}
// 	// Copy the values from the request body to the appropriate fields of the movie
// 	// record.
// 	if input.FirstName != nil {
// 		user.FirstName = *input.FirstName
// 	}
// 	if input.LastName != nil {
// 		user.LastName = *input.LastName
// 	}
// 	if input.Email != nil {
// 		user.Email = *input.Email
// 	}
// 	if input.Password != nil {
// 		user.HashedPassword = *input.Password
// 	}
// 	if input.DOB != nil {
// 		user.DOB = *input.DOB
// 	}

// 	// Validate the updated movie record, sending the client a 422 Unprocessable Entity
// 	// response if any checks fail.
// 	v := validator.New()
// 	if models.ValidateUser(v, user); !v.Valid() {
// 		app.failedValidationResponse(w, r, v.Errors)
// 		return
// 	}
// 	// Pass the updated movie record to our new Update() method.
// 	err = app.models.Users.Update(user)
// 	if err != nil {

// 		switch {
// 		case errors.Is(err, models.ErrEditConflict):
// 			app.editConflictResponse(w, r)
// 		default:
// 			app.serverErrorResponse(w, r, err)
// 		}
// 		return
// 	}
// 	// Write the updated movie record in a JSON response.
// 	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
// 	if err != nil {
// 		app.serverErrorResponse(w, r, err)
// 	}
// }

// func (app *application) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
// 	id, err := app.readUUIDParam(r)
// 	if err != nil {
// 		app.notFoundResponse(w, r)
// 		return
// 	}
// 	// Delete the user from the database, sending a 404 Not Found response to the
// 	// client if there isn't a matching record.
// 	err = app.models.Users.Delete(id)
// 	if err != nil {
// 		switch {
// 		case errors.Is(err, models.ErrRecordNotFound):
// 			app.notFoundResponse(w, r)
// 		default:
// 			app.serverErrorResponse(w, r, err)
// 		}
// 		return
// 	}
// 	// Return a 200 OK status code along with a success message.
// 	err = app.writeJSON(w, http.StatusOK, envelope{"message": "user successfully deleted"}, nil)
// 	if err != nil {
// 		app.serverErrorResponse(w, r, err)
// 	}
// }

// func (app *application) listUsersHandler(w http.ResponseWriter, r *http.Request) {
// 	// To keep things consistent with our other handlers, we'll define an input struct
// 	// to hold the expected values from the request query string.
// 	var input struct {
// 		FirstName string
// 		LastName  string
// 		Email     string
// 		data.Filters
// 	}
// 	// Initialize a new Validator instance.
// 	v := validator.New()
// 	// Call r.URL.Query() to get the url.Values map containing the query string data.
// 	qs := r.URL.Query()
// 	// Use our helpers to extract the title and genres query string values, falling back
// 	// to defaults of an empty string and an empty slice respectively if they are not
// 	// provided by the client.
// 	input.FirstName = app.readString(qs, "firstName", "")
// 	input.LastName = app.readString(qs, "lastName", "")
// 	// Get the page and page_size query string values as integers. Notice that we set
// 	// the default page value to 1 and default page_size to 20, and that we pass the
// 	// validator instance as the final argument here.
// 	input.Email = app.readString(qs, "email", "")
// 	// input.DOB = app.readDate(qs, "dob", time.Time{}, v)
// 	input.Filters.Page = app.readInt(qs, "page", 1, v)
// 	input.Filters.PageSize = app.readInt(qs, "pageSize", 20, v)
// 	// Extract the sort query string value, falling back to "id" if it is not provided
// 	// by the client (which will imply a ascending sort on movie ID).
// 	input.Filters.Sort = app.readString(qs, "sort", "createdAt")
// 	input.Filters.SortSafelist = []string{"createdAt", "id", "firstName", "lastName", "email", "dob", "-id", "-firstname", "-lastName", "-email", "-dob"}

// 	// Check the Validator instance for any errors and use the failedValidationResponse()
// 	// helper to send the client a response if necessary.
// 	if data.ValidateFilters(v, input.Filters); !v.Valid() {
// 		app.failedValidationResponse(w, r, v.Errors)
// 		return
// 	}

// 	// Call the GetAll() method to retrieve the movies, passing in the various filter
// 	// parameters.
// 	users, metadata, err := app.models.Users.GetAll(input.FirstName, input.LastName, input.Email, input.Filters)
// 	if err != nil {
// 		app.serverErrorResponse(w, r, err)
// 		return
// 	}
// 	// Send a JSON response containing the movie data.
// 	err = app.writeJSON(w, http.StatusOK, envelope{"users": users, "metadata": metadata}, nil)
// 	if err != nil {
// 		app.serverErrorResponse(w, r, err)
// 	}
// }
