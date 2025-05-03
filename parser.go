package binding

import (
	"strconv"
	"strings"
)

// autoParse tries to infer and parse the type from a string
func autoParse(input string) interface{} {
	input = strings.TrimSpace(input)

	// Check if it's a comma-separated list
	if strings.Contains(input, ",") && strings.HasPrefix(input, "[") && strings.HasSuffix(input, "]") {
		parts := strings.Split(input, ",")
		var parsed []interface{}
		for _, part := range parts {
			parsed = append(parsed, autoParseSingle(strings.TrimSpace(part)))
		}
		return parsed
	}

	return autoParseSingle(input)
}

// autoParseSingle handles individual values (non-array)
func autoParseSingle(input string) interface{} {

	// Try parsing as int
	if i, err := strconv.Atoi(input); err == nil {
		return i
	}

	// Try parsing as float
	if f, err := strconv.ParseFloat(input, 64); err == nil {
		return f
	}

	// Try parsing as bool
	if b, err := strconv.ParseBool(input); err == nil {
		return b
	}

	// Default to string
	return input
}

func flattenMap(input map[string]interface{}, prefix string) map[string]interface{} {
	flatMap := make(map[string]interface{})

	for key, value := range input {
		// Create the new key by combining the prefix and the current key
		newKey := key
		if prefix != "" {
			newKey = prefix + "." + key
		}

		// If the value is a nested map, recurse
		switch v := value.(type) {
		case map[string]interface{}:
			// Recursively flatten the nested map
			for k, val := range flattenMap(v, newKey) {
				flatMap[k] = val
			}
		default:
			// Otherwise, just add the key-value pair to the flat map
			flatMap[newKey] = value
		}
	}

	return flatMap
}
