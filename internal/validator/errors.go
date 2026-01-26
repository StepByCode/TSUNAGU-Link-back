package validator

// FieldError represents a validation error for a specific field
type FieldError struct {
	Message string `json:"message"`
	Value   any    `json:"value,omitempty"`
}

// ValidationErrors is a map of field names to their validation errors
type ValidationErrors map[string]FieldError

// ToResponse converts validation errors to HTTP response format
func (ve ValidationErrors) ToResponse() map[string]interface{} {
	return map[string]interface{}{
		"error":  "validation failed",
		"fields": ve,
	}
}

// HasErrors returns true if there are any validation errors
func (ve ValidationErrors) HasErrors() bool {
	return len(ve) > 0
}
