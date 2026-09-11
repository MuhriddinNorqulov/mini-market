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
	if se, ok := errors.AsType[*SafeError](err); ok {
		return se.Code == code
	}
	if r, ok := errors.AsType[*Response](err); ok {
		return r.Code == code
	}
	return false
}
