package parser_test

import (
	"errors"
	"testing"

	"github.com/devbackend/crutch-db/internal/operation"
	"github.com/devbackend/crutch-db/internal/parser"
	"github.com/stretchr/testify/assert"
)

func TestParser_ParseValidOperations(t *testing.T) {
	cases := []struct {
		name          string
		str           string
		expectedType  operation.Type
		expectedKey   string
		expectedValue interface{}
	}{
		{
			name:         "GET operation",
			str:          "GET getKeyName",
			expectedType: operation.Get,
			expectedKey:  "getKeyName",
		},
		{
			name:          "SET operation",
			str:           "SET setKeyName 123",
			expectedType:  operation.Set,
			expectedKey:   "setKeyName",
			expectedValue: 123,
		},
		{
			name:         "DELETE operation",
			str:          "DELETE delKeyName",
			expectedType: operation.Delete,
			expectedKey:  "delKeyName",
		},
		{
			name:         "KEYS operation",
			str:          "KEYS",
			expectedType: operation.Keys,
		}, {
			name:         "GET operation - changed letter case",
			str:          "Get getKeyName",
			expectedType: operation.Get,
			expectedKey:  "getKeyName",
		},
		{
			name:          "SET operation - changed letter case",
			str:           "sEt setKeyName 123",
			expectedType:  operation.Set,
			expectedKey:   "setKeyName",
			expectedValue: 123,
		},
		{
			name:         "DELETE operation - changed letter case",
			str:          "delETe delKeyName",
			expectedType: operation.Delete,
			expectedKey:  "delKeyName",
		},
		{
			name:         "KEYS operation - changed letter case",
			str:          "keYS",
			expectedType: operation.Keys,
		},
	}
	for _, c := range cases {
		if c.name == "" {
			t.Errorf("test case name required!")
			continue
		}

		t.Run(c.name, func(t *testing.T) {
			p := parser.New()

			op, _ := p.Parse(c.str)

			assert.Equal(t, c.expectedType, op.Type)
			assert.Equal(t, c.expectedKey, op.Key)
			assert.Equal(t, c.expectedValue, op.Value)
		})
	}
}

func TestParser_ParseInvalidOperations(t *testing.T) {
	cases := []struct {
		name        string
		str         string
		expectedErr error
	}{
		{
			name:        "Empty request",
			str:         "",
			expectedErr: errors.New("empty request"),
		},
		{
			name:        "GET without key name",
			str:         "GET",
			expectedErr: errors.New("empty key mame"),
		},
		{
			name:        "SET without key name",
			str:         "SET",
			expectedErr: errors.New("empty key mame"),
		},
		{
			name:        "DELETE without key name",
			str:         "DELETE",
			expectedErr: errors.New("empty key mame"),
		},
		{
			name:        "SET without value",
			str:         "SET keyName",
			expectedErr: errors.New("empty value for set"),
		},
		{
			name:        "SET without value by space",
			str:         "SET keyName ",
			expectedErr: errors.New("empty value for set"),
		},
		{
			name:        "Unknown operation",
			str:         "FOOBAR",
			expectedErr: errors.New("unknown operation FOOBAR"),
		},
	}
	for _, c := range cases {
		if c.name == "" {
			t.Errorf("test case name required!")
			continue
		}

		t.Run(c.name, func(t *testing.T) {
			p := parser.New()

			_, err := p.Parse(c.str)

			assert.Equal(t, c.expectedErr, err)
		})
	}
}

func TestParser_ParseSetValidValues(t *testing.T) {
	cases := []struct {
		name     string
		value    string
		expected interface{}
	}{
		{
			name:     "positive int",
			value:    "12345",
			expected: 12345,
		},
		{
			name:     "negative int",
			value:    "-12345",
			expected: -12345,
		},
		{
			name:     "zero",
			value:    "0",
			expected: 0,
		},
		{
			name:     "positive float",
			value:    "123.45",
			expected: 123.45,
		},
		{
			name:     "negative float",
			value:    "-123.45",
			expected: -123.45,
		},
		{
			name:     "empty string",
			value:    "''",
			expected: "",
		},
		{
			name:     "valid string: one word",
			value:    "'hello!'",
			expected: "hello!",
		},
		{
			name:     "valid string: phrase",
			value:    "'hello world!'",
			expected: "hello world!",
		},
		{
			name:     "shielding for quote",
			value:    "'it\\'s my life!'",
			expected: "it's my life!",
		},
		{
			name:     "empty array",
			value:    "[]",
			expected: []interface{}{},
		},
		{
			name:     "array of int values",
			value:    "[1 ,2,3, -4]",
			expected: []interface{}{1, 2, 3, -4},
		},
		{
			name:     "array of string values",
			value:    "['hello' ,'world','from', 'me']",
			expected: []interface{}{"hello", "world", "from", "me"},
		},
	}
	for _, c := range cases {
		if c.name == "" {
			t.Errorf("test case name required!")
			continue
		}

		t.Run(c.name, func(t *testing.T) {
			p := parser.New()

			op, _ := p.Parse("SET key-name " + c.value)

			assert.Equal(t, c.expected, op.Value)
		})
	}
}

func TestParser_ParseSetInvalidValues(t *testing.T) {
	cases := []struct {
		name     string
		value    string
		expected error
	}{
		{
			name:     "incomplete string",
			value:    "'hello world!",
			expected: errors.New("incomplete string"),
		},
		{
			name:     "incomplete string - only quote",
			value:    "'",
			expected: errors.New("incomplete string"),
		},
		{
			name:     "wrong string - quote without shielding",
			value:    "'it's my life!'",
			expected: errors.New("bad string"),
		},
		{
			name:     "incomplete array - only open parenthesis",
			value:    "[",
			expected: errors.New("bad array syntax"),
		},
		{
			name:     "incomplete array",
			value:    "[1, 2, 3",
			expected: errors.New("bad array syntax"),
		},
		{
			name:     "incomplete string in array",
			value:    "['hello', 'world]",
			expected: errors.New("bad array syntax: incomplete string"),
		},
	}
	for _, c := range cases {
		if c.name == "" {
			t.Errorf("test case name required!")
			continue
		}

		t.Run(c.name, func(t *testing.T) {
			p := parser.New()

			_, err := p.Parse("SET key-name " + c.value)

			assert.Equal(t, c.expected, err)
		})
	}
}
