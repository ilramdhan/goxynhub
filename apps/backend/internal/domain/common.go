package domain

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

// JSONMap is a custom type for JSONB columns
type JSONMap map[string]interface{}

// Value implements the driver.Valuer interface for database serialization
func (j JSONMap) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	b, err := json.Marshal(j)
	if err != nil {
		return nil, fmt.Errorf("JSONMap.Value: %w", err)
	}
	return string(b), nil
}

// Scan implements the sql.Scanner interface for database deserialization
func (j *JSONMap) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New("JSONMap.Scan: unsupported type")
	}
	return json.Unmarshal(bytes, j)
}

// JSONArray is a custom type for JSONB array columns
type JSONArray []interface{}

// Value implements the driver.Valuer interface
func (j JSONArray) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	b, err := json.Marshal(j)
	if err != nil {
		return nil, fmt.Errorf("JSONArray.Value: %w", err)
	}
	return string(b), nil
}

// Scan implements the sql.Scanner interface
func (j *JSONArray) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New("JSONArray.Scan: unsupported type")
	}
	return json.Unmarshal(bytes, j)
}

// StringArray is a custom type for text[] columns
type StringArray []string

// Value implements the driver.Valuer interface
func (s StringArray) Value() (driver.Value, error) {
	if s == nil {
		return nil, nil
	}
	b, err := json.Marshal(s)
	if err != nil {
		return nil, fmt.Errorf("StringArray.Value: %w", err)
	}
	return string(b), nil
}

// Scan implements the sql.Scanner interface
func (s *StringArray) Scan(value interface{}) error {
	if value == nil {
		*s = nil
		return nil
	}
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New("StringArray.Scan: unsupported type")
	}
	return json.Unmarshal(bytes, s)
}

// Pagination holds pagination parameters
type Pagination struct {
	Page    int `json:"page" form:"page"`
	PerPage int `json:"per_page" form:"per_page"`
}

// Normalize ensures pagination values are within valid ranges
func (p *Pagination) Normalize() {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.PerPage < 1 || p.PerPage > 100 {
		p.PerPage = 20
	}
}

// Offset returns the SQL offset for the current page
func (p *Pagination) Offset() int {
	return (p.Page - 1) * p.PerPage
}

// PaginatedResult holds paginated query results
type PaginatedResult[T any] struct {
	Data       []T `json:"data"`
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// NewPaginatedResult creates a new PaginatedResult
func NewPaginatedResult[T any](data []T, total int, pagination Pagination) PaginatedResult[T] {
	totalPages := total / pagination.PerPage
	if total%pagination.PerPage > 0 {
		totalPages++
	}
	return PaginatedResult[T]{
		Data:       data,
		Page:       pagination.Page,
		PerPage:    pagination.PerPage,
		Total:      total,
		TotalPages: totalPages,
	}
}

// Common domain errors
var (
	ErrNotFound          = errors.New("resource not found")
	ErrAlreadyExists     = errors.New("resource already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrAccountLocked     = errors.New("account is temporarily locked")
	ErrAccountInactive   = errors.New("account is inactive")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrForbidden         = errors.New("forbidden")
	ErrInvalidToken      = errors.New("invalid or expired token")
	ErrTokenRevoked      = errors.New("token has been revoked")
	ErrValidation        = errors.New("validation error")
)
