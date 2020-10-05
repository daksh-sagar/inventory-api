package web

import "net/http"

type ErrorResponse struct {
	Error string `json:"error"`
}

type Error struct {
	Err    error
	Status int
}

func (e *Error) Error() string {
	return e.Err.Error()
}

// NewRequestError is used when a known error is encountered
func NewRequestError(err error, status int) error {
	return &Error{
		Err:    err,
		Status: status,
	}
}

// RespondError knows how to handle errors going to the client
func RespondError(w http.ResponseWriter, err error) error {
	if webErr, ok := err.(*Error); ok {
		res := ErrorResponse{
			Error: webErr.Error(),
		}
		return Respond(w, res, webErr.Status)
	}

	res := ErrorResponse{
		Error: http.StatusText(http.StatusInternalServerError),
	}
	return Respond(w, res, http.StatusInternalServerError)
}
