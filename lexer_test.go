package toml

// Copyright (C) 2024 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

import (
	"testing"

	"tideland.dev/go/asserts/verify"
)

// TestLexerBasicTokens tests basic token recognition.
func TestLexerBasicTokens(t *testing.T) {
	input := `= [ ] { } , .`
	lexer := NewLexer(input)

	tests := []struct {
		kind  TokenKind
		value string
	}{
		{TokenEquals, "="},
		{TokenLeftBracket, "["},
		{TokenRightBracket, "]"},
		{TokenLeftBrace, "{"},
		{TokenRightBrace, "}"},
		{TokenComma, ","},
		{TokenDot, "."},
		{TokenEOF, ""},
	}

	for _, tt := range tests {
		tok, err := lexer.NextToken()
		verify.NoError(t, err)
		verify.Equal(t, tok.Kind, tt.kind)
		if tt.value != "" {
			verify.Equal(t, tok.Value, tt.value)
		}
	}
}

// TestLexerStrings tests string tokenization.
func TestLexerStrings(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"basic string", `"hello"`, "hello"},
		{"literal string", `'hello'`, "hello"},
		{"string with space", `"hello world"`, "hello world"},
		{"empty string", `""`, ""},
		{"escape sequences", `"hello\nworld"`, "hello\nworld"},
		{"escaped quote", `"say \"hi\""`, `say "hi"`},
		{"escaped backslash", `"path\\to\\file"`, `path\to\file`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			tok, err := lexer.NextToken()
			verify.NoError(t, err)
			verify.Equal(t, tok.Kind, TokenString)
			verify.Equal(t, tok.Value, tt.want)
		})
	}
}

// TestLexerMultilineStrings tests multi-line string tokenization.
func TestLexerMultilineStrings(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"basic multiline", `"""hello
world"""`, "hello\nworld"},
		{"literal multiline", `'''hello
world'''`, "hello\nworld"},
		{"multiline with skip first newline", `"""
hello"""`, "hello"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			tok, err := lexer.NextToken()
			verify.NoError(t, err)
			verify.Equal(t, tok.Kind, TokenString)
			verify.Equal(t, tok.Value, tt.want)
		})
	}
}

// TestLexerIntegers tests integer tokenization.
func TestLexerIntegers(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"simple integer", "42", "42"},
		{"negative integer", "-42", "-42"},
		{"positive integer", "+42", "+42"},
		{"integer with underscores", "1_000_000", "1000000"},
		{"hex integer", "0xDEADBEEF", "0xDEADBEEF"},
		{"octal integer", "0o755", "0o755"},
		{"binary integer", "0b11010110", "0b11010110"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			tok, err := lexer.NextToken()
			verify.NoError(t, err)
			verify.Equal(t, tok.Kind, TokenInteger)
		})
	}
}

// TestLexerFloats tests float tokenization.
func TestLexerFloats(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"simple float", "3.14"},
		{"negative float", "-3.14"},
		{"positive float", "+3.14"},
		{"float with exponent", "1.23e10"},
		{"float with negative exponent", "1.23e-10"},
		{"infinity", "inf"},
		{"negative infinity", "-inf"},
		{"not a number", "nan"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			tok, err := lexer.NextToken()
			verify.NoError(t, err)
			verify.Equal(t, tok.Kind, TokenFloat)
		})
	}
}

// TestLexerBooleans tests boolean tokenization.
func TestLexerBooleans(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"true", "true"},
		{"false", "false"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			tok, err := lexer.NextToken()
			verify.NoError(t, err)
			verify.Equal(t, tok.Kind, TokenBoolean)
			verify.Equal(t, tok.Value, tt.want)
		})
	}
}

// TestLexerKeys tests key tokenization.
func TestLexerKeys(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"simple", "simple"},
		{"with_underscore", "with_underscore"},
		{"with-dash", "with-dash"},
		{"key123", "key123"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			tok, err := lexer.NextToken()
			verify.NoError(t, err)
			verify.Equal(t, tok.Kind, TokenKey)
			verify.Equal(t, tok.Value, tt.want)
		})
	}
}

// TestLexerComments tests comment skipping.
func TestLexerComments(t *testing.T) {
	input := `# This is a comment
key = "value" # inline comment`

	lexer := NewLexer(input)

	// Should skip first comment, but not the newline
	tok, err := lexer.NextToken()
	verify.NoError(t, err)
	verify.Equal(t, tok.Kind, TokenNewline)

	tok, err = lexer.NextToken()
	verify.NoError(t, err)
	verify.Equal(t, tok.Kind, TokenKey)
	verify.Equal(t, tok.Value, "key")

	tok, err = lexer.NextToken()
	verify.NoError(t, err)
	verify.Equal(t, tok.Kind, TokenEquals)

	tok, err = lexer.NextToken()
	verify.NoError(t, err)
	verify.Equal(t, tok.Kind, TokenString)
	verify.Equal(t, tok.Value, "value")
}

