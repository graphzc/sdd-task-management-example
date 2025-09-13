package echoutil

import (
	"context"
	"net/http"
	"reflect"

	"github.com/graphzc/sdd-task-management-example/internal/dto"
	"github.com/graphzc/sdd-task-management-example/internal/utils/servererr"
	"github.com/labstack/echo/v4"
)

// Wrap wraps a business logic function with Echo handler
// Supports both func(req) (res, err) and func() (res, err) signatures
// For functions with request: func(req) (res, err) where req is the request DTO and res is the response DTO
// For functions without request: func() (res, err) where res is the response DTO
func WrapWithStatus[Req any, Res any](fn func(context.Context, Req) (Res, error), status int) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req Req
		reqType := reflect.TypeOf(req)

		// Check if Req is an empty struct (no request needed)
		isEmpty := isEmptyStruct(reqType)

		if !isEmpty {
			// Bind request body to DTO only if request type is not empty
			if err := c.Bind(&req); err != nil {
				return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
					Code:    servererr.ErrorCodeBadRequest.String(),
					Message: "Invalid request format",
				})
			}

			// Validate request if validator is available
			if err := c.Validate(req); err != nil {
				return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
					Code:    servererr.ErrorCodeBadRequest.String(),
					Message: err.Error(),
				})
			}
		}

		// Call business logic function
		res, err := fn(c.Request().Context(), req)
		if err != nil {
			// Handle server errors
			if serverErr, ok := err.(*servererr.ServerError); ok {
				return c.JSON(serverErr.Code.HTTPStatus(), dto.ErrorResponse{
					Code:    serverErr.Code.String(),
					Message: serverErr.Message,
				})
			}

			// Handle generic errors
			return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Code:    servererr.ErrorCodeInternalServerError.String(),
				Message: err.Error(),
			})
		}

		// Return success response
		return c.JSON(status, res)
	}
}

// isEmptyStruct checks if a type is an empty struct
func isEmptyStruct(t reflect.Type) bool {
	if t == nil {
		return true
	}

	// Handle pointer types
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// Check if it's a struct with no fields
	if t.Kind() == reflect.Struct {
		return t.NumField() == 0
	}

	return false
}
