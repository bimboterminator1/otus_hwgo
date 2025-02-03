package storage

import (
	"errors"
	"fmt"
)

var (
	// Base errors.
	ErrNotFound          = errors.New("not found")
	ErrAlreadyExists     = errors.New("already exists")
	ErrDatabaseOperation = errors.New("database operation failed")
	ErrValidation        = errors.New("validation failed")
)

// Event-specific errors.
func EventNotFound(id int64) error {
	return fmt.Errorf("event %d: %w", id, ErrNotFound)
}

func EventAlreadyExists(id int64) error {
	return fmt.Errorf("event %d: %w", id, ErrAlreadyExists)
}

// Database-specific errors.
func DatabaseError(op string, err error) error {
	return fmt.Errorf("%s: %w: %w", op, ErrDatabaseOperation, err)
}

func ValidationError(msg string) error {
	return fmt.Errorf("%s: %w", msg, ErrValidation)
}

// Error checking helpers.
func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}

func IsAlreadyExists(err error) bool {
	return errors.Is(err, ErrAlreadyExists)
}

func IsDatabaseError(err error) bool {
	return errors.Is(err, ErrDatabaseOperation)
}

func IsValidationError(err error) bool {
	return errors.Is(err, ErrValidation)
}
