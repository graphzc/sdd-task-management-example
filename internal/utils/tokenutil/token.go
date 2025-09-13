package tokenutil

import (
	"errors"
	"strings"

	"github.com/graphzc/sdd-task-management-example/internal/utils/servererr"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

var (
	ErrInvalidToken = errors.New("invalid token")
)

func SplitBearerToken(bearer string) (string, error) {
	bearer = strings.TrimSpace(bearer)
	splittedToken := strings.Split(bearer, "Bearer ")
	if len(splittedToken) != 2 {
		return "", ErrInvalidToken
	}

	token := splittedToken[1]

	return token, nil
}

func GetTokenFromEchoHeader(c echo.Context) (string, error) {
	bearer := c.Request().Header.Get("Authorization")

	bearer = strings.TrimSpace(bearer)
	token, err := SplitBearerToken(bearer)
	if err != nil {
		return "", err
	}

	return token, nil
}

func GetProfileOnEchoContext[T any](c echo.Context, key string) (*T, error) {
	profile := c.Get(key)

	if profile == nil {
		log.Ctx(c.Request().Context()).
			Error().Msg("No user found in context")

		return nil, servererr.NewError(
			servererr.ErrorCodeUnauthorized,
			"Invalid user",
		)
	}

	userClaims, ok := profile.(*T)
	if !ok {
		log.Ctx(c.Request().Context()).
			Error().Msg("Invalid user claims type in context")

		return nil, servererr.NewError(
			servererr.ErrorCodeInternalServerError,
			"Invalid user claims type",
		)
	}

	return userClaims, nil
}
