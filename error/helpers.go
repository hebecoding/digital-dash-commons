package error

import "errors"

// Error Definitions
const (
	errEmptyInput = "input cannot be nil/empty"
)

// Errors
var (
	ErrEmptyInput = errors.New(errEmptyInput)
)
