package middlewares

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/graphzc/sdd-task-management-example/internal/config"
	"github.com/graphzc/sdd-task-management-example/internal/infrastructure/auth"
	"github.com/graphzc/sdd-task-management-example/internal/utils/servererr"
	"github.com/graphzc/sdd-task-management-example/internal/utils/tokenutil"
	"github.com/labstack/echo/v4"
)

type authMiddleware struct {
	configs *config.Config
}

type AuthMiddleware interface {
	Middleware(next echo.HandlerFunc) echo.HandlerFunc
}

// @WireSet("Middleware")
func NewAuthMiddleware(configs *config.Config) AuthMiddleware {
	return &authMiddleware{
		configs: configs,
	}
}

func (a *authMiddleware) Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString, err := tokenutil.GetTokenFromEchoHeader(c)
		if err != nil {
			return servererr.NewError(
				servererr.ErrorCodeUnauthorized,
				err.Error(),
			)
		}

		// Parse and validate the token
		token, err := jwt.ParseWithClaims(tokenString, &auth.JWTClaims{}, func(token *jwt.Token) (any, error) {
			// Validate the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(a.configs.JWT.AccessTokenSecret), nil
		})
		if err != nil {
			return servererr.NewError(
				servererr.ErrorCodeUnauthorized,
				err.Error(),
			)
		}

		// Validate claims
		claims, ok := token.Claims.(*auth.JWTClaims)
		if !ok || !token.Valid {
			return servererr.NewError(
				servererr.ErrorCodeUnauthorized,
				"invalid token claims",
			)
		}

		// Set claims and user ID to context
		c.Set("profile", claims)
		c.Set("user_id", claims.UserID)

		return next(c)
	}
}
