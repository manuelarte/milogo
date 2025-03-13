package parser_test

import (
	"testing"

	"github.com/manuelarte/milogo/internal/parser"

	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		fields   string
		expected parser.JSONFieldObject
	}{
		"querying simple fields": {
			fields:   "name,age",
			expected: parser.JSONFieldObject{"name": parser.JSONFieldValue{}, "age": parser.JSONFieldValue{}},
		},
		"querying simple and complex field": {
			fields: "name,address(street,number)",
			expected: parser.JSONFieldObject{
				"name": parser.JSONFieldValue{},
				"address": parser.JSONFieldObject{
					"street": parser.JSONFieldValue{},
					"number": parser.JSONFieldValue{},
				},
			},
		},
		"querying simple and nested complex fields": {
			fields: "a,b(c,d,e(f,g))",
			expected: parser.JSONFieldObject{
				"a": parser.JSONFieldValue{},
				"b": parser.JSONFieldObject{
					"c": parser.JSONFieldValue{}, "d": parser.JSONFieldValue{},
					"e": parser.JSONFieldObject{"f": parser.JSONFieldValue{}, "g": parser.JSONFieldValue{}},
				},
			},
		},
		"querying complex fields in the middle": {
			fields: "a,b(c,d,e),f,g",
			expected: parser.JSONFieldObject{
				"a": parser.JSONFieldValue{},
				"b": parser.JSONFieldObject{
					"c": parser.JSONFieldValue{}, "d": parser.JSONFieldValue{},
					"e": parser.JSONFieldValue{},
				},
				"f": parser.JSONFieldValue{}, "g": parser.JSONFieldValue{},
			},
		},
		"querying complex fields in several places": {
			fields: "a,b(c,d),e(f,g)",
			expected: parser.JSONFieldObject{
				"a": parser.JSONFieldValue{},
				"b": parser.JSONFieldObject{"c": parser.JSONFieldValue{}, "d": parser.JSONFieldValue{}},
				"e": parser.JSONFieldObject{"f": parser.JSONFieldValue{}, "g": parser.JSONFieldValue{}},
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			prs := parser.NewParser()
			actual, err := prs.Parse(test.fields)
			if assert.NoError(t, err) {
				assert.Equal(t, test.expected, actual)
			}
		})
	}
}
