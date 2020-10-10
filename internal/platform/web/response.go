package web

import (
	"encoding/json"
	"github.com/pkg/errors"
	"net/http"
)

// Respond responds to the client with marshalled val
func Respond(w http.ResponseWriter, val interface{}, statusCode int) error {
	if statusCode == http.StatusNoContent {
		w.WriteHeader(statusCode)
		return nil
	}

	data, err := json.Marshal(val)

	if err != nil {
		return errors.Wrap(err, "marshalling to json")
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(statusCode)
	if _, err := w.Write(data); err != nil {
		return errors.Wrap(err, "writing to client")
	}

	return nil
}

// RespondError sends an error response back to the client.
func RespondError(w http.ResponseWriter, err error) error {

	// If the error was of the type *Error, the handler has
	// a specific status code and error to return.
	if webErr, ok := errors.Cause(err).(*Error); ok {
		er := ErrorResponse{
			Error:  webErr.Err.Error(),
			Fields: webErr.Fields,
		}
		if err := Respond(w, er, webErr.Status); err != nil {
			return err
		}
		return nil
	}

	// If not, the handler sent any arbitrary error value so use 500.
	er := ErrorResponse{
		Error: http.StatusText(http.StatusInternalServerError),
	}
	if err := Respond(w, er, http.StatusInternalServerError); err != nil {
		return err
	}
	return nil
}
