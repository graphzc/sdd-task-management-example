package servererr

import (
	"github.com/graphzc/sdd-task-management-example/internal/dto"
	"github.com/labstack/echo/v4"
)

func EchoHTTPErrorHandler(err error, c echo.Context) {
	code := ErrorCodeInternalServerError
	message := err.Error()

	// Handle echo HTTP errors
	if echoErr, ok := err.(*echo.HTTPError); ok {
		switch echoErr.Code {
		case 400:
			code = ErrorCodeBadRequest
			message = "Bad request"
		case 401:
			code = ErrorCodeUnauthorized
			message = "Unauthorized"
		case 403:
			code = ErrorCodeForbidden
			message = "Forbidden"
		case 404:
			code = ErrorCodeNotFound
			message = "Route not found"
		}
	} else if serverErr, ok := err.(*ServerError); ok {
		code = serverErr.Code
	}

	c.JSON(code.HTTPStatus(), dto.ErrorResponse{
		Code:    code.String(),
		Message: message,
	})
}
