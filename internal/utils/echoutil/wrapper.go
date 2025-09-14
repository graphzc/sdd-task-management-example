package echoutil

import (
	"context"
	"net/http"
	"reflect"
	"time"

	"github.com/graphzc/sdd-task-management-example/internal/dto"
	"github.com/graphzc/sdd-task-management-example/internal/utils/servererr"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

// WrapWithStatus wraps a business logic function with Echo handler
// Supports both func(req) (res, err) and func() (res, err) signatures
// For functions with request: func(req) (res, err) where req is the request DTO and res is the response DTO
// For functions without request: func() (res, err) where res is the response DTO
func WrapWithStatus[Req any, Res any](fn func(context.Context, Req) (Res, error), status int) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()

		// Log incoming request
		log.Info().
			Str("method", c.Request().Method).
			Str("path", c.Request().URL.Path).
			Str("remote_addr", c.RealIP()).
			Str("user_agent", c.Request().UserAgent()).
			Msg("Incoming request")

		var req Req
		reqType := reflect.TypeOf(req)

		// Check if Req is an empty struct (no request needed)
		isEmpty := isEmptyStruct(reqType)

		if !isEmpty {
			// Bind request body to DTO only if request type is not empty
			if err := c.Bind(&req); err != nil {
				log.Error().
					Str("method", c.Request().Method).
					Str("path", c.Request().URL.Path).
					Str("remote_addr", c.RealIP()).
					Str("user_agent", c.Request().UserAgent()).
					Msg("Failed to bind request")

				return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
					Code:    servererr.ErrorCodeBadRequest.String(),
					Message: "Invalid request format",
				})
			}

			// Validate request if validator is available
			if err := c.Validate(req); err != nil {
				log.Error().
					Str("method", c.Request().Method).
					Str("path", c.Request().URL.Path).
					Str("remote_addr", c.RealIP()).
					Str("user_agent", c.Request().UserAgent()).
					Msg("Request validation failed")

				return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
					Code:    servererr.ErrorCodeBadRequest.String(),
					Message: err.Error(),
				})
			}
		} else {
			// Even for empty requests, we might need to bind path parameters
			if err := c.Bind(&req); err != nil {
				// Ignore binding errors for truly empty structs
				if !isEmptyStruct(reqType) {
					log.Error().
						Str("method", c.Request().Method).
						Str("path", c.Request().URL.Path).
						Str("remote_addr", c.RealIP()).
						Str("user_agent", c.Request().UserAgent()).
						Msg("Failed to bind request")

					return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
						Code:    servererr.ErrorCodeBadRequest.String(),
						Message: "Invalid request format",
					})
				}
			}
		}

		// Create context with user ID if available
		ctx := c.Request().Context()
		userID, err := GetUserIDFromEchoContext(c)
		if err == nil && userID != "" {
			ctx = SetUserIDInContext(ctx, userID)
		}

		// Call business logic function
		res, err := fn(ctx, req)

		duration := time.Since(start)

		if err != nil {
			// Log error response
			log.Error().
				Str("method", c.Request().Method).
				Str("path", c.Request().URL.Path).
				Dur("duration", duration).
				Err(err).
				Msg("Request completed with error")

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

		// Log successful response
		log.Info().
			Str("method", c.Request().Method).
			Str("path", c.Request().URL.Path).
			Dur("duration", duration).
			Int("status", status).
			Msg("Request completed successfully")

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
