# Tideland Go TOML Parser

[![GitHub release](https://img.shields.io/github/release/tideland/go-toml.svg)](https://github.com/tideland/go-toml)
[![GitHub license](https://img.shields.io/badge/license-New%20BSD-blue.svg)](https://raw.githubusercontent.com/tideland/go-toml/master/LICENSE)
[![Go Module](https://img.shields.io/github/go-mod/go-version/tideland/go-toml)](https://github.com/tideland/go-toml/blob/master/go.mod)
[![GoDoc](https://godoc.org/tideland.dev/go/toml?status.svg)](https://pkg.go.dev/mod/tideland.dev/go/toml?tab=packages)
![Workflow](https://github.com/tideland/go-toml/actions/workflows/build.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/tideland/go-toml)](https://goreportcard.com/report/tideland.dev/go/toml)

## Overview

A standalone TOML v1.0.0 parser for Go, following Tideland standards.

## Features

- **Read-only parser** - Parse TOML files into a navigable document structure
- **Dual access patterns** - Flat dot notation and hierarchical navigation
- **Complete TOML v1.0.0 support** - All data types, tables, arrays, and features
- **Type-safe value access** - Explicit type conversion with error handling
- **Comprehensive error reporting** - Line and column information for parse errors
- **No external dependencies** - Pure Go implementation
- **Well tested** - >90% test coverage

## Installation

```bash
go get tideland.dev/go/toml
```

## Quick Start

### Parse a TOML file

```go
package main

import (
    "fmt"
    "log"

    "tideland.dev/go/toml"
)

func main() {
    doc, err := toml.Parse("config.toml")
    if err != nil {
        log.Fatal(err)
    }

    // Flat access
    host, _ := doc.GetString("server.host")
    port, _ := doc.GetInt("server.port")

    fmt.Printf("Server: %s:%d\n", host, port)
}
```

### Example TOML file

```toml
# config.toml
[server]
host = "localhost"
port = 8080
enabled = true

[database]
connection = "postgres://localhost/mydb"
max_connections = 100
```

## Usage

### Flat Access (Dot Notation)

Simple and direct for known paths:

```go
doc, _ := toml.Parse("config.toml")

// Access nested values with dot notation
timeout, _ := doc.GetInt("database.connection.timeout")
retries, _ := doc.GetInt("database.connection.retries")
enabled, _ := doc.GetBool("database.connection.enabled")
```

### Hierarchical Navigation

Useful for exploring structure:

```go
doc, _ := toml.Parse("config.toml")

// Navigate to a section
db, _ := doc.Section("database")

// List subsections
for _, name := range db.Sections() {
    section, _ := db.Section(name)
    fmt.Printf("Section: %s\n", section.Name())
}

// List keys in a section
for _, key := range db.Keys() {
    fmt.Printf("Key: %s\n", key)
}
```

### Table Arrays

Access arrays of tables:

```go
// TOML:
// [[fruits]]
// name = "apple"
// color = "red"
//
// [[fruits]]
// name = "banana"
// color = "yellow"

fruits, _ := doc.TableArray("fruits")
for _, fruit := range fruits {
    name, _ := fruit.GetString("name")
    color, _ := fruit.GetString("color")
    fmt.Printf("%s is %s\n", name, color)
}
```

### Type Conversion

All value accessors return errors for type mismatches:

```go
// Typed getters
str, err := doc.GetString("key")
num, err := doc.GetInt("key")
flt, err := doc.GetFloat("key")
bln, err := doc.GetBool("key")
tim, err := doc.GetTime("key")
arr, err := doc.GetArray("key")
```

### Configuration

Customize parser behavior:

```go
cfg := toml.NewConfig().
    SetStrict(true).
    SetMaxDepth(100)

doc, err := toml.ParseWithConfig("config.toml", cfg)
```

## API Overview

### Document Methods

```go
// Flat access
GetString(path string) (string, error)
GetInt(path string) (int64, error)
GetFloat(path string) (float64, error)
GetBool(path string) (bool, error)
GetTime(path string) (time.Time, error)
GetArray(path string) ([]any, error)

// Hierarchical access
Section(name string) (*Section, error)
Sections() []string
Keys() []string
TableArray(name string) ([]*Section, error)

// Utilities
Has(path string) bool
Path() string
```

### Section Methods

```go
// Value access
GetString(key string) (string, error)
GetInt(key string) (int64, error)
GetFloat(key string) (float64, error)
GetBool(key string) (bool, error)
GetTime(key string) (time.Time, error)
GetArray(key string) ([]any, error)

// Navigation
Section(name string) (*Section, error)
Sections() []string
Keys() []string
TableArray(name string) ([]*Section, error)

// Metadata
Name() string
FullPath() []string
Parent() *Section
Has(key string) bool
```

## Supported TOML Features

### Data Types

- **Strings**: Basic, multi-line, literal, multi-line literal
- **Integers**: Decimal, hex (`0x`), octal (`0o`), binary (`0b`)
- **Floats**: Standard, exponent, special values (`inf`, `-inf`, `nan`)
- **Booleans**: `true`, `false`
- **Datetimes**: Offset date-time, local date-time, local date, local time
- **Arrays**: Homogeneous and mixed
- **Inline tables**: `{key = value, ...}`

### Tables

- Simple tables: `[section]`
- Nested tables: `[a.b.c]`
- Dotted keys: `a.b.c = value`
- Array of tables: `[[array]]`

### Other Features

- Comments (`#`)
- Escape sequences in strings
- Underscores in numbers (`1_000_000`)
- Multi-line strings with backslash line endings

## Error Handling

Detailed error information with context:

```go
doc, err := toml.Parse("bad.toml")
if err != nil {
    if parseErr, ok := err.(*toml.ParseError); ok {
        fmt.Printf("Error at %s:%d:%d: %v (%v)\n",
            parseErr.Path,
            parseErr.Line,
            parseErr.Col,
            parseErr.Err,
            parseErr.Code)
    }
}
```

Error codes:
- `ErrSyntax` - Syntax error in TOML
- `ErrType` - Type conversion error
- `ErrNotFound` - Key or section not found
- `ErrDuplicate` - Duplicate key or table
- `ErrInvalidPath` - Invalid dot-notation path
- `ErrIO` - File I/O error
- `ErrInvalid` - Invalid configuration

## Testing

Run tests:

```bash
make test
```

Run with coverage:

```bash
make coverage
```

## License

Copyright (C) 2024 Frank Mueller / Tideland / Oldenburg / Germany

All rights reserved. Use of this source code is governed by the new BSD license.

## See Also

- [TOML v1.0.0 Specification](https://toml.io/en/v1.0.0)
- [Tideland Go Packages](https://tideland.dev/go)
