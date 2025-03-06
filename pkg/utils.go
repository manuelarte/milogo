package pkg

func Filter(jsonData interface{}, partialResponseFields JSONFieldObject) error {
	if casted, ok := jsonData.(map[string]interface{}); ok {
		return filterMap(casted, partialResponseFields)
	} else if array, okCast := jsonData.([]interface{}); okCast {
		for _, item := range array {
			if innerCasted, okMap := item.(map[string]interface{}); okMap {
				if err := filterMap(innerCasted, partialResponseFields); err != nil {
					return err
				}
			}
		}

		return nil
	}

	return ErrUnrecognizedFormat
}

func filterMap(jsonData map[string]interface{}, partialResponseFields JSONFieldObject) error {
	for key, value := range jsonData {
		if _, ok := partialResponseFields[key]; !ok {
			delete(jsonData, key)
		} else {
			if values, okCast := value.([]map[string]interface{}); okCast {
				for _, value := range values {
					if nestedPartialResponse, isFieldObject := partialResponseFields[key].(JSONFieldObject); isFieldObject {
						return filterMap(value, nestedPartialResponse)
					}
				}
			} else {
				if casted, isFieldObject := partialResponseFields[key].(JSONFieldObject); isFieldObject {
					if nestedObject, isMap := value.(map[string]interface{}); isMap {
						return filterMap(nestedObject, casted)
					}
					return NotAnObjectError(key)
				}
			}
		}
	}

	return nil
}
