package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Danik14/library/internal/data"
	"github.com/Danik14/library/internal/models"
	"github.com/Danik14/library/internal/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func (app *application) createBookHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Title  string   `json:"title"`
		Author string   `json:"author"`
		Year   string   `json:"year"`
		Pages  string   `json:"pages"`
		Genres []string `json:"genres"`
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

	var ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	book.ID = primitive.NewObjectID()
	_, anyerr := app.models.Books.DB.InsertOne(ctx, book)
	if anyerr != nil {
		app.serverErrorResponse(w, r, anyerr)
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/books/%d", book.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"book": book}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	fmt.Println(book)
}

func (app *application) showBookHandler(w http.ResponseWriter, r *http.Request) {
	objectId, err := app.readPrimitiveObjectIdParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	// Call the Get() method to fetch the data for a specific book. We also need to
	// use the errors.Is() function to check if it returns a data.ErrRecordNotFound
	// error, in which case we send a 404 Not Found response to the client.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	book := models.Book{}

	err = app.models.Books.DB.FindOne(ctx, bson.M{"_id": objectId}).Decode(&book)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			app.notFoundResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
			app.logError(r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"book": book}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listBooksHandler(w http.ResponseWriter, r *http.Request) {
	// To keep things consistent with our other handlers, we'll define an input struct
	// to hold the expected values from the request query string.
	var input struct {
		Title  string
		Author string
		// Year     int
		Genres []string
		data.Filters
	}
	// Initialize a new Validator instance.
	v := validator.New()
	// Call r.URL.Query() to get the url.Values map containing the query string data.
	qs := r.URL.Query()
	// Use our helpers to extract the title and genres query string values, falling back
	// to defaults of an empty string and an empty slice respectively if they are not
	// provided by the client.
	input.Title = app.readString(qs, "title", "")
	input.Author = app.readString(qs, "author", "")
	// input.Year = app.readInt(qs, "year", 0, v)
	input.Genres = app.readCSV(qs, "genres", []string{})
	// Get the page and page_size query string values as integers. Notice that we set
	// the default page value to 1 and default page_size to 20, and that we pass the
	// validator instance as the final argument here.
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	// Extract the sort query string value, falling back to "id" if it is not provided
	// by the client (which will imply a ascending sort on movie ID).
	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "title", "author", "year", "runtime", "-id", "-title", "-author", "-year", "-runtime"}

	// Check the Validator instance for any errors and use the failedValidationResponse()
	// helper to send the client a response if necessary.
	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	var ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	searchquerydb, err := app.models.Books.DB.Find(ctx, bson.M{})
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	books := []models.Book{}
	err = searchquerydb.All(ctx, &books)
	if err != nil {
		app.logError(r, err)
		app.serverErrorResponse(w, r, err)
		return
	}
	defer searchquerydb.Close(ctx)

	// Send a JSON response containing the movie data.
	err = app.writeJSON(w, http.StatusOK, envelope{"books": books}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {
	// user, err := models.NewUser("Danik", "Slave", "danik_slave@gmail.com", "123", time.Now(), 1)
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
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(input.Password), 10)
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
		HashedPassword: hashedPass,
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

	var ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now()
	_, anyerr := app.models.Users.DB.InsertOne(ctx, user)
	if anyerr != nil {
		app.serverErrorResponse(w, r, anyerr)
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/user/%d", user.ID))

	// Write a JSON response with a 201 Created status code, the movie data in the
	// response body, and the Location header.
	err = app.writeJSON(w, http.StatusCreated, envelope{"user": user}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showUserHandler(w http.ResponseWriter, r *http.Request) {
	objectId, err := app.readPrimitiveObjectIdParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	// Call the Get() method to fetch the data for a specific book. We also need to
	// use the errors.Is() function to check if it returns a data.ErrRecordNotFound
	// error, in which case we send a 404 Not Found response to the client.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	cursor, err := app.models.Users.DB.Find(ctx, bson.M{"_id": objectId})
	if err != nil {
		app.serverErrorResponse(w, r, err)
		app.logError(r, err)
		return
	}

	user := []models.User{}

	err = cursor.All(ctx, &user)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		app.logError(r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) sign_in(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string
		Password string
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.logError(r, err)
		app.badRequestResponse(w, r, err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user models.User
	app.models.Users.DB.FindOne(ctx, bson.M{"email": input.Email}).Decode(&user)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		app.logError(r, err)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(input.Password))
	if err == nil {
		app.writeJSON(w, http.StatusOK, envelope{"message": "Succesfully authorized!"}, nil)
	} else {
		app.logError(r, err)
		app.errorResponse(w, r, http.StatusUnauthorized, "Invalid email or password!")
		return
	}
}

func (app *application) listUsersHandler(w http.ResponseWriter, r *http.Request) {
	// To keep things consistent with our other handlers, we'll define an input struct
	// to hold the expected values from the request query string.
	var input struct {
		FirstName string
		LastName  string
		Email     string
		data.Filters
	}
	// Initialize a new Validator instance.
	v := validator.New()
	// Call r.URL.Query() to get the url.Values map containing the query string data.
	qs := r.URL.Query()
	// Use our helpers to extract the title and genres query string values, falling back
	// to defaults of an empty string and an empty slice respectively if they are not
	// provided by the client.
	input.FirstName = app.readString(qs, "firstName", "")
	input.LastName = app.readString(qs, "lastName", "")
	// Get the page and page_size query string values as integers. Notice that we set
	// the default page value to 1 and default page_size to 20, and that we pass the
	// validator instance as the final argument here.
	input.Email = app.readString(qs, "email", "")
	// input.DOB = app.readDate(qs, "dob", time.Time{}, v)
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "pageSize", 20, v)
	// Extract the sort query string value, falling back to "id" if it is not provided
	// by the client (which will imply a ascending sort on movie ID).
	input.Filters.Sort = app.readString(qs, "sort", "createdAt")
	input.Filters.SortSafelist = []string{"createdAt", "id", "firstName", "lastName", "email", "dob", "-id", "-firstname", "-lastName", "-email", "-dob"}

	// Check the Validator instance for any errors and use the failedValidationResponse()
	// helper to send the client a response if necessary.
	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	var ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	searchquerydb, err := app.models.Users.DB.Find(ctx, bson.M{})
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	users := []models.User{}
	err = searchquerydb.All(ctx, &users)
	if err != nil {
		app.logError(r, err)
		app.serverErrorResponse(w, r, err)
		return
	}
	defer searchquerydb.Close(ctx)

	err = app.writeJSON(w, http.StatusOK, envelope{"users": users}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readPrimitiveObjectIdParam(r)

	// Extract the movie ID from the URL.

	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	if err != nil {
		app.notFoundResponse(w, r)
	}
	// Declare an input struct to hold the expected data from the client.
	var input struct {
		FirstName *string                        `json:"firstName"`
		LastName  *string                        `json:"lastName"`
		Email     *string                        `json:"email"`
		Password  *string                        `json:"password"`
		Marks     *map[primitive.ObjectID]string `json:"marks"`
		DOB       *time.Time                     `json:"dob"` // date of birth
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	var ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var user models.User
	app.models.Users.DB.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if input.FirstName != nil {
		user.FirstName = *input.FirstName
	}
	if input.LastName != nil {
		user.LastName = *input.LastName
	}
	if input.Email != nil {
		user.Email = *input.Email
	}
	if input.Password != nil {
		hashedPass, err := bcrypt.GenerateFromPassword([]byte(*input.Password), 10)
		if err != nil {
			app.logError(r, err)
			app.badRequestResponse(w, r, err)
			return
		}
		user.HashedPassword = hashedPass
	}
	if input.Marks != nil {
		user.Marks = *input.Marks
	}
	if input.DOB != nil {
		user.DOB = *input.DOB
	}

	// Validate the updated movie record, sending the client a 422 Unprocessable Entity
	// response if any checks fail.
	v := validator.New()
	if models.ValidateUser(v, &user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	cursor, err := app.models.Users.DB.UpdateByID(ctx, id, bson.M{"$set": user})
	if err != nil {

		switch {
		case errors.Is(err, models.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// Write the updated movie record in a JSON response.
	err = app.writeJSON(w, http.StatusOK, envelope{"user": cursor}, nil)
}
func (app *application) updateBookHandler(w http.ResponseWriter, r *http.Request) {
	objectId, err := app.readPrimitiveObjectIdParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	// Fetch the existing movie record from the database, sending a 404 Not Found
	// response to the client if we couldn't find a matching record.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	book := models.Book{}

	err = app.models.Books.DB.FindOne(ctx, bson.M{"_id": objectId}).Decode(&book)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			app.notFoundResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
			app.logError(r, err)
		}
		return
	}

	// Declare an input struct to hold the expected data from the client.
	var input struct {
		Title  *string  `json:"title"`
		Author *string  `json:"author"`
		Year   *string  `json:"year"`
		Pages  *string  `json:"pages"`
		Genres []string `json:"genres"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// If the input.Title value is nil then we know that no corresponding "title" key/
	// value pair was provided in the JSON request body. So we move on and leave the
	// movie record unchanged. Otherwise, we update the movie record with the new title
	// value. Importantly, because input.Title is a now a pointer to a string, we need
	// to dereference the pointer using the * operator to get the underlying value
	// before assigning it to our movie record.
	if input.Title != nil {
		book.Title = *input.Title
	}
	if input.Author != nil {
		book.Author = *input.Author
	}
	if input.Year != nil {
		book.Year = *input.Year
	}
	if input.Pages != nil {
		book.Pages = *input.Pages
	}
	if input.Genres != nil {
		book.Genres = input.Genres
	}

	// Validate the updated movie record, sending the client a 422 Unprocessable Entity
	// response if any checks fail.
	v := validator.New()
	if models.ValidateBook(v, &book); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Update the book record in the database.
	update := bson.M{
		"$set": bson.M{
			"title":  book.Title,
			"author": book.Author,
			"year":   book.Year,
			"pages":  book.Pages,
			"genres": book.Genres,
		},
	}

	result, err := app.models.Books.DB.UpdateOne(
		ctx,
		bson.M{"_id": objectId},
		update,
	)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		app.logError(r, err)
		return
	}

	// Check if no documents were updated.
	if result.ModifiedCount == 0 {
		app.notFoundResponse(w, r)
		return
	}

	// Write the updated movie record in a JSON response.
	err = app.writeJSON(w, http.StatusOK, envelope{"book": book}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteBookHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the movie ID from the URL.
	objectId, err := app.readPrimitiveObjectIdParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	// Delete the book from the database.
	result, err := app.models.Books.DB.DeleteOne(context.Background(), bson.M{"_id": objectId})
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// If no documents were deleted, send a 404 Not Found response.
	if result.DeletedCount == 0 {
		app.notFoundResponse(w, r)
		return
	}

	// Return a 200 OK status code along with a success message.
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "book successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readPrimitiveObjectIdParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	// Delete the user from the database, sending a 404 Not Found response to the
	// client if there isn't a matching record.
	var ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err = app.models.Users.DB.DeleteOne(ctx, bson.M{"_id": id})
	// fmt.Println(cursor)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}

		return
	}
	// Return a 200 OK status code along with a success message.
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "user successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
