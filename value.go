package toml

// Copyright (C) 2024 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

import (
	"fmt"
	"time"
)

// ValueKind represents the type of a TOML value.
type ValueKind int

const (
	// KindString represents a string value.
	KindString ValueKind = iota
	// KindInteger represents an integer value.
	KindInteger
	// KindFloat represents a floating-point value.
	KindFloat
	// KindBoolean represents a boolean value.
	KindBoolean
	// KindDatetime represents a datetime value.
	KindDatetime
	// KindArray represents an array value.
	KindArray
	// KindTable represents an inline table.
	KindTable
)

// String returns the string representation of the value kind.
func (vk ValueKind) String() string {
	switch vk {
	case KindString:
		return "string"
	case KindInteger:
		return "integer"
	case KindFloat:
		return "float"
	case KindBoolean:
		return "boolean"
	case KindDatetime:
		return "datetime"
	case KindArray:
		return "array"
	case KindTable:
		return "table"
	default:
		return "unknown"
	}
}

// Value represents a TOML value with its type information.
type Value struct {
	kind ValueKind
	raw  any
}

// NewValue creates a new Value with the given kind and raw data.
func NewValue(kind ValueKind, raw any) *Value {
	return &Value{
		kind: kind,
		raw:  raw,
	}
}

// NewStringValue creates a new string Value.
func NewStringValue(s string) *Value {
	return &Value{kind: KindString, raw: s}
}

// NewIntegerValue creates a new integer Value.
func NewIntegerValue(i int64) *Value {
	return &Value{kind: KindInteger, raw: i}
}

// NewFloatValue creates a new float Value.
func NewFloatValue(f float64) *Value {
	return &Value{kind: KindFloat, raw: f}
}

// NewBooleanValue creates a new boolean Value.
func NewBooleanValue(b bool) *Value {
	return &Value{kind: KindBoolean, raw: b}
}

// NewDatetimeValue creates a new datetime Value.
func NewDatetimeValue(t time.Time) *Value {
	return &Value{kind: KindDatetime, raw: t}
}

// NewArrayValue creates a new array Value.
func NewArrayValue(a []any) *Value {
	return &Value{kind: KindArray, raw: a}
}

// NewTableValue creates a new table Value (for inline tables).
func NewTableValue(t map[string]*Value) *Value {
	return &Value{kind: KindTable, raw: t}
}

// Kind returns the kind of the value.
func (v *Value) Kind() ValueKind {
	return v.kind
}

// Raw returns the raw value.
func (v *Value) Raw() any {
	return v.raw
}

// String converts the value to a string.
// Returns an error if the value is not a string.
func (v *Value) String() (string, error) {
	if v.kind != KindString {
		return "", NewError("string conversion", fmt.Errorf("expected string, got %v", v.kind), ErrType)
	}
	return v.raw.(string), nil
}

// Int converts the value to an int64.
// Returns an error if the value is not an integer.
func (v *Value) Int() (int64, error) {
	if v.kind != KindInteger {
		return 0, NewError("int conversion", fmt.Errorf("expected integer, got %v", v.kind), ErrType)
	}
	return v.raw.(int64), nil
}

// Float converts the value to a float64.
// Returns an error if the value is not a float.
func (v *Value) Float() (float64, error) {
	if v.kind != KindFloat {
		return 0, NewError("float conversion", fmt.Errorf("expected float, got %v", v.kind), ErrType)
	}
	return v.raw.(float64), nil
}

// Bool converts the value to a boolean.
// Returns an error if the value is not a boolean.
func (v *Value) Bool() (bool, error) {
	if v.kind != KindBoolean {
		return false, NewError("bool conversion", fmt.Errorf("expected boolean, got %v", v.kind), ErrType)
	}
	return v.raw.(bool), nil
}

// Time converts the value to a time.Time.
// Returns an error if the value is not a datetime.
func (v *Value) Time() (time.Time, error) {
	if v.kind != KindDatetime {
		return time.Time{}, NewError("time conversion", fmt.Errorf("expected datetime, got %v", v.kind), ErrType)
	}
	return v.raw.(time.Time), nil
}

// Array converts the value to a slice.
// Returns an error if the value is not an array.
func (v *Value) Array() ([]any, error) {
	if v.kind != KindArray {
		return nil, NewError("array conversion", fmt.Errorf("expected array, got %v", v.kind), ErrType)
	}
	return v.raw.([]any), nil
}

// Table converts the value to a map (for inline tables).
// Returns an error if the value is not a table.
func (v *Value) Table() (map[string]*Value, error) {
	if v.kind != KindTable {
		return nil, NewError("table conversion", fmt.Errorf("expected table, got %v", v.kind), ErrType)
	}
	return v.raw.(map[string]*Value), nil
}
