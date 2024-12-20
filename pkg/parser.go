package pkg

const (
	defaultFieldSeparator = ","
)

type JsonField interface{}

var _ JsonField = new(JsonFieldValue)
var _ JsonField = new(JsonFieldObject)

// JsonFieldValue Holder that indicates that the field is just a value.
type JsonFieldValue struct{}

type JsonFieldObject map[string]JsonField

func NewParser() Parser {
	return Parser{
		fieldSeparator: defaultFieldSeparator,
	}
}

type Parser struct {
	fieldSeparator string
}

// Parse String to get what fields are present.
func (p Parser) Parse(fields string) (JsonFieldObject, error) {
	if fields == "" {
		return nil, ErrFieldsIsEmpty
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

func (p Parser) parseChunk(chunk string, index *int, openParenthesis *int) (JsonFieldObject, error) {
	if chunk == "" {
		return nil, ErrFieldIsEmpty
	}
	toReturn := JsonFieldObject{}
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
			newJsonField, err := p.parseChunk(chunk, index, openParenthesis)
			if err != nil {
				return toReturn, err
			}
			toReturn[field] = newJsonField
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
		return nil, ErrUnbalancedParenthesis
	}

	return toReturn, nil
}

func (p Parser) addFieldValue(field string, object JsonFieldObject) error {
	if field == "" {
		return ErrFieldIsEmpty
	}
	object[field] = JsonFieldValue{}

	return nil
}
