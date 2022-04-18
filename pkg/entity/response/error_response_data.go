package response

import (
	"encoding/json"
	"io"
)

type ErrorResponseData []ErrorResponseValue

type ErrorResponseValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (er *ErrorResponseData) ToJSON(w io.Writer) error {
	return json.NewEncoder(w).Encode(er)
}

func NewErrorResponseValue(key string, value string) ErrorResponseValue {
	return ErrorResponseValue{Key: key, Value: value}
}

func NewErrorResponseData(errorResponseValues ...ErrorResponseValue) ErrorResponseData {
	errors := ErrorResponseData{}

	for _, v := range errorResponseValues {
		errors = append(errors, v)
	}

	return errors
}