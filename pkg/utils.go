package pkg

func Filter(jsonData interface{}, partialResponseFields JsonFieldObject) error {
	if casted, ok := jsonData.(map[string]interface{}); ok {
		return filterMap(casted, partialResponseFields)
	} else if array, ok := jsonData.([]interface{}); ok {
		for _, item := range array {
			if casted, ok := item.(map[string]interface{}); ok {
				if err := filterMap(casted, partialResponseFields); err != nil {
					return err
				}
			}
		}

		return nil
	}

	return ErrUnrecognizedFormat
}

func filterMap(jsonData map[string]interface{}, partialResponseFields JsonFieldObject) error {
	for key, value := range jsonData {
		if _, ok := partialResponseFields[key]; !ok {
			delete(jsonData, key)
		} else {
			if values, ok := value.([]map[string]interface{}); ok {
				for _, value := range values {
					if nestedPartialResponse, ok := partialResponseFields[key].(JsonFieldObject); ok {
						return filterMap(value, nestedPartialResponse)
					}
				}
			} else {
				if casted, ok := partialResponseFields[key].(JsonFieldObject); ok {
					if nestedObject, ok := value.(map[string]interface{}); ok {
						return filterMap(nestedObject, casted)
					} else {
						return NotAnObjectError(key)
					}
				}
			}
		}
	}

	return nil
}
