// All material is licensed under the Apache License Version 2.0, January 2004
// http://www.apache.org/licenses/LICENSE-2.0

// Current Status Codes:
//		200 OK           : StatusOK                  : Call is success and returning data.
//		204 No Content   : StatusNoContent           : Call is success and returns no data.
//		400 Bad Request  : StatusBadRequest          : Invalid post data (syntax or semantics).
//		401 Unauthorized : StatusUnauthorized        : Authentication failure.
//		404 Not Found    : StatusNotFound            : Invalid URL or identifier.
//		500 Internal     : StatusInternalServerError : Application specific beyond scope of user.

package web

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/pkg/errors"
)

// Invalid describes a validation error belonging to a specific field.
type Invalid struct {
	Fld string `json:"field_name"`
	Err string `json:"error"`
}

// InvalidError is a custom error type for invalid fields.
type InvalidError []Invalid

// Error implements the error interface for InvalidError.
func (err InvalidError) Error() string {
	var str string
	for _, v := range err {
		str = fmt.Sprintf("%s,{%s:%s}", str, v.Fld, v.Err)
	}
	return str
}

// JSONError is the response for errors that occur within the API.
type JSONError struct {
	Error  string       `json:"error"`
	Fields InvalidError `json:"fields,omitempty"`
}

var (
	// ErrNotAuthorized occurs when the call is not authorized.
	ErrNotAuthorized = errors.New("Not authorized")

	// ErrDBNotConfigured occurs when the DB is not initialized.
	ErrDBNotConfigured = errors.New("DB not initialized")

	// ErrNotFound is abstracting the mgo not found error.
	ErrNotFound = errors.New("Entity not found")

	// ErrInvalidID occurs when an ID is not in a valid form.
	ErrInvalidID = errors.New("ID is not in it's proper form")

	// ErrValidation occurs when there are validation errors.
	ErrValidation = errors.New("Validation errors occurred")
)

// Error handles all error responses for the API.
func Error(cxt context.Context, w http.ResponseWriter, err error) {
	switch errors.Cause(err) {
	case ErrNotFound:
		RespondError(cxt, w, err, http.StatusNotFound)
		return

	case ErrInvalidID:
		RespondError(cxt, w, err, http.StatusBadRequest)
		return

	case ErrValidation:
		RespondError(cxt, w, err, http.StatusBadRequest)
		return

	case ErrNotAuthorized:
		RespondError(cxt, w, err, http.StatusUnauthorized)
		return
	}

	switch e := errors.Cause(err).(type) {
	case InvalidError:
		v := JSONError{
			Error:  "field validation failure",
			Fields: e,
		}

		http.Error(w, "Invalid"+v.Error, http.StatusBadRequest)
		return
	}

	RespondError(cxt, w, err, http.StatusInternalServerError)
}

// RespondError sends JSON describing the error
func RespondError(ctx context.Context, w http.ResponseWriter, err error, code int) {

	// Set the status code for the request logger middleware.
	v := ctx.Value(KeyValues).(*Values)
	log.Printf("ERROR (%d) : %s",
		v.StatusCode,
		v.TraceID,
	)

	http.Error(w, "Something bad happened. \n"+v.TraceID+" <-- This error id can help us resolve the issue for you. \nPlease try again or, write to us at //TODO MAIL ID", code)

}

// Respond sends JSON to the client.
// If code is StatusNoContent, v is expected to be nil.
func Respond(ctx context.Context, w http.ResponseWriter, name string, data interface{}, code int) {

	// Set the status code for the request logger middleware.
	v := ctx.Value(KeyValues).(*Values)
	v.StatusCode = code

	// Just set the status code and we are done.
	if code == http.StatusNoContent {
		w.WriteHeader(code)
		return
	}

	// Set the content type.
	w.Header().Set("Content-Type", "text/html")

	// Write the status code to the response and context.
	//w.WriteHeader(code)

	err := Execute(w, name, data)
	if err != nil {
		//http.Error(w, "Internal Server Error", 500)
		// If template is unable to execute, return a internal server error
		RespondError(ctx, w, err, http.StatusInternalServerError)
	}
}
