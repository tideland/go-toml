package toml

// Copyright (C) 2024 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

import "fmt"

// ErrorCode represents different categories of errors that can occur during TOML parsing.
type ErrorCode int

const (
	// ErrNone indicates no error occurred.
	ErrNone ErrorCode = iota
	// ErrSyntax indicates a syntax error in the TOML file.
	ErrSyntax
	// ErrType indicates a type conversion error.
	ErrType
	// ErrNotFound indicates a key or section was not found.
	ErrNotFound
	// ErrDuplicate indicates a duplicate key or table definition.
	ErrDuplicate
	// ErrInvalidPath indicates an invalid dot-notation path.
	ErrInvalidPath
	// ErrIO indicates a file I/O error.
	ErrIO
	// ErrInvalid indicates an invalid configuration or parameter.
	ErrInvalid
)

// String returns the string representation of the error code.
func (ec ErrorCode) String() string {
	switch ec {
	case ErrNone:
		return "none"
	case ErrSyntax:
		return "syntax"
	case ErrType:
		return "type"
	case ErrNotFound:
		return "not found"
	case ErrDuplicate:
		return "duplicate"
	case ErrInvalidPath:
		return "invalid path"
	case ErrIO:
		return "io"
	case ErrInvalid:
		return "invalid"
	default:
		return "unknown"
	}
}

// ParseError represents an error that occurred during TOML parsing or value access.
type ParseError struct {
	// Op is the operation that failed (e.g., "parse", "get", "lexer").
	Op string
	// Path is the TOML path where the error occurred (e.g., "database.connection.timeout").
	Path string
	// Line is the line number where the error occurred (0 if not applicable).
	Line int
	// Col is the column number where the error occurred (0 if not applicable).
	Col int
	// Err is the underlying error.
	Err error
	// Code is the error code category.
	Code ErrorCode
}

// Error returns the string representation of the parse error.
func (e *ParseError) Error() string {
	if e.Line > 0 {
		if e.Path != "" {
			return fmt.Sprintf("toml %s at %s:%d:%d: %v (%v)",
				e.Op, e.Path, e.Line, e.Col, e.Err, e.Code)
		}
		return fmt.Sprintf("toml %s at %d:%d: %v (%v)",
			e.Op, e.Line, e.Col, e.Err, e.Code)
	}
	if e.Path != "" {
		return fmt.Sprintf("toml %s at %s: %v (%v)",
			e.Op, e.Path, e.Err, e.Code)
	}
	return fmt.Sprintf("toml %s: %v (%v)", e.Op, e.Err, e.Code)
}

// Unwrap returns the underlying error.
func (e *ParseError) Unwrap() error {
	return e.Err
}

// NewError creates a new ParseError with the given operation, error, and code.
func NewError(op string, err error, code ErrorCode) *ParseError {
	return &ParseError{
		Op:   op,
		Err:  err,
		Code: code,
	}
}

// NewErrorWithPath creates a new ParseError with operation, path, error, and code.
func NewErrorWithPath(op, path string, err error, code ErrorCode) *ParseError {
	return &ParseError{
		Op:   op,
		Path: path,
		Err:  err,
		Code: code,
	}
}

// NewErrorWithLocation creates a new ParseError with full location information.
func NewErrorWithLocation(op, path string, line, col int, err error, code ErrorCode) *ParseError {
	return &ParseError{
		Op:   op,
		Path: path,
		Line: line,
		Col:  col,
		Err:  err,
		Code: code,
	}
}
