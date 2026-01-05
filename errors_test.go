package toml

// Copyright (C) 2024 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

import (
	"errors"
	"testing"

	"tideland.dev/go/asserts/verify"
)

// TestErrorCodeString tests ErrorCode String method.
func TestErrorCodeString(t *testing.T) {
	tests := []struct {
		code ErrorCode
		want string
	}{
		{ErrNone, "none"},
		{ErrSyntax, "syntax"},
		{ErrType, "type"},
		{ErrNotFound, "not found"},
		{ErrDuplicate, "duplicate"},
		{ErrInvalidPath, "invalid path"},
		{ErrIO, "io"},
		{ErrInvalid, "invalid"},
		{ErrorCode(999), "unknown"},
	}

	for _, tt := range tests {
		got := tt.code.String()
		verify.Equal(t, got, tt.want)
	}
}

// TestParseErrorError tests ParseError Error method.
func TestParseErrorError(t *testing.T) {
	tests := []struct {
		name string
		err  *ParseError
	}{
		{
			"with line and path",
			&ParseError{
				Op:   "parse",
				Path: "test.toml",
				Line: 10,
				Col:  5,
				Err:  errors.New("test error"),
				Code: ErrSyntax,
			},
		},
		{
			"with line no path",
			&ParseError{
				Op:   "lexer",
				Line: 5,
				Col:  3,
				Err:  errors.New("test error"),
				Code: ErrSyntax,
			},
		},
		{
			"with path no line",
			&ParseError{
				Op:   "get",
				Path: "key.path",
				Err:  errors.New("not found"),
				Code: ErrNotFound,
			},
		},
		{
			"no path no line",
			&ParseError{
				Op:   "validate",
				Err:  errors.New("invalid"),
				Code: ErrInvalid,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := tt.err.Error()
			verify.True(t, len(msg) > 0)
		})
	}
}

// TestParseErrorUnwrap tests ParseError Unwrap method.
func TestParseErrorUnwrap(t *testing.T) {
	innerErr := errors.New("inner error")
	parseErr := &ParseError{
		Op:   "test",
		Err:  innerErr,
		Code: ErrSyntax,
	}

	unwrapped := parseErr.Unwrap()
	verify.Equal(t, unwrapped, innerErr)
}

// TestNewErrorConstructors tests error constructor functions.
func TestNewErrorConstructors(t *testing.T) {
	err1 := NewError("op1", errors.New("err1"), ErrSyntax)
	verify.Equal(t, err1.Op, "op1")
	verify.Equal(t, err1.Code, ErrSyntax)

	err2 := NewErrorWithPath("op2", "path2", errors.New("err2"), ErrNotFound)
	verify.Equal(t, err2.Op, "op2")
	verify.Equal(t, err2.Path, "path2")
	verify.Equal(t, err2.Code, ErrNotFound)

	err3 := NewErrorWithLocation("op3", "path3", 10, 5, errors.New("err3"), ErrInvalid)
	verify.Equal(t, err3.Op, "op3")
	verify.Equal(t, err3.Path, "path3")
	verify.Equal(t, err3.Line, 10)
	verify.Equal(t, err3.Col, 5)
	verify.Equal(t, err3.Code, ErrInvalid)
}
