package toml

// Copyright (C) 2024 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

import (
	"testing"

	"tideland.dev/go/asserts/verify"
)

// TestParseSimpleKeyValue tests parsing simple key-value pairs.
func TestParseSimpleKeyValue(t *testing.T) {
	input := `title = "TOML Example"
count = 42
pi = 3.14
enabled = true`

	doc, err := ParseString(input, "test.toml")
	verify.NoError(t, err)
	verify.NotNil(t, doc)

	// Test string value
	title, err := doc.GetString("title")
	verify.NoError(t, err)
	verify.Equal(t, title, "TOML Example")

	// Test integer value
	count, err := doc.GetInt("count")
	verify.NoError(t, err)
	verify.Equal(t, count, int64(42))

	// Test float value
	pi, err := doc.GetFloat("pi")
	verify.NoError(t, err)
	verify.Equal(t, pi, 3.14)

	// Test boolean value
	enabled, err := doc.GetBool("enabled")
	verify.NoError(t, err)
	verify.True(t, enabled)
}

// TestParseSimpleTable tests parsing a simple table.
func TestParseSimpleTable(t *testing.T) {
	input := `[server]
host = "localhost"
port = 8080`

	doc, err := ParseString(input, "test.toml")
	verify.NoError(t, err)

	// Test hierarchical access
	server, err := doc.Section("server")
	verify.NoError(t, err)
	verify.NotNil(t, server)

	host, err := server.GetString("host")
	verify.NoError(t, err)
	verify.Equal(t, host, "localhost")

	port, err := server.GetInt("port")
	verify.NoError(t, err)
	verify.Equal(t, port, int64(8080))

	// Test flat access with dot notation
	host2, err := doc.GetString("server.host")
	verify.NoError(t, err)
	verify.Equal(t, host2, "localhost")
}

// TestParseNestedTables tests parsing nested tables.
func TestParseNestedTables(t *testing.T) {
	input := `[database.connection]
timeout = 30
retries = 3`

	doc, err := ParseString(input, "test.toml")
	verify.NoError(t, err)

	// Hierarchical navigation
	database, err := doc.Section("database")
	verify.NoError(t, err)

	connection, err := database.Section("connection")
	verify.NoError(t, err)

	timeout, err := connection.GetInt("timeout")
	verify.NoError(t, err)
	verify.Equal(t, timeout, int64(30))

	// Flat access
	retries, err := doc.GetInt("database.connection.retries")
	verify.NoError(t, err)
	verify.Equal(t, retries, int64(3))
}

// TestParseArray tests parsing arrays.
func TestParseArray(t *testing.T) {
	input := `numbers = [1, 2, 3, 4, 5]
strings = ["a", "b", "c"]
mixed = [1, "two", 3.0]`

	doc, err := ParseString(input, "test.toml")
	verify.NoError(t, err)

	numbers, err := doc.GetArray("numbers")
	verify.NoError(t, err)
	verify.Equal(t, len(numbers), 5)
	verify.Equal(t, numbers[0].(int64), int64(1))
	verify.Equal(t, numbers[4].(int64), int64(5))

	strings, err := doc.GetArray("strings")
	verify.NoError(t, err)
	verify.Equal(t, len(strings), 3)
	verify.Equal(t, strings[0], "a")

	mixed, err := doc.GetArray("mixed")
	verify.NoError(t, err)
	verify.Equal(t, len(mixed), 3)
}

// TestParseInlineTable tests parsing inline tables.
func TestParseInlineTable(t *testing.T) {
	input := `point = {x = 1, y = 2}
person = {name = "John", age = 30}`

	doc, err := ParseString(input, "test.toml")
	verify.NoError(t, err)

	point, err := doc.root.keys["point"].Table()
	verify.NoError(t, err)
	verify.Equal(t, len(point), 2)

	x, err := point["x"].Int()
	verify.NoError(t, err)
	verify.Equal(t, x, int64(1))

	person, err := doc.root.keys["person"].Table()
	verify.NoError(t, err)

	name, err := person["name"].String()
	verify.NoError(t, err)
	verify.Equal(t, name, "John")
}

// TestParseDottedKeys tests parsing dotted keys.
func TestParseDottedKeys(t *testing.T) {
	input := `physical.color = "orange"
physical.shape = "round"
site."google.com" = true`

	doc, err := ParseString(input, "test.toml")
	verify.NoError(t, err)

	// Dotted keys create nested structure
	physical, err := doc.Section("physical")
	verify.NoError(t, err)

	color, err := physical.GetString("color")
	verify.NoError(t, err)
	verify.Equal(t, color, "orange")

	shape, err := physical.GetString("shape")
	verify.NoError(t, err)
	verify.Equal(t, shape, "round")

	// Using flat access
	color2, err := doc.GetString("physical.color")
	verify.NoError(t, err)
	verify.Equal(t, color2, "orange")
}

