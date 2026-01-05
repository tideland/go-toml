package toml

// Copyright (C) 2024 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

import (
	"testing"
	"time"

	"tideland.dev/go/asserts/verify"
)

// TestSectionHas tests the Section Has method.
func TestSectionHas(t *testing.T) {
	input := `[section]
key = "value"
number = 42`

	doc, err := ParseString(input, "test.toml")
	verify.NoError(t, err)

	section, err := doc.Section("section")
	verify.NoError(t, err)

	verify.True(t, section.Has("key"))
	verify.True(t, section.Has("number"))
	verify.False(t, section.Has("missing"))
}

// TestSectionGetFloat tests Section GetFloat method.
func TestSectionGetFloat(t *testing.T) {
	input := `[section]
pi = 3.14159`

	doc, err := ParseString(input, "test.toml")
	verify.NoError(t, err)

	section, err := doc.Section("section")
	verify.NoError(t, err)

	pi, err := section.GetFloat("pi")
	verify.NoError(t, err)
	verify.True(t, pi > 3.14 && pi < 3.15)

	// Test missing key
	_, err = section.GetFloat("missing")
	verify.Error(t, err)
}

// TestSectionGetBool tests Section GetBool method.
func TestSectionGetBool(t *testing.T) {
	input := `[section]
enabled = true
disabled = false`

	doc, err := ParseString(input, "test.toml")
	verify.NoError(t, err)

	section, err := doc.Section("section")
	verify.NoError(t, err)

	enabled, err := section.GetBool("enabled")
	verify.NoError(t, err)
	verify.True(t, enabled)

	disabled, err := section.GetBool("disabled")
	verify.NoError(t, err)
	verify.False(t, disabled)

	// Test missing key
	_, err = section.GetBool("missing")
	verify.Error(t, err)
}

// TestSectionGetTime tests Section GetTime method.
func TestSectionGetTime(t *testing.T) {
	input := `[section]
timestamp = 1979-05-27T07:32:00Z`

	doc, err := ParseString(input, "test.toml")
	verify.NoError(t, err)

	section, err := doc.Section("section")
	verify.NoError(t, err)

	ts, err := section.GetTime("timestamp")
	verify.NoError(t, err)
	verify.Equal(t, ts.Year(), 1979)

	// Test missing key
	_, err = section.GetTime("missing")
	verify.Error(t, err)
}

// TestSectionGetArray tests Section GetArray method.
func TestSectionGetArray(t *testing.T) {
	input := `[section]
numbers = [1, 2, 3]`

	doc, err := ParseString(input, "test.toml")
	verify.NoError(t, err)

	section, err := doc.Section("section")
	verify.NoError(t, err)

	arr, err := section.GetArray("numbers")
	verify.NoError(t, err)
	verify.Equal(t, len(arr), 3)

	// Test missing key
	_, err = section.GetArray("missing")
	verify.Error(t, err)
}

// TestDocumentPath tests Document Path method.
func TestDocumentPath(t *testing.T) {
	doc, err := Parse("testdata/valid/basic.toml")
	verify.NoError(t, err)

	path := doc.Path()
	verify.Equal(t, path, "testdata/valid/basic.toml")
}

// TestDocumentGetTime tests Document GetTime method.
func TestDocumentGetTime(t *testing.T) {
	input := `timestamp = 2024-01-15T10:30:00Z

[section]
date = 2024-12-25T00:00:00Z`

	doc, err := ParseString(input, "test.toml")
	verify.NoError(t, err)

	ts, err := doc.GetTime("timestamp")
	verify.NoError(t, err)
	verify.Equal(t, ts.Year(), 2024)

	// Nested timestamp
	sectionDate, err := doc.GetTime("section.date")
	verify.NoError(t, err)
	verify.Equal(t, sectionDate.Month(), time.December)
}

// TestSectionNotFound tests error when section not found.
func TestSectionNotFound(t *testing.T) {
	input := `[existing]
key = "value"`

	doc, err := ParseString(input, "test.toml")
	verify.NoError(t, err)

	section, err := doc.Section("existing")
	verify.NoError(t, err)
	verify.NotNil(t, section)

	_, err = section.Section("missing")
	verify.Error(t, err)
}

// TestTableArrayNotFound tests error when table array not found.
func TestTableArrayNotFound(t *testing.T) {
	input := `[[existing]]
key = "value"`

	doc, err := ParseString(input, "test.toml")
	verify.NoError(t, err)

	_, err = doc.Section("existing")
	verify.Error(t, err) // It's a table array, not a section

	_, err = doc.TableArray("missing")
	verify.Error(t, err)
}

// TestValueTimeError tests Time conversion error.
func TestValueTimeError(t *testing.T) {
	v := NewStringValue("not a time")
	_, err := v.Time()
	verify.Error(t, err)
}

// TestValueTableError tests Table conversion error.
func TestValueTableError(t *testing.T) {
	v := NewStringValue("not a table")
	_, err := v.Table()
	verify.Error(t, err)
}

// TestResolvePathErrors tests various resolvePath error cases.
func TestResolvePathErrors(t *testing.T) {
	input := `[section]
key = "value"`

	doc, err := ParseString(input, "test.toml")
	verify.NoError(t, err)

	// Empty path
	_, err = doc.GetString("")
	verify.Error(t, err)

	// Section doesn't exist
	_, err = doc.GetString("missing.key")
	verify.Error(t, err)

	// Key doesn't exist
	_, err = doc.GetString("section.missing")
	verify.Error(t, err)
}
