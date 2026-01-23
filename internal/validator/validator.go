package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validator wraps go-playground/validator with custom error handling
type Validator struct {
	validate *validator.Validate
}

// NewValidator creates a new Validator instance
func NewValidator() *Validator {
	return &Validator{
		validate: validator.New(),
	}
}

// Validate performs validation on the given data and returns ValidationErrors
func (v *Validator) Validate(data interface{}) ValidationErrors {
	errs := ValidationErrors{}

	err := v.validate.Struct(data)
	if err == nil {
		return errs
	}

	// Type assertion to validator.ValidationErrors
	validationErrs, ok := err.(validator.ValidationErrors)
	if !ok {
		// Unexpected error type
		errs["_general"] = FieldError{
			Message: "validation error",
		}
		return errs
	}

	// Convert validator.ValidationErrors to our custom format
	for _, fieldErr := range validationErrs {
		fieldName := fieldErr.Field()
		// Convert field name to lowercase for JSON consistency
		fieldName = strings.ToLower(string(fieldName[0])) + fieldName[1:]

		errs[fieldName] = FieldError{
			Message: v.formatErrorMessage(fieldErr),
			Value:   fieldErr.Value(),
		}
	}

	return errs
}

// formatErrorMessage converts validator.FieldError to a human-readable message
func (v *Validator) formatErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "this field is required"
	case "email":
		return "must be a valid email address"
	case "min":
		if fe.Type().String() == "string" {
			return fmt.Sprintf("must be at least %s characters long", fe.Param())
		}
		return fmt.Sprintf("must be at least %s", fe.Param())
	case "max":
		if fe.Type().String() == "string" {
			return fmt.Sprintf("must be at most %s characters long", fe.Param())
		}
		return fmt.Sprintf("must be at most %s", fe.Param())
	case "len":
		return fmt.Sprintf("must be exactly %s characters long", fe.Param())
	case "eq":
		return fmt.Sprintf("must be equal to %s", fe.Param())
	case "ne":
		return fmt.Sprintf("must not be equal to %s", fe.Param())
	case "gt":
		return fmt.Sprintf("must be greater than %s", fe.Param())
	case "gte":
		return fmt.Sprintf("must be greater than or equal to %s", fe.Param())
	case "lt":
		return fmt.Sprintf("must be less than %s", fe.Param())
	case "lte":
		return fmt.Sprintf("must be less than or equal to %s", fe.Param())
	case "oneof":
		return fmt.Sprintf("must be one of: %s", fe.Param())
	case "uuid":
		return "must be a valid UUID"
	case "url":
		return "must be a valid URL"
	case "uri":
		return "must be a valid URI"
	case "alpha":
		return "must contain only alphabetic characters"
	case "alphanum":
		return "must contain only alphanumeric characters"
	case "numeric":
		return "must be a numeric value"
	default:
		return fmt.Sprintf("validation failed on '%s' tag", fe.Tag())
	}
}