// TestParseComments tests that comments are properly ignored.
func TestParseComments(t *testing.T) {
	input := `# This is a comment
key = "value" # inline comment
# Another comment
[section] # section comment
name = "test"`

	doc, err := ParseString(input, "test.toml")
	verify.NoError(t, err)

	val, err := doc.GetString("key")
	verify.NoError(t, err)
	verify.Equal(t, val, "value")

	section, err := doc.Section("section")
	verify.NoError(t, err)

	name, err := section.GetString("name")
	verify.NoError(t, err)
	verify.Equal(t, name, "test")
}

// TestParseMultilineString tests multi-line strings.
func TestParseMultilineString(t *testing.T) {
	input := `str = """
Line one
Line two
Line three"""

literal = '''
No escape \n here'''`

	doc, err := ParseString(input, "test.toml")
	verify.NoError(t, err)

	str, err := doc.GetString("str")
	verify.NoError(t, err)
	verify.True(t, len(str) > 0)

	literal, err := doc.GetString("literal")
	verify.NoError(t, err)
	verify.True(t, len(literal) > 0)
}

// TestParseSpecialNumbers tests special float values.
func TestParseSpecialNumbers(t *testing.T) {
	input := `inf_pos = inf
inf_neg = -inf
not_a_num = nan`

	doc, err := ParseString(input, "test.toml")
	verify.NoError(t, err)

	// Just verify they parse without error
	_, err = doc.GetFloat("inf_pos")
	verify.NoError(t, err)

	_, err = doc.GetFloat("inf_neg")
	verify.NoError(t, err)

	_, err = doc.GetFloat("not_a_num")
	verify.NoError(t, err)
}

// TestParseFromFile tests parsing from an actual file.
func TestParseFromFile(t *testing.T) {
	doc, err := Parse("testdata/valid/basic.toml")
	verify.NoError(t, err)
	verify.NotNil(t, doc)

	title, err := doc.GetString("title")
	verify.NoError(t, err)
	verify.Equal(t, title, "TOML Example")

	host, err := doc.GetString("server.host")
	verify.NoError(t, err)
	verify.Equal(t, host, "localhost")
}

// TestParseNestedFromFile tests parsing nested structure from file.
func TestParseNestedFromFile(t *testing.T) {
	doc, err := Parse("testdata/valid/nested.toml")
	verify.NoError(t, err)

	timeout, err := doc.GetInt("database.connection.timeout")
	verify.NoError(t, err)
	verify.Equal(t, timeout, int64(30))

	maxSize, err := doc.GetInt("database.pool.max_size")
	verify.NoError(t, err)
	verify.Equal(t, maxSize, int64(20))
}

// TestDocumentHas tests the Has method.
func TestDocumentHas(t *testing.T) {
	input := `key = "value"
[section]
inner = "test"`

	doc, err := ParseString(input, "test.toml")
	verify.NoError(t, err)

	verify.True(t, doc.Has("key"))
	verify.True(t, doc.Has("section.inner"))
	verify.False(t, doc.Has("nonexistent"))
	verify.False(t, doc.Has("section.missing"))
}

// TestDocumentKeysAndSections tests Keys and Sections methods.
func TestDocumentKeysAndSections(t *testing.T) {
	input := `key1 = "value1"
key2 = "value2"

[section1]
a = 1

[section2]
b = 2`

	doc, err := ParseString(input, "test.toml")
	verify.NoError(t, err)

	keys := doc.Keys()
	verify.Equal(t, len(keys), 2)

	sections := doc.Sections()
	verify.Equal(t, len(sections), 2)
}

// TestSectionMetadata tests Section metadata methods.
func TestSectionMetadata(t *testing.T) {
	input := `[database.connection]
timeout = 30`

	doc, err := ParseString(input, "test.toml")
	verify.NoError(t, err)

	db, err := doc.Section("database")
	verify.NoError(t, err)

	verify.Equal(t, db.Name(), "database")
	verify.Equal(t, len(db.FullPath()), 1)
	verify.Equal(t, db.FullPath()[0], "database")

	conn, err := db.Section("connection")
	verify.NoError(t, err)

	verify.Equal(t, conn.Name(), "connection")
	verify.Equal(t, len(conn.FullPath()), 2)
	verify.Equal(t, conn.FullPath()[0], "database")
	verify.Equal(t, conn.FullPath()[1], "connection")

	verify.Equal(t, conn.Parent(), db)
}

