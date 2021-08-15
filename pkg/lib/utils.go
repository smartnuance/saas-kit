package lib

import "encoding/json"

func PP(v interface{}) string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		return string(b)
	}
	return err.Error()
}
