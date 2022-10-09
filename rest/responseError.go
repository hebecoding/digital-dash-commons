package rest

type ErrorResponse struct {
	Status       int    `json:"status,omitempty"`
	ErrorMessage string `json:"errorMessage"`
	Error        error
}

func (e *ErrorResponse) String() string {
	return e.ErrorMessage
}
