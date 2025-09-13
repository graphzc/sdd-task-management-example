package servererr

type ErrorCode string

const (
	ErrorCodeInternalServerError ErrorCode = "INTERNAL_SERVER_ERROR"
	ErrorCodeBadRequest          ErrorCode = "BAD_REQUEST"
	ErrorCodeNotFound            ErrorCode = "NOT_FOUND"
	ErrorCodeUnauthorized        ErrorCode = "UNAUTHORIZED"
	ErrorCodeForbidden           ErrorCode = "FORBIDDEN"
	ErrorCodeConflict            ErrorCode = "CONFLICT"
	ErrorCodeTooManyRequests     ErrorCode = "TOO_MANY_REQUESTS"
	ErrorCodeServiceUnavailable  ErrorCode = "SERVICE_UNAVAILABLE"
)

func (e ErrorCode) String() string {
	return string(e)
}

func (e ErrorCode) HTTPStatus() int {
	mapErrorToHTTPStatus := map[ErrorCode]int{
		ErrorCodeInternalServerError: 500,
		ErrorCodeBadRequest:          400,
		ErrorCodeNotFound:            404,
		ErrorCodeUnauthorized:        401,
		ErrorCodeForbidden:           403,
		ErrorCodeConflict:            409,
		ErrorCodeTooManyRequests:     429,
		ErrorCodeServiceUnavailable:  503,
	}

	if status, ok := mapErrorToHTTPStatus[e]; ok {
		return status
	}
	return 500
}
