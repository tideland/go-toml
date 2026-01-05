package toml

// Copyright (C) 2024 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

import (
	"testing"

	"tideland.dev/go/asserts/verify"
)

// TestConfigDefaults tests default configuration values.
func TestConfigDefaults(t *testing.T) {
	cfg := NewConfig()

	verify.False(t, cfg.Strict())
	verify.False(t, cfg.AllowDuplicateKeys())
	verify.Equal(t, cfg.MaxDepth(), 100)
	verify.NoError(t, cfg.Validate())
}

// TestConfigSetters tests configuration setters.
func TestConfigSetters(t *testing.T) {
	cfg := NewConfig().
		SetStrict(true).
		SetAllowDuplicateKeys(true).
		SetMaxDepth(50)

	verify.True(t, cfg.Strict())
	verify.True(t, cfg.AllowDuplicateKeys())
	verify.Equal(t, cfg.MaxDepth(), 50)
	verify.NoError(t, cfg.Validate())
}

// TestConfigErrorAccumulation tests error accumulation.
func TestConfigErrorAccumulation(t *testing.T) {
	cfg := NewConfig().
		SetMaxDepth(0). // Invalid
		SetMaxDepth(-1) // Invalid

	err := cfg.Validate()
	verify.Error(t, err)
}

// TestConfigFluentAPI tests fluent API chaining.
func TestConfigFluentAPI(t *testing.T) {
	cfg := NewConfig()

	// Chaining should work
	result := cfg.SetStrict(true).SetAllowDuplicateKeys(false).SetMaxDepth(200)
	verify.Equal(t, result, cfg) // Same instance
	verify.True(t, cfg.Strict())
	verify.False(t, cfg.AllowDuplicateKeys())
	verify.Equal(t, cfg.MaxDepth(), 200)
}

// TestConfigError tests Error method (alias for Validate).
func TestConfigError(t *testing.T) {
	cfg := NewConfig()
	verify.NoError(t, cfg.Error())

	cfg.SetMaxDepth(0)
	verify.Error(t, cfg.Error())
}

// TestParseWithValidConfig tests parsing with valid configuration.
func TestParseWithValidConfig(t *testing.T) {
	cfg := NewConfig().SetStrict(true)

	doc, err := ParseWithConfig("testdata/valid/basic.toml", cfg)
	verify.NoError(t, err)
	verify.NotNil(t, doc)

	title, err := doc.GetString("title")
	verify.NoError(t, err)
	verify.Equal(t, title, "TOML Example")
}

// TestParseWithInvalidConfig tests parsing with invalid configuration.
func TestParseWithInvalidConfig(t *testing.T) {
	cfg := NewConfig().SetMaxDepth(0) // Invalid

	_, err := ParseWithConfig("testdata/valid/basic.toml", cfg)
	verify.Error(t, err)
}
