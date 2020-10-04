package web

import (
	"encoding/json"
	"github.com/pkg/errors"
	"net/http"
)

// Respond responds to the client with marshalled val
func Respond(w http.ResponseWriter, val interface{}, statusCode int) error {
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
