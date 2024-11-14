package pkg

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParser(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		fields   string
		expected JsonFieldObject
	}{
		"querying simple fields": {
			fields:   "name,age",
			expected: JsonFieldObject{"name": JsonFieldValue{}, "age": JsonFieldValue{}},
		},
		"querying simple and complex field": {
			fields:   "name,address(street,number)",
			expected: JsonFieldObject{"name": JsonFieldValue{}, "address": JsonFieldObject{"street": JsonFieldValue{}, "number": JsonFieldValue{}}},
		},
		"querying simple and nested complex fields": {
			fields: "a,b(c,d,e(f,g))",
			expected: JsonFieldObject{"a": JsonFieldValue{},
				"b": JsonFieldObject{"c": JsonFieldValue{}, "d": JsonFieldValue{},
					"e": JsonFieldObject{"f": JsonFieldValue{}, "g": JsonFieldValue{}}}},
		},
		"querying complex fields in the middle": {
			fields: "a,b(c,d,e),f,g",
			expected: JsonFieldObject{"a": JsonFieldValue{},
				"b": JsonFieldObject{"c": JsonFieldValue{}, "d": JsonFieldValue{}, "e": JsonFieldValue{}},
				"f": JsonFieldValue{}, "g": JsonFieldValue{}},
		},
		"querying complex fields in several places": {
			fields: "a,b(c,d),e(f,g)",
			expected: JsonFieldObject{"a": JsonFieldValue{},
				"b": JsonFieldObject{"c": JsonFieldValue{}, "d": JsonFieldValue{}},
				"e": JsonFieldObject{"f": JsonFieldValue{}, "g": JsonFieldValue{}}},
		},
	}
	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			parser := NewParser()
			actual, err := parser.Parse(test.fields)
			if assert.NoError(t, err) {
				assert.Equal(t, test.expected, actual)
			}
		})
	}
}
