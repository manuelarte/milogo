package fieldparser_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/manuelarte/milogo/internal/fieldparser"
	"github.com/manuelarte/milogo/pkg/errors"
)

func TestFilter(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		jsonData        map[string]any
		jsonFieldObject fieldparser.JSONFieldObject
		expected        map[string]any
	}{
		"no nested json,querying one not existing field": {
			jsonData:        map[string]any{"name": "Manuel", "age": 99},
			jsonFieldObject: fieldparser.JSONFieldObject{"surname": fieldparser.JSONFieldValue{}},
			expected:        map[string]any{},
		},
		"no nested json,querying one field": {
			jsonData:        map[string]any{"name": "Manuel", "age": 99},
			jsonFieldObject: fieldparser.JSONFieldObject{"name": fieldparser.JSONFieldValue{}},
			expected:        map[string]any{"name": "Manuel"},
		},
		"no nested json,querying two field": {
			jsonData: map[string]any{"name": "Manuel", "age": 99},
			jsonFieldObject: fieldparser.JSONFieldObject{
				"name": fieldparser.JSONFieldValue{},
				"age":  fieldparser.JSONFieldValue{},
			},
			expected: map[string]any{"name": "Manuel", "age": 99},
		},
		"nested json,querying all nested field": {
			jsonData: map[string]any{
				"name": "Manuel", "age": 99,
				"address": map[string]any{"street": "mystreet", "number": 1},
			},
			jsonFieldObject: fieldparser.JSONFieldObject{
				"name": fieldparser.JSONFieldValue{}, "age": fieldparser.JSONFieldValue{},
				"address": fieldparser.JSONFieldValue{},
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
			jsonFieldObject: fieldparser.JSONFieldObject{
				"name": fieldparser.JSONFieldValue{},
				"age":  fieldparser.JSONFieldValue{}, "addresses": fieldparser.JSONFieldValue{},
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
			jsonFieldObject: fieldparser.JSONFieldObject{
				"name": fieldparser.JSONFieldValue{}, "age": fieldparser.JSONFieldValue{},
				"addresses": fieldparser.JSONFieldObject{"street": fieldparser.JSONFieldValue{}},
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

			err := fieldparser.Filter(test.jsonData, test.jsonFieldObject)
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
		jsonFieldObject fieldparser.JSONFieldObject
		expected        error
	}{
		"querying not existing nested json": {
			jsonData: map[string]any{"name": "Manuel", "age": 99},
			jsonFieldObject: fieldparser.JSONFieldObject{"name": fieldparser.JSONFieldObject{
				"a": fieldparser.JSONFieldValue{},
				"b": fieldparser.JSONFieldValue{},
			}},
			expected: errors.NotAnObjectError("name"),
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := fieldparser.Filter(test.jsonData, test.jsonFieldObject)
			if assert.Error(t, err) {
				assert.Equal(t, test.expected, err)
			}
		})
	}
}
