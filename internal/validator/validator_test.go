package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestStruct struct {
	Email    string `validate:"required,email"`
	Name     string `validate:"required"`
	Password string `validate:"required,min=8"`
	Age      int    `validate:"gte=0,lte=120"`
}

func TestValidator_Validate_Success(t *testing.T) {
	v := NewValidator()

	data := TestStruct{
		Email:    "test@example.com",
		Name:     "Test User",
		Password: "password123",
		Age:      25,
	}

	errs := v.Validate(&data)

	assert.Empty(t, errs)
	assert.False(t, errs.HasErrors())
}

func TestValidator_Validate_RequiredField(t *testing.T) {
	v := NewValidator()

	data := TestStruct{
		Email:    "",
		Name:     "Test User",
		Password: "password123",
		Age:      25,
	}

	errs := v.Validate(&data)

	require.True(t, errs.HasErrors())
	assert.Contains(t, errs, "email")
	assert.Equal(t, "this field is required", errs["email"].Message)
}

func TestValidator_Validate_InvalidEmail(t *testing.T) {
	v := NewValidator()

	data := TestStruct{
		Email:    "invalid-email",
		Name:     "Test User",
		Password: "password123",
		Age:      25,
	}

	errs := v.Validate(&data)

	require.True(t, errs.HasErrors())
	assert.Contains(t, errs, "email")
	assert.Equal(t, "must be a valid email address", errs["email"].Message)
	assert.Equal(t, "invalid-email", errs["email"].Value)
}

func TestValidator_Validate_MinLength(t *testing.T) {
	v := NewValidator()

	data := TestStruct{
		Email:    "test@example.com",
		Name:     "Test User",
		Password: "pass",
		Age:      25,
	}

	errs := v.Validate(&data)

	require.True(t, errs.HasErrors())
	assert.Contains(t, errs, "password")
	assert.Equal(t, "must be at least 8 characters long", errs["password"].Message)
}

func TestValidator_Validate_NumericRange(t *testing.T) {
	v := NewValidator()

	data := TestStruct{
		Email:    "test@example.com",
		Name:     "Test User",
		Password: "password123",
		Age:      150,
	}

	errs := v.Validate(&data)

	require.True(t, errs.HasErrors())
	assert.Contains(t, errs, "age")
	assert.Contains(t, errs["age"].Message, "must be less than or equal to")
}

func TestValidator_Validate_MultipleErrors(t *testing.T) {
	v := NewValidator()

	data := TestStruct{
		Email:    "invalid",
		Name:     "",
		Password: "short",
		Age:      -5,
	}

	errs := v.Validate(&data)

	require.True(t, errs.HasErrors())
	assert.Len(t, errs, 4)
	assert.Contains(t, errs, "email")
	assert.Contains(t, errs, "name")
	assert.Contains(t, errs, "password")
	assert.Contains(t, errs, "age")
}

func TestValidationErrors_ToResponse(t *testing.T) {
	errs := ValidationErrors{
		"email": FieldError{
			Message: "must be a valid email address",
			Value:   "invalid",
		},
		"password": FieldError{
			Message: "must be at least 8 characters long",
		},
	}

	response := errs.ToResponse()

	assert.Equal(t, "validation failed", response["error"])
	assert.Equal(t, errs, response["fields"])
}

func TestValidationErrors_HasErrors(t *testing.T) {
	tests := []struct {
		name     string
		errors   ValidationErrors
		expected bool
	}{
		{
			name:     "No errors",
			errors:   ValidationErrors{},
			expected: false,
		},
		{
			name: "Has errors",
			errors: ValidationErrors{
				"email": FieldError{Message: "invalid"},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.errors.HasErrors())
		})
	}
}
