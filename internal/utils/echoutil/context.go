package echoutil

import (
	"context"
	"errors"

	"github.com/graphzc/sdd-task-management-example/internal/domain/constants"
	"github.com/labstack/echo/v4"
)

// GetUserIDFromContext extracts user ID from Echo context
func GetUserIDFromEchoContext(c echo.Context) (string, error) {
	userID := c.Get(string(constants.UserIDContextKey))
	if userID == nil {
		return "", errors.New("user ID not found in context")
	}

	userIDStr, ok := userID.(string)
	if !ok {
		return "", errors.New("invalid user ID type in context")
	}

	return userIDStr, nil
}

// SetUserIDInContext sets user ID in standard context
func SetUserIDInContext(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, constants.UserIDContextKey, userID)
}

// GetUserIDFromContext extracts user ID from standard context
func GetUserIDFromContext(ctx context.Context) (string, error) {
	userID := ctx.Value(constants.UserIDContextKey)
	if userID == nil {
		return "", errors.New("user ID not found in context")
	}

	userIDStr, ok := userID.(string)
	if !ok {
		return "", errors.New("invalid user ID type in context")
	}

	return userIDStr, nil
}
