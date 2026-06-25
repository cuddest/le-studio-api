package handler

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

// FieldError represents a structured validation error returned to the client.
type FieldError struct {
	Field   string `json:"field"`
	Rule    string `json:"rule"`
	Param   string `json:"param,omitempty"`
	Message string `json:"message"`
}

// formatValidationErrors converts a validator.ValidationErrors into a JSON-friendly slice.
// If the error is not a validation error, it returns nil.
func formatValidationErrors(err error) []FieldError {
	var ves validator.ValidationErrors
	if !errors.As(err, &ves) {
		return nil
	}
	out := make([]FieldError, 0, len(ves))
	for _, fe := range ves {
		out = append(out, FieldError{
			Field:   fe.Field(),
			Rule:    fe.Tag(),
			Param:   fe.Param(),
			Message: humanizeValidationMessage(fe),
		})
	}
	return out
}

func humanizeValidationMessage(fe validator.FieldError) string {
	field := fe.Field()
	switch fe.Tag() {
	case "required":
		return field + " is required"
	case "min":
		return field + " must be at least " + fe.Param() + " characters"
	case "max":
		return field + " must be at most " + fe.Param() + " characters"
	case "email":
		return field + " must be a valid email"
	case "oneof":
		return field + " must be one of: " + fe.Param()
	case "uuid", "uuid4":
		return field + " must be a valid UUID"
	default:
		return field + " failed " + fe.Tag() + " validation"
	}
}
