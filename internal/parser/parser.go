package parser

import "github.com/manuelarte/milogo/pkg/errors"

type JSONField interface{}

var (
	_ JSONField = new(JSONFieldValue)
	_ JSONField = new(JSONFieldObject)
)

// JSONFieldValue Holder that indicates that the field is just a value.
type JSONFieldValue struct{}

type JSONFieldObject map[string]JSONField

func NewParser() Parser {
	const defaultFieldSeparator = ","
	return Parser{
		fieldSeparator: defaultFieldSeparator,
	}
}

type Parser struct {
	fieldSeparator string
}

// Parse String to get what fields are present.
func (p Parser) Parse(fields string) (JSONFieldObject, error) {
	if fields == "" {
		return nil, errors.ErrFieldsIsEmpty
	}
	openParenthesis := 0
	index := 0
	toReturn, err := p.parseChunk(fields, &index, &openParenthesis)
	if err != nil {
		return nil, err
	}
	// check wrong number of parenthesis

	return toReturn, nil
}

//nolint:gocognit
func (p Parser) parseChunk(chunk string, index *int, openParenthesis *int) (JSONFieldObject, error) {
	if chunk == "" {
		return nil, errors.ErrFieldIsEmpty
	}
	toReturn := JSONFieldObject{}
	field := ""
	for *index < len(chunk) {
		char := string(chunk[*index])
		switch char {
		case p.fieldSeparator:
			if field != "" {
				err := p.addFieldValue(field, toReturn)
				if err != nil {
					return toReturn, err
				}
				field = ""
			}
		case "(":
			*openParenthesis++
			*index++
			newJSONField, err := p.parseChunk(chunk, index, openParenthesis)
			if err != nil {
				return toReturn, err
			}
			toReturn[field] = newJSONField
			field = ""
		case ")":
			*openParenthesis--
			if field != "" {
				err := p.addFieldValue(field, toReturn)
				if err != nil {
					return toReturn, err
				}
			}

			return toReturn, nil
		default:
			field += char
		}
		*index++
	}
	if len(field) > 0 {
		err := p.addFieldValue(field, toReturn)
		if err != nil {
			return toReturn, err
		}
	}
	if *openParenthesis != 0 {
		return nil, errors.ErrUnbalancedParenthesis
	}

	return toReturn, nil
}

func (p Parser) addFieldValue(field string, object JSONFieldObject) error {
	if field == "" {
		return errors.ErrFieldIsEmpty
	}
	object[field] = JSONFieldValue{}

	return nil
}
