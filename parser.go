package toml

// Copyright (C) 2024 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

// Parser parses TOML input into a Document.
type Parser struct {
	lexer   *Lexer
	current Token
	doc     *Document
	section *Section
}

// newParser creates a new parser for the given input.
func newParser(input string, filepath string) *Parser {
	p := &Parser{
		lexer: NewLexer(input),
		doc:   newDocument(filepath),
	}
	p.section = p.doc.root
	return p
}

// Parse parses a TOML file from disk.
func Parse(filepath string) (*Document, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, NewErrorWithPath("parse", filepath, err, ErrIO)
	}
	return ParseString(string(data), filepath)
}

// ParseString parses a TOML string.
func ParseString(input string, filepath string) (*Document, error) {
	p := newParser(input, filepath)
	if err := p.parse(); err != nil {
		return nil, err
	}
	return p.doc, nil
}

// parse is the main parsing loop.
func (p *Parser) parse() error {
	// Advance to first token
	if err := p.next(); err != nil {
		return err
	}

	for p.current.Kind != TokenEOF {
		// Skip newlines at top level
		if p.current.Kind == TokenNewline {
			if err := p.next(); err != nil {
				return err
			}
			continue
		}

		// Handle table headers [section] or [[array]]
		if p.current.Kind == TokenLeftBracket {
			if err := p.parseTable(); err != nil {
				return err
			}
			continue
		}

		// Handle key-value pairs
		if p.current.Kind == TokenKey || p.current.Kind == TokenString {
			if err := p.parseKeyValue(); err != nil {
				return err
			}
			continue
		}

		return p.error(fmt.Errorf("unexpected token: %s", p.current.Kind))
	}

	return nil
}

// next advances to the next token.
func (p *Parser) next() error {
	tok, err := p.lexer.NextToken()
	if err != nil {
		return err
	}
	p.current = tok
	return nil
}

// expect checks if the current token matches the expected kind and advances.
func (p *Parser) expect(kind TokenKind) error {
	if p.current.Kind != kind {
		return p.error(fmt.Errorf("expected %s, got %s", kind, p.current.Kind))
	}
	return p.next()
}

// error creates a parse error with current position.
func (p *Parser) error(err error) error {
	return NewErrorWithLocation("parse", p.doc.filepath,
		p.current.Line, p.current.Col, err, ErrSyntax)
}

// parseTable parses a table header: [section] or [[array]].
func (p *Parser) parseTable() error {
	if err := p.expect(TokenLeftBracket); err != nil {
		return err
	}

	// Check for array of tables [[...]]
	isArray := false
	if p.current.Kind == TokenLeftBracket {
		isArray = true
		if err := p.next(); err != nil {
			return err
		}
	}

	// Parse table path
	path, err := p.parseTablePath()
	if err != nil {
		return err
	}

	// Expect closing bracket(s)
	if err := p.expect(TokenRightBracket); err != nil {
		return err
	}
	if isArray {
		if err := p.expect(TokenRightBracket); err != nil {
			return err
		}
	}

	// Skip optional newline
	if p.current.Kind == TokenNewline {
		if err := p.next(); err != nil {
			return err
		}
	}

	// Navigate to or create the section
	if isArray {
		p.section = p.navigateToTableArray(path)
	} else {
		sec, err := p.navigateToTable(path)
		if err != nil {
			return err
		}
		p.section = sec
	}

	return nil
}

// parseTablePath parses a dotted table path like "a.b.c".
func (p *Parser) parseTablePath() ([]string, error) {
	var path []string

	for {
		if p.current.Kind != TokenKey && p.current.Kind != TokenString {
			return nil, p.error(fmt.Errorf("expected key in table path"))
		}

		path = append(path, p.current.Value)
		if err := p.next(); err != nil {
			return nil, err
		}

		if p.current.Kind == TokenDot {
			if err := p.next(); err != nil {
				return nil, err
			}
			continue
		}

		break
	}

	return path, nil
}

// navigateToTable navigates to a table, creating it if necessary.
func (p *Parser) navigateToTable(path []string) (*Section, error) {
	current := p.doc.root

	for _, name := range path {
		current = current.getOrCreateSection(name)
	}

	return current, nil
}

// navigateToTableArray navigates to a table array, creating a new element.
func (p *Parser) navigateToTableArray(path []string) *Section {
	current := p.doc.root

	// Navigate to parent
	for i := 0; i < len(path)-1; i++ {
		current = current.getOrCreateSection(path[i])
	}

	// Add array element
	name := path[len(path)-1]
	return current.addTableArrayElement(name)
}

// parseKeyValue parses a key-value pair.
func (p *Parser) parseKeyValue() error {
	// Parse key (might be dotted)
	keys, err := p.parseKey()
	if err != nil {
		return err
	}

	// Expect equals
	if err := p.expect(TokenEquals); err != nil {
		return err
	}

	// Parse value
	value, err := p.parseValue()
	if err != nil {
		return err
	}

	// Navigate to the right section for dotted keys
	section := p.section
	for i := 0; i < len(keys)-1; i++ {
		section = section.getOrCreateSection(keys[i])
	}

	// Set the value
	finalKey := keys[len(keys)-1]
	section.setValue(finalKey, value)

	// Skip optional newline
	if p.current.Kind == TokenNewline {
		if err := p.next(); err != nil {
			return err
		}
	}

	return nil
}

