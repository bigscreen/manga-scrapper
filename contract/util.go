package contract

import "encoding/json"

func getJson(i interface{}) string {
	j, err := json.Marshal(i)
	if err == nil {
		return string(j)
	}
	return ""
}