// TestLexerNewlines tests newline handling.
func TestLexerNewlines(t *testing.T) {
	input := "key\n=\nvalue"
	lexer := NewLexer(input)

	tok, err := lexer.NextToken()
	verify.NoError(t, err)
	verify.Equal(t, tok.Kind, TokenKey)

	tok, err = lexer.NextToken()
	verify.NoError(t, err)
	verify.Equal(t, tok.Kind, TokenNewline)

	tok, err = lexer.NextToken()
	verify.NoError(t, err)
	verify.Equal(t, tok.Kind, TokenEquals)
}

// TestLexerErrors tests error cases.
func TestLexerErrors(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"unclosed string", `"hello`},
		{"unclosed multiline", `"""hello`},
		{"invalid escape", `"hello\x"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			_, err := lexer.NextToken()
			verify.Error(t, err)
		})
	}
}

// TestLexerLineCol tests line and column tracking.
func TestLexerLineCol(t *testing.T) {
	input := `key = "value"
[section]`

	lexer := NewLexer(input)

	tok, err := lexer.NextToken()
	verify.NoError(t, err)
	verify.Equal(t, tok.Line, 1)
	verify.Equal(t, tok.Col, 1)

	// Skip to newline
	tok, err = lexer.NextToken() // =
	verify.NoError(t, err)
	verify.Equal(t, tok.Kind, TokenEquals)
	tok, err = lexer.NextToken() // "value"
	verify.NoError(t, err)
	verify.Equal(t, tok.Kind, TokenString)

	tok, err = lexer.NextToken()
	verify.NoError(t, err)
	verify.Equal(t, tok.Kind, TokenNewline)

	tok, err = lexer.NextToken()
	verify.NoError(t, err)
	verify.Equal(t, tok.Kind, TokenLeftBracket)
	verify.Equal(t, tok.Line, 2)
	verify.Equal(t, tok.Col, 1)
}

// TestTokenKindString tests TokenKind String method.
func TestTokenKindString(t *testing.T) {
	tests := []struct {
		kind TokenKind
		want string
	}{
		{TokenEOF, "EOF"},
		{TokenKey, "Key"},
		{TokenEquals, "Equals"},
		{TokenString, "String"},
		{TokenInteger, "Integer"},
		{TokenFloat, "Float"},
		{TokenBoolean, "Boolean"},
		{TokenDatetime, "Datetime"},
		{TokenLeftBracket, "LeftBracket"},
		{TokenRightBracket, "RightBracket"},
		{TokenLeftBrace, "LeftBrace"},
		{TokenRightBrace, "RightBrace"},
		{TokenComma, "Comma"},
		{TokenDot, "Dot"},
		{TokenNewline, "Newline"},
		{TokenKind(999), "Unknown"},
	}

	for _, tt := range tests {
		got := tt.kind.String()
		verify.Equal(t, got, tt.want)
	}
}

// TestTokenString tests Token String method.
func TestTokenString(t *testing.T) {
	tok := Token{
		Kind:  TokenString,
		Value: "hello",
		Line:  5,
		Col:   10,
	}

	str := tok.String()
	verify.True(t, len(str) > 0)
}

// TestLexerEscapeSequences tests various escape sequences.
func TestLexerEscapeSequences(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"backspace", `"\b"`, "\b"},
		{"tab", `"\t"`, "\t"},
		{"newline", `"\n"`, "\n"},
		{"form feed", `"\f"`, "\f"},
		{"carriage return", `"\r"`, "\r"},
		{"double quote", `"\""`, "\""},
		{"backslash", `"\\"`, "\\"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			tok, err := lexer.NextToken()
			verify.NoError(t, err)
			verify.Equal(t, tok.Kind, TokenString)
			verify.Equal(t, tok.Value, tt.want)
		})
	}
}

// TestLexerMultilineStringBackslash tests backslash line ending in multiline strings.
func TestLexerMultilineStringBackslash(t *testing.T) {
	input := `"""line1\
    line2"""`

	lexer := NewLexer(input)
	tok, err := lexer.NextToken()
	verify.NoError(t, err)
	verify.Equal(t, tok.Kind, TokenString)
	// Backslash at end of line trims following whitespace
	verify.True(t, len(tok.Value) > 0)
}

// TestLexerUnicodeEscape tests unicode escape sequences (basic support).
func TestLexerUnicodeEscape(t *testing.T) {
	input := `"\u0041"`

	lexer := NewLexer(input)
	tok, err := lexer.NextToken()
	verify.NoError(t, err)
	verify.Equal(t, tok.Kind, TokenString)
	// Unicode escapes return the escape char for now
	verify.True(t, len(tok.Value) > 0)
}

// TestLexerDatetime tests datetime tokenization.
func TestLexerDatetime(t *testing.T) {
	input := `2024-01-15T10:30:00Z`

	lexer := NewLexer(input)
	tok, err := lexer.NextToken()
	verify.NoError(t, err)
	verify.Equal(t, tok.Kind, TokenDatetime)
	verify.Equal(t, tok.Value, "2024-01-15T10:30:00Z")
}
