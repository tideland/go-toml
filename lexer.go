package toml

// Copyright (C) 2024 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

import (
	"fmt"
	"strings"
	"unicode"
)

// TokenKind represents the type of a token.
type TokenKind int

const (
	// TokenEOF indicates end of file.
	TokenEOF TokenKind = iota
	// TokenKey represents a key (bare, quoted, or dotted).
	TokenKey
	// TokenEquals represents the '=' character.
	TokenEquals
	// TokenString represents a string value.
	TokenString
	// TokenInteger represents an integer value.
	TokenInteger
	// TokenFloat represents a float value.
	TokenFloat
	// TokenBoolean represents a boolean value.
	TokenBoolean
	// TokenDatetime represents a datetime value.
	TokenDatetime
	// TokenLeftBracket represents '['.
	TokenLeftBracket
	// TokenRightBracket represents ']'.
	TokenRightBracket
	// TokenLeftBrace represents '{'.
	TokenLeftBrace
	// TokenRightBrace represents '}'.
	TokenRightBrace
	// TokenComma represents ','.
	TokenComma
	// TokenDot represents '.'.
	TokenDot
	// TokenNewline represents a newline.
	TokenNewline
)

// String returns the string representation of the token kind.
func (tk TokenKind) String() string {
	switch tk {
	case TokenEOF:
		return "EOF"
	case TokenKey:
		return "Key"
	case TokenEquals:
		return "Equals"
	case TokenString:
		return "String"
	case TokenInteger:
		return "Integer"
	case TokenFloat:
		return "Float"
	case TokenBoolean:
		return "Boolean"
	case TokenDatetime:
		return "Datetime"
	case TokenLeftBracket:
		return "LeftBracket"
	case TokenRightBracket:
		return "RightBracket"
	case TokenLeftBrace:
		return "LeftBrace"
	case TokenRightBrace:
		return "RightBrace"
	case TokenComma:
		return "Comma"
	case TokenDot:
		return "Dot"
	case TokenNewline:
		return "Newline"
	default:
		return "Unknown"
	}
}

// Token represents a lexical token in TOML.
type Token struct {
	Kind  TokenKind
	Value string
	Line  int
	Col   int
}

// String returns the string representation of the token.
func (t Token) String() string {
	return fmt.Sprintf("%s(%q) at %d:%d", t.Kind, t.Value, t.Line, t.Col)
}

// Lexer tokenizes TOML input.
type Lexer struct {
	input []rune
	pos   int
	line  int
	col   int
}

// NewLexer creates a new lexer for the given input.
func NewLexer(input string) *Lexer {
	return &Lexer{
		input: []rune(input),
		pos:   0,
		line:  1,
		col:   1,
	}
}

// NextToken returns the next token from the input.
func (l *Lexer) NextToken() (Token, error) {
	l.skipWhitespace()

	if l.isEOF() {
		return Token{Kind: TokenEOF, Line: l.line, Col: l.col}, nil
	}

	// Skip comments
	if l.peek() == '#' {
		l.skipComment()
		return l.NextToken()
	}

	// Save position for token
	line, col := l.line, l.col

	ch := l.peek()

	// Single character tokens
	switch ch {
	case '\n':
		l.advance()
		return Token{Kind: TokenNewline, Value: "\n", Line: line, Col: col}, nil
	case '=':
		l.advance()
		return Token{Kind: TokenEquals, Value: "=", Line: line, Col: col}, nil
	case '[':
		l.advance()
		return Token{Kind: TokenLeftBracket, Value: "[", Line: line, Col: col}, nil
	case ']':
		l.advance()
		return Token{Kind: TokenRightBracket, Value: "]", Line: line, Col: col}, nil
	case '{':
		l.advance()
		return Token{Kind: TokenLeftBrace, Value: "{", Line: line, Col: col}, nil
	case '}':
		l.advance()
		return Token{Kind: TokenRightBrace, Value: "}", Line: line, Col: col}, nil
	case ',':
		l.advance()
		return Token{Kind: TokenComma, Value: ",", Line: line, Col: col}, nil
	case '.':
		l.advance()
		return Token{Kind: TokenDot, Value: ".", Line: line, Col: col}, nil
	}

	// String values
	if ch == '"' || ch == '\'' {
		return l.scanString(line, col)
	}

	// Numbers, booleans, datetimes, or keys
	if unicode.IsLetter(ch) || ch == '_' {
		return l.scanKeyOrKeyword(line, col)
	}

	// Special handling for +/- with inf/nan
	if (ch == '+' || ch == '-') && !l.isEOF() {
		next := l.peekAhead(1)
		if next == 'i' || next == 'n' {
			// Could be +inf, -inf, +nan, -nan
			return l.scanKeyOrKeyword(line, col)
		}
		return l.scanNumber(line, col)
	}

	if unicode.IsDigit(ch) {
		return l.scanNumber(line, col)
	}

	return Token{}, NewErrorWithLocation("lexer", "", line, col,
		fmt.Errorf("unexpected character: %q", ch), ErrSyntax)
}

