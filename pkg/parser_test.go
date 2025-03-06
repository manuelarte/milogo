package pkg_test

import (
	"testing"

	"github.com/manuelarte/milogo/pkg"

	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		fields   string
		expected pkg.JSONFieldObject
	}{
		"querying simple fields": {
			fields:   "name,age",
			expected: pkg.JSONFieldObject{"name": pkg.JSONFieldValue{}, "age": pkg.JSONFieldValue{}},
		},
		"querying simple and complex field": {
			fields: "name,address(street,number)",
			expected: pkg.JSONFieldObject{
				"name":    pkg.JSONFieldValue{},
				"address": pkg.JSONFieldObject{"street": pkg.JSONFieldValue{}, "number": pkg.JSONFieldValue{}},
			},
		},
		"querying simple and nested complex fields": {
			fields: "a,b(c,d,e(f,g))",
			expected: pkg.JSONFieldObject{
				"a": pkg.JSONFieldValue{},
				"b": pkg.JSONFieldObject{
					"c": pkg.JSONFieldValue{}, "d": pkg.JSONFieldValue{},
					"e": pkg.JSONFieldObject{"f": pkg.JSONFieldValue{}, "g": pkg.JSONFieldValue{}},
				},
			},
		},
		"querying complex fields in the middle": {
			fields: "a,b(c,d,e),f,g",
			expected: pkg.JSONFieldObject{
				"a": pkg.JSONFieldValue{},
				"b": pkg.JSONFieldObject{"c": pkg.JSONFieldValue{}, "d": pkg.JSONFieldValue{}, "e": pkg.JSONFieldValue{}},
				"f": pkg.JSONFieldValue{}, "g": pkg.JSONFieldValue{},
			},
		},
		"querying complex fields in several places": {
			fields: "a,b(c,d),e(f,g)",
			expected: pkg.JSONFieldObject{
				"a": pkg.JSONFieldValue{},
				"b": pkg.JSONFieldObject{"c": pkg.JSONFieldValue{}, "d": pkg.JSONFieldValue{}},
				"e": pkg.JSONFieldObject{"f": pkg.JSONFieldValue{}, "g": pkg.JSONFieldValue{}},
			},
		},
	}
	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			parser := pkg.NewParser()
			actual, err := parser.Parse(test.fields)
			if assert.NoError(t, err) {
				assert.Equal(t, test.expected, actual)
			}
		})
	}
}
