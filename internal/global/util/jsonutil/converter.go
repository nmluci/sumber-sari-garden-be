package jsonutil

import "encoding/json"

func ConvertToJSON(data interface{}) (b []byte, err error) {
	b, err = json.Marshal(data)
	return
}
