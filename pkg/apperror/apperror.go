package apperror

import "fmt"

// AppError defines structured app error.
type AppError struct {
	Code    string
	Message string
	Status  int
	Err     error
}

// Error satisfies error interface.
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}
