package configutil

import "encoding/json"

func Format(c interface{}) string {
	bs, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(bs)
}
