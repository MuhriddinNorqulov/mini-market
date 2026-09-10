package utils

import "encoding/json"

func JsonUnmarshal[T any](data []byte) (T, error) {
	var t T
	if err := json.Unmarshal(data, &t); err != nil {
		var zero T
		return zero, err
	}
	return t, nil
}

func JsonMarshal[T any](t T) ([]byte, error) {
	return json.Marshal(t)
}
