package ctx

import (
	"fmt"
	"strconv"
)

func GetBody[T any](c Context) (*T, error) {
	var data = new(T)
	if err := c.Bind(data); err != nil {
		return nil, err
	}
	if err := c.Validate(data); err != nil {
		return nil, err
	}
	return data, nil
}

func GetIntQueryParam(c Context, param string, defaultValue int) int {
	value := c.QueryParam(param)
	if value == "" {
		return defaultValue
	}

	v, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return v
}

func GetBoolQueryParam(c Context, param string) *bool {
	switch c.QueryParam(param) {
	case "true":
		return new(true)
	case "false":
		return new(false)
	default:
		return nil
	}
}

func GetIntPathParam(c Context, param string) (int, error) {
	value := c.Param(param)
	if value == "" {
		return 0, fmt.Errorf("%s is required", param)
	}

	id, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("invalid %s: %v", param, err)
	}
	return id, nil
}
