package gotv

import "errors"

// Predefined errors for GOTV-related operations.
var (
	ErrInvalidAuth      = errors.New("invalid authentication")
	ErrFragmentNotFound = errors.New("fragment not found")
	ErrMatchNotFound    = errors.New("match not found")
)
