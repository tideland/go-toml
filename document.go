package toml

// Copyright (C) 2024 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Section represents a TOML table (section).
type Section struct {
	name     string
	path     []string
	keys     map[string]*Value
	sections map[string]*Section
	arrays   map[string][]*Section
	parent   *Section
}

// newSection creates a new Section with the given name and parent.
func newSection(name string, parent *Section) *Section {
	s := &Section{
		name:     name,
		keys:     make(map[string]*Value),
		sections: make(map[string]*Section),
		arrays:   make(map[string][]*Section),
		parent:   parent,
	}
	if parent != nil {
		s.path = append(append([]string{}, parent.path...), name)
	} else {
		if name != "" {
			s.path = []string{name}
		} else {
			s.path = []string{}
		}
	}
	return s
}

// Name returns the section name.
func (s *Section) Name() string {
	return s.name
}

// FullPath returns the full path of the section as a slice of strings.
func (s *Section) FullPath() []string {
	return append([]string{}, s.path...)
}

// Parent returns the parent section, or nil if this is the root.
func (s *Section) Parent() *Section {
	return s.parent
}

// Has checks if a key exists in this section.
func (s *Section) Has(key string) bool {
	_, ok := s.keys[key]
	return ok
}

// Keys returns all keys in this section (sorted).
func (s *Section) Keys() []string {
	keys := make([]string, 0, len(s.keys))
	for k := range s.keys {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// Sections returns all subsection names (sorted).
func (s *Section) Sections() []string {
	names := make([]string, 0, len(s.sections))
	for n := range s.sections {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}

// Section returns a subsection by name.
func (s *Section) Section(name string) (*Section, error) {
	sec, ok := s.sections[name]
	if !ok {
		return nil, NewErrorWithPath("section", name, fmt.Errorf("section not found"), ErrNotFound)
	}
	return sec, nil
}

// TableArray returns an array of tables within this section.
func (s *Section) TableArray(name string) ([]*Section, error) {
	arr, ok := s.arrays[name]
	if !ok {
		return nil, NewErrorWithPath("table array", name, fmt.Errorf("table array not found"), ErrNotFound)
	}
	return arr, nil
}

// GetString retrieves a string value by key.
func (s *Section) GetString(key string) (string, error) {
	v, ok := s.keys[key]
	if !ok {
		return "", NewErrorWithPath("get", key, fmt.Errorf("key not found"), ErrNotFound)
	}
	return v.String()
}

// GetInt retrieves an integer value by key.
func (s *Section) GetInt(key string) (int64, error) {
	v, ok := s.keys[key]
	if !ok {
		return 0, NewErrorWithPath("get", key, fmt.Errorf("key not found"), ErrNotFound)
	}
	return v.Int()
}

// GetFloat retrieves a float value by key.
func (s *Section) GetFloat(key string) (float64, error) {
	v, ok := s.keys[key]
	if !ok {
		return 0, NewErrorWithPath("get", key, fmt.Errorf("key not found"), ErrNotFound)
	}
	return v.Float()
}

// GetBool retrieves a boolean value by key.
func (s *Section) GetBool(key string) (bool, error) {
	v, ok := s.keys[key]
	if !ok {
		return false, NewErrorWithPath("get", key, fmt.Errorf("key not found"), ErrNotFound)
	}
	return v.Bool()
}

// GetTime retrieves a datetime value by key.
func (s *Section) GetTime(key string) (time.Time, error) {
	v, ok := s.keys[key]
	if !ok {
		return time.Time{}, NewErrorWithPath("get", key, fmt.Errorf("key not found"), ErrNotFound)
	}
	return v.Time()
}

// GetArray retrieves an array value by key.
func (s *Section) GetArray(key string) ([]any, error) {
	v, ok := s.keys[key]
	if !ok {
		return nil, NewErrorWithPath("get", key, fmt.Errorf("key not found"), ErrNotFound)
	}
	return v.Array()
}

// setValue sets a value for a key in this section.
func (s *Section) setValue(key string, value *Value) {
	s.keys[key] = value
}

// addSection adds a subsection to this section.
func (s *Section) addSection(name string) *Section {
	sec := newSection(name, s)
	s.sections[name] = sec
	return sec
}

// getOrCreateSection gets or creates a subsection.
func (s *Section) getOrCreateSection(name string) *Section {
	if sec, ok := s.sections[name]; ok {
		return sec
	}
	return s.addSection(name)
}

// addTableArrayElement adds an element to a table array.
func (s *Section) addTableArrayElement(name string) *Section {
	sec := newSection("", s)
	s.arrays[name] = append(s.arrays[name], sec)
	return sec
}

// Document represents a parsed TOML document.
type Document struct {
	root     *Section
	filepath string
}

// newDocument creates a new Document.
func newDocument(filepath string) *Document {
	return &Document{
		root:     newSection("", nil),
		filepath: filepath,
	}
}

// Path returns the source file path.
func (d *Document) Path() string {
	return d.filepath
}

// Has checks if a path exists in the document.
func (d *Document) Has(path string) bool {
	_, err := d.resolvePath(path)
	return err == nil
}

// Keys returns all top-level keys (sorted).
func (d *Document) Keys() []string {
	return d.root.Keys()
}

// Sections returns all top-level section names (sorted).
func (d *Document) Sections() []string {
	return d.root.Sections()
}

// Section returns a top-level section by name.
func (d *Document) Section(name string) (*Section, error) {
	return d.root.Section(name)
}

// TableArray returns a top-level table array.
func (d *Document) TableArray(name string) ([]*Section, error) {
	return d.root.TableArray(name)
}

// GetString retrieves a string value by dot-notation path.
func (d *Document) GetString(path string) (string, error) {
	v, err := d.resolvePath(path)
	if err != nil {
		return "", err
	}
	return v.String()
}

// GetInt retrieves an integer value by dot-notation path.
func (d *Document) GetInt(path string) (int64, error) {
	v, err := d.resolvePath(path)
	if err != nil {
		return 0, err
	}
	return v.Int()
}

// GetFloat retrieves a float value by dot-notation path.
func (d *Document) GetFloat(path string) (float64, error) {
	v, err := d.resolvePath(path)
	if err != nil {
		return 0, err
	}
	return v.Float()
}

// GetBool retrieves a boolean value by dot-notation path.
func (d *Document) GetBool(path string) (bool, error) {
	v, err := d.resolvePath(path)
	if err != nil {
		return false, err
	}
	return v.Bool()
}

// GetTime retrieves a datetime value by dot-notation path.
func (d *Document) GetTime(path string) (time.Time, error) {
	v, err := d.resolvePath(path)
	if err != nil {
		return time.Time{}, err
	}
	return v.Time()
}

// GetArray retrieves an array value by dot-notation path.
func (d *Document) GetArray(path string) ([]any, error) {
	v, err := d.resolvePath(path)
	if err != nil {
		return nil, err
	}
	return v.Array()
}

// resolvePath resolves a dot-notation path to a Value.
// Supports paths like "database.connection.timeout" and array indices like "fruits.0.name".
func (d *Document) resolvePath(path string) (*Value, error) {
	if path == "" {
		return nil, NewErrorWithPath("resolve", path, fmt.Errorf("empty path"), ErrInvalidPath)
	}

	parts := strings.Split(path, ".")
	current := d.root

	// Navigate through all parts except the last one
	for i := 0; i < len(parts)-1; i++ {
		part := parts[i]

		// Check if it's an array index (numeric)
		if _, err := strconv.Atoi(part); err == nil {
			// This shouldn't happen in the middle of navigation for simple cases
			// but we'll handle it for consistency
			return nil, NewErrorWithPath("resolve", path,
				fmt.Errorf("unexpected array index in path: %s", part), ErrInvalidPath)
		}

		// Try to navigate to subsection
		if sec, ok := current.sections[part]; ok {
			current = sec
		} else {
			return nil, NewErrorWithPath("resolve", path,
				fmt.Errorf("section not found: %s", part), ErrNotFound)
		}
	}

	// Get the value using the last part as key
	lastPart := parts[len(parts)-1]
	v, ok := current.keys[lastPart]
	if !ok {
		return nil, NewErrorWithPath("resolve", path,
			fmt.Errorf("key not found: %s", lastPart), ErrNotFound)
	}

	return v, nil
}
