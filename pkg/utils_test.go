package pkg_test

import (
	"testing"

	"github.com/manuelarte/milogo/pkg"

	"github.com/stretchr/testify/assert"
)

func TestFilter(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		jsonData        map[string]interface{}
		jsonFieldObject pkg.JSONFieldObject
		expected        map[string]interface{}
	}{
		"no nested json,querying one not existing field": {
			jsonData:        map[string]interface{}{"name": "Manuel", "age": 99},
			jsonFieldObject: pkg.JSONFieldObject{"surname": pkg.JSONFieldValue{}},
			expected:        map[string]interface{}{},
		},
		"no nested json,querying one field": {
			jsonData:        map[string]interface{}{"name": "Manuel", "age": 99},
			jsonFieldObject: pkg.JSONFieldObject{"name": pkg.JSONFieldValue{}},
			expected:        map[string]interface{}{"name": "Manuel"},
		},
		"no nested json,querying two field": {
			jsonData:        map[string]interface{}{"name": "Manuel", "age": 99},
			jsonFieldObject: pkg.JSONFieldObject{"name": pkg.JSONFieldValue{}, "age": pkg.JSONFieldValue{}},
			expected:        map[string]interface{}{"name": "Manuel", "age": 99},
		},
		"nested json,querying all nested field": {
			jsonData: map[string]interface{}{
				"name": "Manuel", "age": 99,
				"address": map[string]interface{}{"street": "mystreet", "number": 1},
			},
			jsonFieldObject: pkg.JSONFieldObject{
				"name": pkg.JSONFieldValue{}, "age": pkg.JSONFieldValue{},
				"address": pkg.JSONFieldValue{},
			},
			expected: map[string]interface{}{
				"name": "Manuel", "age": 99,
				"address": map[string]interface{}{"street": "mystreet", "number": 1},
			},
		},
		"nested array json,querying all nested fields": {
			jsonData: map[string]interface{}{
				"name": "Manuel", "age": 99,
				"addresses": []map[string]interface{}{{"street": "mystreet", "number": 1}},
			},
			jsonFieldObject: pkg.JSONFieldObject{
				"name": pkg.JSONFieldValue{},
				"age":  pkg.JSONFieldValue{}, "addresses": pkg.JSONFieldValue{},
			},
			expected: map[string]interface{}{
				"name": "Manuel", "age": 99,
				"addresses": []map[string]interface{}{{"street": "mystreet", "number": 1}},
			},
		},
		"nested array json,querying one nested field": {
			jsonData: map[string]interface{}{
				"name": "Manuel", "age": 99,
				"addresses": []map[string]interface{}{{"street": "mystreet", "number": 1}},
			},
			jsonFieldObject: pkg.JSONFieldObject{
				"name": pkg.JSONFieldValue{}, "age": pkg.JSONFieldValue{},
				"addresses": pkg.JSONFieldObject{"street": pkg.JSONFieldValue{}},
			},
			expected: map[string]interface{}{
				"name": "Manuel", "age": 99,
				"addresses": []map[string]interface{}{{"street": "mystreet"}},
			},
		},
	}
	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			err := pkg.Filter(test.jsonData, test.jsonFieldObject)
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
		jsonData        map[string]interface{}
		jsonFieldObject pkg.JSONFieldObject
		expected        error
	}{
		"querying not existing nested json": {
			jsonData: map[string]interface{}{"name": "Manuel", "age": 99},
			jsonFieldObject: pkg.JSONFieldObject{"name": pkg.JSONFieldObject{
				"a": pkg.JSONFieldValue{},
				"b": pkg.JSONFieldValue{},
			}},
			expected: pkg.NotAnObjectError("name"),
		},
	}
	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			err := pkg.Filter(test.jsonData, test.jsonFieldObject)
			if assert.Error(t, err) {
				assert.Equal(t, test.expected, err)
			}
		})
	}
}
