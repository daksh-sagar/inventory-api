package web

// FieldError is used to indicate an error with a specific request field.
type FieldError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

type ErrorResponse struct {
	Error  string       `json:"error"`
	Fields []FieldError `json:"fields,omitempty"`
}

type Error struct {
	Err    error
	Status int
	Fields []FieldError
}

func (e *Error) Error() string {
	return e.Err.Error()
}

// NewRequestError is used when a known error is encountered
func NewRequestError(err error, status int) error {
	return &Error{
		err, status, nil,
	}
}
