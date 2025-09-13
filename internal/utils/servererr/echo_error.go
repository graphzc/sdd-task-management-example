package servererr

import (
	"github.com/graphzc/sdd-task-management-example/internal/dto"
	"github.com/labstack/echo/v4"
)

func EchoHTTPErrorHandler(err error, c echo.Context) {
	code := ErrorCodeInternalServerError
	message := err.Error()

	if serverErr, ok := err.(*ServerError); ok {
		code = serverErr.Code
	}

	c.JSON(code.HTTPStatus(), dto.ErrorResponse{
		Code:    code.String(),
		Message: message,
	})
}