// peek returns the current character without advancing.
func (l *Lexer) peek() rune {
	if l.isEOF() {
		return 0
	}
	return l.input[l.pos]
}

// peekAhead looks ahead n characters without advancing.
func (l *Lexer) peekAhead(n int) rune {
	pos := l.pos + n
	if pos >= len(l.input) {
		return 0
	}
	return l.input[pos]
}

// advance moves to the next character.
func (l *Lexer) advance() rune {
	if l.isEOF() {
		return 0
	}
	ch := l.input[l.pos]
	l.pos++
	if ch == '\n' {
		l.line++
		l.col = 1
	} else {
		l.col++
	}
	return ch
}

// isEOF checks if we've reached the end of input.
func (l *Lexer) isEOF() bool {
	return l.pos >= len(l.input)
}

// skipWhitespace skips whitespace characters (except newlines).
func (l *Lexer) skipWhitespace() {
	for !l.isEOF() {
		ch := l.peek()
		if ch == ' ' || ch == '\t' || ch == '\r' {
			l.advance()
		} else {
			break
		}
	}
}

// skipComment skips from # to end of line.
func (l *Lexer) skipComment() {
	for !l.isEOF() && l.peek() != '\n' {
		l.advance()
	}
}

// scanString scans a string value (basic, literal, or multi-line).
func (l *Lexer) scanString(line, col int) (Token, error) {
	quote := l.advance() // Consume opening quote

	// Check for multi-line string (""" or ''')
	if l.peek() == quote && l.peekAhead(1) == quote {
		l.advance() // Second quote
		l.advance() // Third quote
		return l.scanMultilineString(quote, line, col)
	}

	var value strings.Builder
	for !l.isEOF() {
		ch := l.peek()

		if ch == quote {
			l.advance() // Consume closing quote
			return Token{Kind: TokenString, Value: value.String(), Line: line, Col: col}, nil
		}

		if ch == '\n' {
			return Token{}, NewErrorWithLocation("lexer", "", line, col,
				fmt.Errorf("unclosed string"), ErrSyntax)
		}

		if ch == '\\' && quote == '"' {
			// Handle escape sequences in basic strings
			l.advance() // Consume backslash
			escaped, err := l.scanEscape(line, col)
			if err != nil {
				return Token{}, err
			}
			value.WriteString(escaped)
		} else {
			value.WriteRune(ch)
			l.advance()
		}
	}

	return Token{}, NewErrorWithLocation("lexer", "", line, col,
		fmt.Errorf("unclosed string"), ErrSyntax)
}

// scanMultilineString scans a multi-line string.
func (l *Lexer) scanMultilineString(quote rune, line, col int) (Token, error) {
	var value strings.Builder

	// Skip first newline if present
	if l.peek() == '\n' {
		l.advance()
	}

	for !l.isEOF() {
		ch := l.peek()

		// Check for closing """  or '''
		if ch == quote && l.peekAhead(1) == quote && l.peekAhead(2) == quote {
			l.advance() // First quote
			l.advance() // Second quote
			l.advance() // Third quote
			return Token{Kind: TokenString, Value: value.String(), Line: line, Col: col}, nil
		}

		if ch == '\\' && quote == '"' {
			// Handle escape sequences in multi-line basic strings
			l.advance() // Consume backslash

			// Line ending backslash trims whitespace
			if l.peek() == '\n' {
				l.advance() // Consume newline
				// Skip following whitespace
				for !l.isEOF() && (l.peek() == ' ' || l.peek() == '\t' || l.peek() == '\n' || l.peek() == '\r') {
					l.advance()
				}
				continue
			}

			escaped, err := l.scanEscape(line, col)
			if err != nil {
				return Token{}, err
			}
			value.WriteString(escaped)
		} else {
			value.WriteRune(ch)
			l.advance()
		}
	}

	return Token{}, NewErrorWithLocation("lexer", "", line, col,
		fmt.Errorf("unclosed multi-line string"), ErrSyntax)
}

