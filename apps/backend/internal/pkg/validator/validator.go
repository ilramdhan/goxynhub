package validator

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validator wraps the go-playground validator
type Validator struct {
	validate *validator.Validate
}

// ValidationError represents a single field validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// New creates a new Validator instance with custom validators registered
func New() *Validator {
	v := validator.New()

	// Register custom validators
	v.RegisterValidation("slug", validateSlug)       //nolint:errcheck
	v.RegisterValidation("hostname", validateHostname) //nolint:errcheck

	// Use JSON field names instead of struct field names
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return &Validator{validate: v}
}

// Validate validates a struct and returns a slice of ValidationErrors
func (v *Validator) Validate(s interface{}) []ValidationError {
	err := v.validate.Struct(s)
	if err == nil {
		return nil
	}

	var errors []ValidationError
	for _, e := range err.(validator.ValidationErrors) {
		errors = append(errors, ValidationError{
			Field:   e.Field(),
			Message: formatMessage(e),
		})
	}
	return errors
}

// ValidateVar validates a single variable
func (v *Validator) ValidateVar(field interface{}, tag string) error {
	return v.validate.Var(field, tag)
}

// formatMessage creates a human-readable error message for a validation error
func formatMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", e.Field())
	case "email":
		return fmt.Sprintf("%s must be a valid email address", e.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", e.Field(), e.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", e.Field(), e.Param())
	case "url":
		return fmt.Sprintf("%s must be a valid URL", e.Field())
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", e.Field(), e.Param())
	case "slug":
		return fmt.Sprintf("%s must be a valid URL slug (lowercase letters, numbers, hyphens)", e.Field())
	case "hostname":
		return fmt.Sprintf("%s must be a valid hostname", e.Field())
	case "uuid":
		return fmt.Sprintf("%s must be a valid UUID", e.Field())
	case "numeric":
		return fmt.Sprintf("%s must be a number", e.Field())
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", e.Field(), e.Param())
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s", e.Field(), e.Param())
	default:
		return fmt.Sprintf("%s is invalid", e.Field())
	}
}

// validateSlug validates that a string is a valid URL slug
func validateSlug(fl validator.FieldLevel) bool {
	slug := fl.Field().String()
	if slug == "" {
		return true // Let required handle empty
	}
	matched, _ := regexp.MatchString(`^[a-z0-9]+(?:-[a-z0-9]+)*$`, slug)
	return matched
}

// validateHostname validates that a string is a valid hostname
func validateHostname(fl validator.FieldLevel) bool {
	hostname := fl.Field().String()
	if hostname == "" {
		return true // Let required handle empty
	}
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?)*$`, hostname)
	return matched
}
