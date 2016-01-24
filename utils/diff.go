package utils

import (
	"net/url"
	"strings"
)

func GetDiffFormValues(formValues url.Values, structValues map[string]interface{}) map[string]interface{} {
	Log.Debug("[DIFF}: struct: %v", structValues)
	vals := make(map[string]interface{})
	for key := range formValues {
		Log.Debug("[DIFF]: key: %v", key)
		val, exists := structValues[strings.Title(key)]
		if exists {
			Log.Debug("[DIFF]: val: %v", val)
			vals[key] = val
		}
	}

	return vals
}