// scanEscape scans an escape sequence.
func (l *Lexer) scanEscape(line, col int) (string, error) {
	if l.isEOF() {
		return "", NewErrorWithLocation("lexer", "", line, col,
			fmt.Errorf("incomplete escape sequence"), ErrSyntax)
	}

	ch := l.advance()
	switch ch {
	case 'b':
		return "\b", nil
	case 't':
		return "\t", nil
	case 'n':
		return "\n", nil
	case 'f':
		return "\f", nil
	case 'r':
		return "\r", nil
	case '"':
		return "\"", nil
	case '\\':
		return "\\", nil
	case 'u', 'U':
		// Unicode escape sequences (not fully implemented yet)
		return string(ch), nil
	default:
		return "", NewErrorWithLocation("lexer", "", line, col,
			fmt.Errorf("invalid escape sequence: \\%c", ch), ErrSyntax)
	}
}

// scanKeyOrKeyword scans a bare key or keyword (true, false, inf, nan).
func (l *Lexer) scanKeyOrKeyword(line, col int) (Token, error) {
	var value strings.Builder

	for !l.isEOF() {
		ch := l.peek()
		if unicode.IsLetter(ch) || unicode.IsDigit(ch) || ch == '_' || ch == '-' {
			value.WriteRune(ch)
			l.advance()
		} else {
			break
		}
	}

	str := value.String()

	// Check for keywords
	switch str {
	case "true", "false":
		return Token{Kind: TokenBoolean, Value: str, Line: line, Col: col}, nil
	case "inf", "nan", "+inf", "-inf", "+nan", "-nan":
		return Token{Kind: TokenFloat, Value: str, Line: line, Col: col}, nil
	default:
		return Token{Kind: TokenKey, Value: str, Line: line, Col: col}, nil
	}
}

// scanNumber scans a number (integer, float, or datetime).
func (l *Lexer) scanNumber(line, col int) (Token, error) {
	var value strings.Builder
	hasDecimal := false
	hasExponent := false

	// Handle sign
	if l.peek() == '+' || l.peek() == '-' {
		value.WriteRune(l.advance())
	}

	// Check for special prefixes (0x, 0o, 0b for hex, octal, binary)
	if l.peek() == '0' && !l.isEOF() {
		next := l.peekAhead(1)
		if next == 'x' || next == 'o' || next == 'b' {
			value.WriteRune(l.advance()) // 0
			value.WriteRune(l.advance()) // x, o, or b
			return l.scanSpecialInteger(value.String(), line, col)
		}
	}

	for !l.isEOF() {
		ch := l.peek()

		if unicode.IsDigit(ch) {
			value.WriteRune(ch)
			l.advance()
		} else if ch == '_' {
			// Underscores are allowed in numbers
			l.advance()
		} else if ch == '.' && !hasDecimal && !hasExponent {
			value.WriteRune(ch)
			l.advance()
			hasDecimal = true
		} else if (ch == 'e' || ch == 'E') && !hasExponent {
			value.WriteRune(ch)
			l.advance()
			hasExponent = true
			// Handle exponent sign
			if l.peek() == '+' || l.peek() == '-' {
				value.WriteRune(l.advance())
			}
		} else if ch == '-' || ch == ':' || ch == 'T' || ch == 'Z' {
			// Could be a datetime
			return l.scanDatetime(value.String(), line, col)
		} else {
			break
		}
	}

	str := value.String()

	if hasDecimal || hasExponent {
		return Token{Kind: TokenFloat, Value: str, Line: line, Col: col}, nil
	}
	return Token{Kind: TokenInteger, Value: str, Line: line, Col: col}, nil
}

// scanSpecialInteger scans hex, octal, or binary integers.
func (l *Lexer) scanSpecialInteger(prefix string, line, col int) (Token, error) {
	var value strings.Builder
	value.WriteString(prefix)

	for !l.isEOF() {
		ch := l.peek()
		if unicode.IsDigit(ch) || (ch >= 'a' && ch <= 'f') || (ch >= 'A' && ch <= 'F') {
			value.WriteRune(ch)
			l.advance()
		} else if ch == '_' {
			// Underscores are allowed
			l.advance()
		} else {
			break
		}
	}

	return Token{Kind: TokenInteger, Value: value.String(), Line: line, Col: col}, nil
}

// scanDatetime continues scanning a datetime value.
func (l *Lexer) scanDatetime(prefix string, line, col int) (Token, error) {
	var value strings.Builder
	value.WriteString(prefix)

	for !l.isEOF() {
		ch := l.peek()
		if unicode.IsDigit(ch) || ch == '-' || ch == ':' || ch == 'T' || ch == 'Z' || ch == '.' || ch == '+' {
			value.WriteRune(ch)
			l.advance()
		} else {
			break
		}
	}

	return Token{Kind: TokenDatetime, Value: value.String(), Line: line, Col: col}, nil
}