// parseKey parses a key, which might be dotted like "a.b.c".
func (p *Parser) parseKey() ([]string, error) {
	var keys []string

	for {
		if p.current.Kind != TokenKey && p.current.Kind != TokenString {
			return nil, p.error(fmt.Errorf("expected key"))
		}

		keys = append(keys, p.current.Value)
		if err := p.next(); err != nil {
			return nil, err
		}

		if p.current.Kind == TokenDot {
			if err := p.next(); err != nil {
				return nil, err
			}
			continue
		}

		break
	}

	return keys, nil
}

// parseValue parses a value (string, integer, float, boolean, datetime, array, inline table).
func (p *Parser) parseValue() (*Value, error) {
	switch p.current.Kind {
	case TokenString:
		v := NewStringValue(p.current.Value)
		if err := p.next(); err != nil {
			return nil, err
		}
		return v, nil

	case TokenInteger:
		i, err := p.parseInteger(p.current.Value)
		if err != nil {
			return nil, p.error(err)
		}
		if err := p.next(); err != nil {
			return nil, err
		}
		return NewIntegerValue(i), nil

	case TokenFloat:
		f, err := p.parseFloat(p.current.Value)
		if err != nil {
			return nil, p.error(err)
		}
		if err := p.next(); err != nil {
			return nil, err
		}
		return NewFloatValue(f), nil

	case TokenBoolean:
		b := p.current.Value == "true"
		if err := p.next(); err != nil {
			return nil, err
		}
		return NewBooleanValue(b), nil

	case TokenDatetime:
		t, err := p.parseDatetime(p.current.Value)
		if err != nil {
			return nil, p.error(err)
		}
		if err := p.next(); err != nil {
			return nil, err
		}
		return NewDatetimeValue(t), nil

	case TokenLeftBracket:
		return p.parseArray()

	case TokenLeftBrace:
		return p.parseInlineTable()

	default:
		return nil, p.error(fmt.Errorf("unexpected token for value: %s", p.current.Kind))
	}
}

// parseInteger parses an integer value (decimal, hex, octal, binary).
func (p *Parser) parseInteger(s string) (int64, error) {
	// Remove underscores
	s = strings.ReplaceAll(s, "_", "")

	// Handle different bases
	if strings.HasPrefix(s, "0x") {
		return strconv.ParseInt(s[2:], 16, 64)
	}
	if strings.HasPrefix(s, "0o") {
		return strconv.ParseInt(s[2:], 8, 64)
	}
	if strings.HasPrefix(s, "0b") {
		return strconv.ParseInt(s[2:], 2, 64)
	}

	return strconv.ParseInt(s, 10, 64)
}

// parseFloat parses a float value.
func (p *Parser) parseFloat(s string) (float64, error) {
	// Handle special values
	switch s {
	case "inf", "+inf":
		return math.Inf(1), nil
	case "-inf":
		return math.Inf(-1), nil
	case "nan", "+nan", "-nan":
		return math.NaN(), nil
	}

	// Remove underscores
	s = strings.ReplaceAll(s, "_", "")

	return strconv.ParseFloat(s, 64)
}

// parseDatetime parses a datetime value.
func (p *Parser) parseDatetime(s string) (time.Time, error) {
	// Try different datetime formats supported by TOML v1.0.0
	formats := []string{
		time.RFC3339,                    // Offset Date-Time
		time.RFC3339Nano,                // Offset Date-Time with fractional seconds
		"2006-01-02T15:04:05",           // Local Date-Time
		"2006-01-02T15:04:05.999999999", // Local Date-Time with fractional seconds
		"2006-01-02",                    // Local Date
		"15:04:05",                      // Local Time
		"15:04:05.999999999",            // Local Time with fractional seconds
	}

	for _, format := range formats {
		if t, err := time.Parse(format, s); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("invalid datetime format: %s", s)
}

// parseArray parses an array value.
func (p *Parser) parseArray() (*Value, error) {
	if err := p.expect(TokenLeftBracket); err != nil {
		return nil, err
	}

	var elements []any

	for p.current.Kind != TokenRightBracket {
		// Skip newlines in arrays
		if p.current.Kind == TokenNewline {
			if err := p.next(); err != nil {
				return nil, err
			}
			continue
		}

		// Parse array element
		val, err := p.parseValue()
		if err != nil {
			return nil, err
		}

		// Convert Value to its raw form for array storage
		elements = append(elements, val.raw)

		// Skip optional comma
		if p.current.Kind == TokenComma {
			if err := p.next(); err != nil {
				return nil, err
			}
		}

		// Skip newlines
		for p.current.Kind == TokenNewline {
			if err := p.next(); err != nil {
				return nil, err
			}
		}
	}

	if err := p.expect(TokenRightBracket); err != nil {
		return nil, err
	}

	return NewArrayValue(elements), nil
}

// parseInlineTable parses an inline table.
func (p *Parser) parseInlineTable() (*Value, error) {
	if err := p.expect(TokenLeftBrace); err != nil {
		return nil, err
	}

	table := make(map[string]*Value)

	for p.current.Kind != TokenRightBrace {
		// Parse key
		if p.current.Kind != TokenKey && p.current.Kind != TokenString {
			return nil, p.error(fmt.Errorf("expected key in inline table"))
		}

		key := p.current.Value
		if err := p.next(); err != nil {
			return nil, err
		}

		// Expect equals
		if err := p.expect(TokenEquals); err != nil {
			return nil, err
		}

		// Parse value
		val, err := p.parseValue()
		if err != nil {
			return nil, err
		}

		table[key] = val

		// Skip optional comma
		if p.current.Kind == TokenComma {
			if err := p.next(); err != nil {
				return nil, err
			}
		}
	}

	if err := p.expect(TokenRightBrace); err != nil {
		return nil, err
	}

	return NewTableValue(table), nil
}
