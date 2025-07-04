package parser_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/manuelarte/milogo/internal/parser"
	"github.com/manuelarte/milogo/pkg/errors"
)

func TestFilter(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		jsonData        map[string]any
		jsonFieldObject parser.JSONFieldObject
		expected        map[string]any
	}{
		"no nested json,querying one not existing field": {
			jsonData:        map[string]any{"name": "Manuel", "age": 99},
			jsonFieldObject: parser.JSONFieldObject{"surname": parser.JSONFieldValue{}},
			expected:        map[string]any{},
		},
		"no nested json,querying one field": {
			jsonData:        map[string]any{"name": "Manuel", "age": 99},
			jsonFieldObject: parser.JSONFieldObject{"name": parser.JSONFieldValue{}},
			expected:        map[string]any{"name": "Manuel"},
		},
		"no nested json,querying two field": {
			jsonData:        map[string]any{"name": "Manuel", "age": 99},
			jsonFieldObject: parser.JSONFieldObject{"name": parser.JSONFieldValue{}, "age": parser.JSONFieldValue{}},
			expected:        map[string]any{"name": "Manuel", "age": 99},
		},
		"nested json,querying all nested field": {
			jsonData: map[string]any{
				"name": "Manuel", "age": 99,
				"address": map[string]any{"street": "mystreet", "number": 1},
			},
			jsonFieldObject: parser.JSONFieldObject{
				"name": parser.JSONFieldValue{}, "age": parser.JSONFieldValue{},
				"address": parser.JSONFieldValue{},
			},
			expected: map[string]any{
				"name": "Manuel", "age": 99,
				"address": map[string]any{"street": "mystreet", "number": 1},
			},
		},
		"nested array json,querying all nested fields": {
			jsonData: map[string]any{
				"name": "Manuel", "age": 99,
				"addresses": []map[string]any{{"street": "mystreet", "number": 1}},
			},
			jsonFieldObject: parser.JSONFieldObject{
				"name": parser.JSONFieldValue{},
				"age":  parser.JSONFieldValue{}, "addresses": parser.JSONFieldValue{},
			},
			expected: map[string]any{
				"name": "Manuel", "age": 99,
				"addresses": []map[string]any{{"street": "mystreet", "number": 1}},
			},
		},
		"nested array json,querying one nested field": {
			jsonData: map[string]any{
				"name": "Manuel", "age": 99,
				"addresses": []map[string]any{{"street": "mystreet", "number": 1}},
			},
			jsonFieldObject: parser.JSONFieldObject{
				"name": parser.JSONFieldValue{}, "age": parser.JSONFieldValue{},
				"addresses": parser.JSONFieldObject{"street": parser.JSONFieldValue{}},
			},
			expected: map[string]any{
				"name": "Manuel", "age": 99,
				"addresses": []map[string]any{{"street": "mystreet"}},
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := parser.Filter(test.jsonData, test.jsonFieldObject)
			if assert.NoError(t, err) {
				assert.Equal(t, test.expected, test.jsonData)
			}
		})
	}
}

// TODO Test starting object is an array

func TestFilterErrors(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		jsonData        map[string]any
		jsonFieldObject parser.JSONFieldObject
		expected        error
	}{
		"querying not existing nested json": {
			jsonData: map[string]any{"name": "Manuel", "age": 99},
			jsonFieldObject: parser.JSONFieldObject{"name": parser.JSONFieldObject{
				"a": parser.JSONFieldValue{},
				"b": parser.JSONFieldValue{},
			}},
			expected: errors.NotAnObjectError("name"),
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := parser.Filter(test.jsonData, test.jsonFieldObject)
			if assert.Error(t, err) {
				assert.Equal(t, test.expected, err)
			}
		})
	}
}
