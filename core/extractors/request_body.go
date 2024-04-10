package extractors

import (
	"fmt"

	"github.com/tidwall/gjson"
)

type bodyExtractor func(params ...interface{}) (interface{}, error)

// RequestBody extracts data from the request body
// RequestBody is a HOF that returns a function that extracts data from the request body
// The function returned by RequestBody can be used in the template engine to extract data from the request body
// The function can be used in the template engine with the following signatures:
// {{ requestBody }} -> returns the whole request body
// {{ requestBody "key" }} -> returns the value of the key in the request body
// {{ requestBody "key" "fallback" }} -> returns the value of the key in the request body
func RequestBody(reqBody *string) bodyExtractor {
	return func(params ...interface{}) (interface{}, error) {
		if len(params) == 0 {
			return *reqBody, nil
		}
		var key string
		var fallback interface{}
		key, ok := params[0].(string)
		if !ok {
			return nil, fmt.Errorf("requestBody: invalid key type")
		}
		if len(params) > 1 {
			fallback = params[1]
		}
		res := gjson.Get(*reqBody, key)
		if !res.Exists() {
			if fallback != nil {
				return fallback, nil
			}
			return "", nil
		}
		return res.Raw, nil
	}
}
