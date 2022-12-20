package models

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Pages uint32

// Define an error that our UnmarshalJSON() method can return if we're unable to parse // or convert the JSON string successfully.
var ErrInvalidPagesFormat = errors.New("invalid pages format")

// Implement a MarshalJSON() method on the Runtime type so that it satisfies the
// json.Marshaler interface. This should return the JSON-encoded value for the movie
// runtime (in our case, it will return a string in the format "<runtime> mins").
func (r Pages) MarshalJSON() ([]byte, error) {
	// Generate a string containing the movie runtime in the required format.
	jsonValue := fmt.Sprintf("%d pages", r)
	// Use the strconv.Quote() function on the string to wrap it in double quotes. It
	// needs to be surrounded by double quotes in order to be a valid *JSON string*.
	quotedJSONValue := strconv.Quote(jsonValue)
	// Convert the quoted string value to a byte slice and return it.
	return []byte(quotedJSONValue), nil
}

// Implement a UnmarshalJSON() method on the Runtime type so that it satisfies the
// json.Unmarshaler interface. IMPORTANT: Because UnmarshalJSON() needs to modify the
// receiver (our Runtime type), we must use a pointer receiver for this to work
// correctly. Otherwise, we will only be modifying a copy (which is then discarded when
// this method returns).
func (p *Pages) UnmarshalJSON(jsonValue []byte) error {
	// We expect that the incoming JSON value will be a string in the format
	// "<runtime> mins", and the first thing we need to do is remove the surrounding
	// double-quotes from this string. If we can't unquote it, then we return the
	// ErrInvalidRuntimeFormat error.
	unquotedJSONValue, err := strconv.Unquote(string(jsonValue))
	if err != nil {
		return ErrInvalidPagesFormat
	}
	// Split the string to isolate the part containing the number.
	parts := strings.Split(unquotedJSONValue, " ")
	// Sanity check the parts of the string to make sure it was in the expected format.
	// If it isn't, we return the ErrInvalidRuntimeFormat error again.
	if len(parts) != 2 || parts[1] != "pages" {
		return ErrInvalidPagesFormat
	}
	// Otherwise, parse the string containing the number into an int32. Again, if this
	// fails return the ErrInvalidRuntimeFormat error.
	i, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return ErrInvalidPagesFormat
	}
	// Convert the int32 to a Runtime type and assign this to the receiver. Note that we
	// use the * operator to deference the receiver (which is a pointer to a Runtime
	// type) in order to set the underlying value of the pointer.
	*p = Pages(i)
	return nil
}
