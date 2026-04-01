package fieldparser_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/manuelarte/milogo/internal/fieldparser"
)

func TestParser(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		fields   string
		expected fieldparser.JSONFieldObject
	}{
		"querying simple fields": {
			fields:   "name,age",
			expected: fieldparser.JSONFieldObject{"name": fieldparser.JSONFieldValue{}, "age": fieldparser.JSONFieldValue{}},
		},
		"querying simple and complex field": {
			fields: "name,address(street,number)",
			expected: fieldparser.JSONFieldObject{
				"name": fieldparser.JSONFieldValue{},
				"address": fieldparser.JSONFieldObject{
					"street": fieldparser.JSONFieldValue{},
					"number": fieldparser.JSONFieldValue{},
				},
			},
		},
		"querying simple and nested complex fields": {
			fields: "a,b(c,d,e(f,g))",
			expected: fieldparser.JSONFieldObject{
				"a": fieldparser.JSONFieldValue{},
				"b": fieldparser.JSONFieldObject{
					"c": fieldparser.JSONFieldValue{}, "d": fieldparser.JSONFieldValue{},
					"e": fieldparser.JSONFieldObject{"f": fieldparser.JSONFieldValue{}, "g": fieldparser.JSONFieldValue{}},
				},
			},
		},
		"querying complex fields in the middle": {
			fields: "a,b(c,d,e),f,g",
			expected: fieldparser.JSONFieldObject{
				"a": fieldparser.JSONFieldValue{},
				"b": fieldparser.JSONFieldObject{
					"c": fieldparser.JSONFieldValue{}, "d": fieldparser.JSONFieldValue{},
					"e": fieldparser.JSONFieldValue{},
				},
				"f": fieldparser.JSONFieldValue{}, "g": fieldparser.JSONFieldValue{},
			},
		},
		"querying complex fields in several places": {
			fields: "a,b(c,d),e(f,g)",
			expected: fieldparser.JSONFieldObject{
				"a": fieldparser.JSONFieldValue{},
				"b": fieldparser.JSONFieldObject{"c": fieldparser.JSONFieldValue{}, "d": fieldparser.JSONFieldValue{}},
				"e": fieldparser.JSONFieldObject{"f": fieldparser.JSONFieldValue{}, "g": fieldparser.JSONFieldValue{}},
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			prs := fieldparser.NewParser()

			actual, err := prs.Parse(test.fields)
			if assert.NoError(t, err) {
				assert.Equal(t, test.expected, actual)
			}
		})
	}
}
