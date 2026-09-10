package response

import "errors"

type Response struct {
	Code    Code   `json:"code"`
	Success bool   `json:"success"`
	Payload any    `json:"payload,omitempty"`
	Message string `json:"message,omitempty"`
}

func NewResponse(code Code, success bool, payload any, message string) *Response {
	return &Response{Code: code, Success: success, Payload: payload, Message: message}
}

func (this *Response) Error() string {
	return this.Message
}

func IsErrorCode(err error, code Code) bool {
	if err == nil {
		return false
	}
	var se *SafeError
	if errors.As(err, &se) {
		return se.Code == code
	}
	return false
}
