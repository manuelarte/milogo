package parser

import (
	"github.com/manuelarte/milogo/pkg/errors"
)

func Filter(jsonData any, partialResponseFields JSONFieldObject) error {
	if casted, ok := jsonData.(map[string]any); ok {
		return filterMap(casted, partialResponseFields)
	} else if array, okCast := jsonData.([]any); okCast {
		for _, item := range array {
			if innerCasted, okMap := item.(map[string]any); okMap {
				if err := filterMap(innerCasted, partialResponseFields); err != nil {
					return err
				}
			}
		}

		return nil
	}

	return errors.ErrUnrecognizedFormat
}

//nolint:gocognit // Refactor later
func filterMap(jsonData map[string]any, partialResponseFields JSONFieldObject) error {
	for key, value := range jsonData {
		//nolint:nestif // Refactor later
		if _, ok := partialResponseFields[key]; !ok {
			delete(jsonData, key)
		} else {
			if values, okCast := value.([]map[string]any); okCast {
				for _, value := range values {
					if nestedPartialResponse, isFieldObject := partialResponseFields[key].(JSONFieldObject); isFieldObject {
						return filterMap(value, nestedPartialResponse)
					}
				}
			} else {
				if casted, isFieldObject := partialResponseFields[key].(JSONFieldObject); isFieldObject {
					if nestedObject, isMap := value.(map[string]any); isMap {
						return filterMap(nestedObject, casted)
					}
					return errors.NotAnObjectError(key)
				}
			}
		}
	}

	return nil
}
