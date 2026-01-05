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

// TestValueKindString tests the String method of ValueKind.
func TestValueKindString(t *testing.T) {
	tests := []struct {
		name string
		kind ValueKind
		want string
	}{
		{"string", KindString, "string"},
		{"integer", KindInteger, "integer"},
		{"float", KindFloat, "float"},
		{"boolean", KindBoolean, "boolean"},
		{"datetime", KindDatetime, "datetime"},
		{"array", KindArray, "array"},
		{"table", KindTable, "table"},
		{"unknown", ValueKind(999), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.kind.String()
			verify.Equal(t, got, tt.want)
		})
	}
}

// TestStringValue tests string value creation and conversion.
func TestStringValue(t *testing.T) {
	v := NewStringValue("hello")
	verify.Equal(t, v.Kind(), KindString)

	s, err := v.String()
	verify.NoError(t, err)
	verify.Equal(t, s, "hello")

	// Test wrong type conversions
	_, err = v.Int()
	verify.Error(t, err)

	_, err = v.Float()
	verify.Error(t, err)

	_, err = v.Bool()
	verify.Error(t, err)
}

// TestIntegerValue tests integer value creation and conversion.
func TestIntegerValue(t *testing.T) {
	v := NewIntegerValue(42)
	verify.Equal(t, v.Kind(), KindInteger)

	i, err := v.Int()
	verify.NoError(t, err)
	verify.Equal(t, i, int64(42))

	// Test wrong type conversions
	_, err = v.String()
	verify.Error(t, err)

	_, err = v.Float()
	verify.Error(t, err)

	_, err = v.Bool()
	verify.Error(t, err)
}

// TestFloatValue tests float value creation and conversion.
func TestFloatValue(t *testing.T) {
	v := NewFloatValue(3.14)
	verify.Equal(t, v.Kind(), KindFloat)

	f, err := v.Float()
	verify.NoError(t, err)
	verify.Equal(t, f, 3.14)

	// Test wrong type conversions
	_, err = v.String()
	verify.Error(t, err)

	_, err = v.Int()
	verify.Error(t, err)

	_, err = v.Bool()
	verify.Error(t, err)
}

// TestBooleanValue tests boolean value creation and conversion.
func TestBooleanValue(t *testing.T) {
	v := NewBooleanValue(true)
	verify.Equal(t, v.Kind(), KindBoolean)

	b, err := v.Bool()
	verify.NoError(t, err)
	verify.True(t, b)

	// Test wrong type conversions
	_, err = v.String()
	verify.Error(t, err)

	_, err = v.Int()
	verify.Error(t, err)

	_, err = v.Float()
	verify.Error(t, err)
}

// TestDatetimeValue tests datetime value creation and conversion.
func TestDatetimeValue(t *testing.T) {
	now := time.Now()
	v := NewDatetimeValue(now)
	verify.Equal(t, v.Kind(), KindDatetime)

	tm, err := v.Time()
	verify.NoError(t, err)
	verify.True(t, tm.Equal(now))

	// Test wrong type conversions
	_, err = v.String()
	verify.Error(t, err)

	_, err = v.Int()
	verify.Error(t, err)

	_, err = v.Bool()
	verify.Error(t, err)
}

// TestArrayValue tests array value creation and conversion.
func TestArrayValue(t *testing.T) {
	arr := []any{1, 2, 3}
	v := NewArrayValue(arr)
	verify.Equal(t, v.Kind(), KindArray)

	a, err := v.Array()
	verify.NoError(t, err)
	verify.Equal(t, len(a), 3)
	verify.Equal(t, a[0], 1)
	verify.Equal(t, a[1], 2)
	verify.Equal(t, a[2], 3)

	// Test wrong type conversions
	_, err = v.String()
	verify.Error(t, err)

	_, err = v.Int()
	verify.Error(t, err)

	_, err = v.Bool()
	verify.Error(t, err)
}

// TestTableValue tests table value creation and conversion.
func TestTableValue(t *testing.T) {
	tbl := map[string]*Value{
		"key1": NewStringValue("value1"),
		"key2": NewIntegerValue(42),
	}
	v := NewTableValue(tbl)
	verify.Equal(t, v.Kind(), KindTable)

	tb, err := v.Table()
	verify.NoError(t, err)
	verify.Equal(t, len(tb), 2)

	val1, err := tb["key1"].String()
	verify.NoError(t, err)
	verify.Equal(t, val1, "value1")

	val2, err := tb["key2"].Int()
	verify.NoError(t, err)
	verify.Equal(t, val2, int64(42))

	// Test wrong type conversions
	_, err = v.String()
	verify.Error(t, err)

	_, err = v.Int()
	verify.Error(t, err)

	_, err = v.Array()
	verify.Error(t, err)
}

// TestValueRaw tests the Raw method.
func TestValueRaw(t *testing.T) {
	v := NewStringValue("test")
	raw := v.Raw()
	verify.Equal(t, raw, "test")

	v2 := NewIntegerValue(100)
	raw2 := v2.Raw()
	verify.Equal(t, raw2.(int64), int64(100))
}

// TestNewValue tests the generic NewValue constructor.
func TestNewValue(t *testing.T) {
	v := NewValue(KindString, "generic")
	verify.Equal(t, v.Kind(), KindString)

	s, err := v.String()
	verify.NoError(t, err)
	verify.Equal(t, s, "generic")
}
