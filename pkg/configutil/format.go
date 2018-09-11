package configutil

import "encoding/json"

// Format the given configuration c into formatted JSON
// with 2-space indentation level.
func Format(c interface{}) string {
	bs, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(bs)
}