// TestParseTableArray tests parsing array of tables.
func TestParseTableArray(t *testing.T) {
	input := `[[fruits]]
name = "apple"
color = "red"

[[fruits]]
name = "banana"
color = "yellow"`

	doc, err := ParseString(input, "test.toml")
	verify.NoError(t, err)

	fruits, err := doc.TableArray("fruits")
	verify.NoError(t, err)
	verify.Equal(t, len(fruits), 2)

	name0, err := fruits[0].GetString("name")
	verify.NoError(t, err)
	verify.Equal(t, name0, "apple")

	color0, err := fruits[0].GetString("color")
	verify.NoError(t, err)
	verify.Equal(t, color0, "red")

	name1, err := fruits[1].GetString("name")
	verify.NoError(t, err)
	verify.Equal(t, name1, "banana")

	color1, err := fruits[1].GetString("color")
	verify.NoError(t, err)
	verify.Equal(t, color1, "yellow")
}

// TestParseTableArrayFromFile tests parsing table arrays from file.
func TestParseTableArrayFromFile(t *testing.T) {
	doc, err := Parse("testdata/valid/table_arrays.toml")
	verify.NoError(t, err)

	fruits, err := doc.TableArray("fruits")
	verify.NoError(t, err)
	verify.Equal(t, len(fruits), 3)

	// Verify names
	names := []string{"apple", "banana", "orange"}
	for i, fruit := range fruits {
		name, err := fruit.GetString("name")
		verify.NoError(t, err)
		verify.Equal(t, name, names[i])
	}
}

// TestParseDatetimeFormats tests various datetime formats.
func TestParseDatetimeFormats(t *testing.T) {
	input := `offset = 1979-05-27T07:32:00Z
local = 1979-05-27T07:32:00
date = 1979-05-27
time = 07:32:00`

	doc, err := ParseString(input, "test.toml")
	verify.NoError(t, err)

	offset, err := doc.GetTime("offset")
	verify.NoError(t, err)
	verify.Equal(t, offset.Year(), 1979)

	local, err := doc.GetTime("local")
	verify.NoError(t, err)
	verify.Equal(t, local.Year(), 1979)

	date, err := doc.GetTime("date")
	verify.NoError(t, err)
	verify.Equal(t, date.Year(), 1979)

	timeVal, err := doc.GetTime("time")
	verify.NoError(t, err)
	verify.Equal(t, timeVal.Hour(), 7)
}

// TestParseHexOctalBinary tests hex, octal, and binary integers.
func TestParseHexOctalBinary(t *testing.T) {
	input := `hex = 0xDEADBEEF
octal = 0o755
binary = 0b11010110`

	doc, err := ParseString(input, "test.toml")
	verify.NoError(t, err)

	hex, err := doc.GetInt("hex")
	verify.NoError(t, err)
	verify.True(t, hex > 0)

	octal, err := doc.GetInt("octal")
	verify.NoError(t, err)
	verify.Equal(t, octal, int64(0755))

	binary, err := doc.GetInt("binary")
	verify.NoError(t, err)
	verify.True(t, binary > 0)
}

// TestParseErrors tests various parse error cases.
func TestParseErrors(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"missing bracket", "[section\nkey = value"},
		{"invalid token", "key value"},
		{"bad array", "[1, 2,"},
		{"bad inline table", "{key = value"},
		{"unexpected token", "= value"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseString(tt.input, "test.toml")
			verify.Error(t, err)
		})
	}
}

// TestParseFileNotFound tests error when file doesn't exist.
func TestParseFileNotFound(t *testing.T) {
	_, err := Parse("nonexistent.toml")
	verify.Error(t, err)
}

// TestParseEmptyArray tests parsing empty arrays.
func TestParseEmptyArray(t *testing.T) {
	input := `empty = []`

	doc, err := ParseString(input, "test.toml")
	verify.NoError(t, err)

	arr, err := doc.GetArray("empty")
	verify.NoError(t, err)
	verify.Equal(t, len(arr), 0)
}

// TestParseEmptyInlineTable tests parsing empty inline tables.
func TestParseEmptyInlineTable(t *testing.T) {
	input := `empty = {}`

	doc, err := ParseString(input, "test.toml")
	verify.NoError(t, err)

	val := doc.root.keys["empty"]
	table, err := val.Table()
	verify.NoError(t, err)
	verify.Equal(t, len(table), 0)
}

// TestParseArrayWithNewlines tests arrays spanning multiple lines.
func TestParseArrayWithNewlines(t *testing.T) {
	input := `numbers = [
  1,
  2,
  3
]`

	doc, err := ParseString(input, "test.toml")
	verify.NoError(t, err)

	arr, err := doc.GetArray("numbers")
	verify.NoError(t, err)
	verify.Equal(t, len(arr), 3)
}
