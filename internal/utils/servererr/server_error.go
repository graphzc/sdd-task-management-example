package servererr

type ServerError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

func NewError(code ErrorCode, message string) *ServerError {
	return &ServerError{
		Code:    code,
		Message: message,
	}
}

func (e *ServerError) Error() string {
	return e.Message
}
