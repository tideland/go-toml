// Copyright (C) 2024 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

/*
Package toml provides a standalone TOML v1.0.0 parser for Go.

# Overview

This package implements a read-only TOML parser that supports all features
of the TOML v1.0.0 specification. It provides both flat dot-notation access
and hierarchical navigation for working with TOML documents.

# Basic Usage

Parse a TOML file and access values:

	doc, err := toml.Parse("config.toml")
	if err != nil {
		log.Fatal(err)
	}

	// Flat access with dot notation
	host, _ := doc.GetString("server.host")
	port, _ := doc.GetInt("server.port")

	// Hierarchical navigation
	server, _ := doc.Section("server")
	timeout, _ := server.GetInt("timeout")

# Features

  - Read-only parser (no marshaling/writing)
  - Parse from files
  - Flat access via dot notation: doc.GetString("database.connection.timeout")
  - Hierarchical navigation: doc.Section("database").Section("connection")
  - List sections and keys
  - Support for all TOML v1.0.0 data types
  - Table arrays (array of tables)
  - Nested sections
  - Dotted keys
  - Comprehensive error reporting with line/column information

# Data Types

The parser supports all TOML v1.0.0 data types:

  - Strings (basic, multi-line, literal, multi-line literal)
  - Integers (decimal, hex, octal, binary)
  - Floats (standard, exponent, inf, nan)
  - Booleans (true, false)
  - Datetimes (offset, local date-time, local date, local time)
  - Arrays
  - Inline tables

# Access Patterns

Two complementary ways to access TOML data:

Flat Access (Dot Notation):

	// Simple and direct for known paths
	timeout, err := doc.GetInt("database.connection.timeout")
	retries, err := doc.GetInt("database.connection.retries")

Hierarchical Navigation:

	// Useful for exploring structure
	db, err := doc.Section("database")
	for _, name := range db.Sections() {
		section, _ := db.Section(name)
		fmt.Printf("Section: %s\n", section.Name())
		for _, key := range section.Keys() {
			fmt.Printf("  Key: %s\n", key)
		}
	}

# Table Arrays

TOML table arrays are accessed as slices of sections:

	// [[fruits]]
	// name = "apple"
	// color = "red"
	//
	// [[fruits]]
	// name = "banana"
	// color = "yellow"

	fruits, err := doc.TableArray("fruits")
	for i, fruit := range fruits {
		name, _ := fruit.GetString("name")
		color, _ := fruit.GetString("color")
		fmt.Printf("Fruit %d: %s (%s)\n", i, name, color)
	}

# Configuration

Customize parser behavior with configuration:

	cfg := toml.NewConfig().
		SetStrict(true).
		SetMaxDepth(100)

	doc, err := toml.ParseWithConfig("config.toml", cfg)

Configuration options:

  - Strict mode: Enforce strict TOML v1.0.0 compliance
  - Allow duplicate keys: Last definition wins vs error
  - Max depth: Limit nesting depth to prevent stack overflow

# Error Handling

The package provides detailed error information:

	doc, err := toml.Parse("bad.toml")
	if err != nil {
		if parseErr, ok := err.(*toml.ParseError); ok {
			fmt.Printf("Error at %s:%d:%d: %v\n",
				parseErr.Path, parseErr.Line, parseErr.Col, parseErr.Err)
		}
	}

Error codes help categorize errors:

  - ErrSyntax: Syntax errors in TOML
  - ErrType: Type conversion errors
  - ErrNotFound: Key or section not found
  - ErrDuplicate: Duplicate key or table
  - ErrInvalidPath: Invalid dot-notation path
  - ErrIO: File I/O errors
  - ErrInvalid: Invalid configuration

# Example

Complete example:

	package main

	import (
		"fmt"
		"log"

		"tideland.dev/go/toml"
	)

	func main() {
		// Parse TOML file
		doc, err := toml.Parse("app.toml")
		if err != nil {
			log.Fatal(err)
		}

		// Access configuration values
		appName, _ := doc.GetString("app.name")
		version, _ := doc.GetString("app.version")

		fmt.Printf("%s v%s\n", appName, version)

		// Navigate sections
		server, _ := doc.Section("server")
		host, _ := server.GetString("host")
		port, _ := server.GetInt("port")

		fmt.Printf("Server: %s:%d\n", host, port)

		// List all database connections
		databases, _ := doc.TableArray("databases")
		for i, db := range databases {
			name, _ := db.GetString("name")
			connStr, _ := db.GetString("connection_string")
			fmt.Printf("DB %d: %s -> %s\n", i+1, name, connStr)
		}
	}

# TOML v1.0.0 Compliance

This package implements the complete TOML v1.0.0 specification:

  - All data types (strings, integers, floats, booleans, datetimes, arrays, tables)
  - String escape sequences
  - Multi-line strings (basic and literal)
  - Integer formats (decimal, hex, octal, binary)
  - Float special values (inf, -inf, nan)
  - Datetime formats (offset, local)
  - Arrays (homogeneous and mixed)
  - Inline tables
  - Tables and nested tables
  - Dotted keys
  - Array of tables (table arrays)
  - Comments

# Design Philosophy

The package follows Tideland Go standards:

  - Simple things should be simple (flat dot notation access)
  - Complex things should be possible (hierarchical navigation)
  - Clear error messages with context
  - Type safety with explicit conversions
  - Comprehensive testing (>90% coverage)
  - No external dependencies for core functionality

See the test files for more examples and usage patterns.
*/
package toml
