package toml

// Copyright (C) 2024 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

import "fmt"

// Config configures the TOML parser behavior.
type Config struct {
	strict             bool
	allowDuplicateKeys bool
	maxDepth           int
	err                error
}

// NewConfig creates a new parser configuration with default settings.
func NewConfig() *Config {
	return &Config{
		strict:             false,
		allowDuplicateKeys: false,
		maxDepth:           100,
	}
}

// SetStrict sets whether to enforce strict TOML v1.0.0 compliance.
// In strict mode, any ambiguities or extensions are rejected.
func (c *Config) SetStrict(strict bool) *Config {
	c.strict = strict
	return c
}

// SetAllowDuplicateKeys sets whether duplicate keys are allowed.
// If true, later definitions override earlier ones.
// If false, duplicate keys cause an error.
func (c *Config) SetAllowDuplicateKeys(allow bool) *Config {
	c.allowDuplicateKeys = allow
	return c
}

// SetMaxDepth sets the maximum nesting depth for tables.
// This prevents stack overflow from maliciously deep TOML files.
func (c *Config) SetMaxDepth(depth int) *Config {
	if depth < 1 {
		c.wrapError(fmt.Errorf("max depth must be at least 1"))
		return c
	}
	c.maxDepth = depth
	return c
}

// Strict returns whether strict mode is enabled.
func (c *Config) Strict() bool {
	return c.strict
}

// AllowDuplicateKeys returns whether duplicate keys are allowed.
func (c *Config) AllowDuplicateKeys() bool {
	return c.allowDuplicateKeys
}

// MaxDepth returns the maximum nesting depth.
func (c *Config) MaxDepth() int {
	return c.maxDepth
}

// Validate checks if the configuration is valid.
// Returns any accumulated errors.
func (c *Config) Validate() error {
	return c.err
}

// Error is an alias for Validate.
func (c *Config) Error() error {
	return c.Validate()
}

// wrapError accumulates errors during configuration.
func (c *Config) wrapError(err error) {
	if c.err == nil {
		c.err = err
	} else {
		c.err = fmt.Errorf("%v; %v", c.err, err)
	}
}

// ParseWithConfig parses a TOML file with custom configuration.
func ParseWithConfig(filepath string, cfg *Config) (*Document, error) {
	if err := cfg.Validate(); err != nil {
		return nil, NewError("parse", err, ErrInvalid)
	}
	// For now, configuration doesn't affect parsing logic
	// This is a placeholder for future enhancements
	return Parse(filepath)
}
