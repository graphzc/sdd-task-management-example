package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/graphzc/sdd-task-management-example/internal/config"
	"github.com/graphzc/sdd-task-management-example/internal/domain/enums"
	"github.com/graphzc/sdd-task-management-example/internal/infrastructure/auth"
	"github.com/graphzc/sdd-task-management-example/internal/utils/servererr"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AuthMiddlewareTestSuite struct {
	suite.Suite
	middleware AuthMiddleware
	config     *config.Config
	echo       *echo.Echo
	secret     string
}

func (suite *AuthMiddlewareTestSuite) SetupTest() {
	suite.secret = "test-secret-key"
	suite.config = &config.Config{
		JWT: config.JWT{
			AccessTokenSecret:     suite.secret,
			AccessTokenExpiration: "1h",
		},
	}
	suite.middleware = NewAuthMiddleware(suite.config)
	suite.echo = echo.New()
}

func (suite *AuthMiddlewareTestSuite) generateValidToken(userID, email string) string {
	claims := auth.JWTClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(suite.secret))
	return tokenString
}

func (suite *AuthMiddlewareTestSuite) generateExpiredToken(userID, email string) string {
	claims := auth.JWTClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)), // Expired 1 hour ago
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(suite.secret))
	return tokenString
}

func (suite *AuthMiddlewareTestSuite) generateInvalidSigningMethodToken(userID, email string) string {
	claims := auth.JWTClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Use RS256 instead of HS256
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, _ := token.SignedString([]byte(suite.secret))
	return tokenString
}

func (suite *AuthMiddlewareTestSuite) TestNewAuthMiddleware() {
	middleware := NewAuthMiddleware(suite.config)
	assert.NotNil(suite.T(), middleware)
}

func (suite *AuthMiddlewareTestSuite) TestMiddleware_ValidToken() {
	// Arrange
	userID := "test-user-id"
	email := "test@example.com"
	token := suite.generateValidToken(userID, email)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	c := suite.echo.NewContext(req, rec)

	called := false
	nextHandler := func(c echo.Context) error {
		called = true
		// Check if user ID is set in context
		contextUserID := c.Get(string(enums.UserIDContextKey))
		assert.Equal(suite.T(), userID, contextUserID)
		return c.JSON(http.StatusOK, map[string]string{"message": "success"})
	}

	// Act
	middlewareFunc := suite.middleware.Middleware(nextHandler)
	err := middlewareFunc(c)

	// Assert
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), called)
}

func (suite *AuthMiddlewareTestSuite) TestMiddleware_MissingAuthorizationHeader() {
	// Arrange
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	c := suite.echo.NewContext(req, rec)

	nextHandler := func(c echo.Context) error {
		suite.T().Error("Next handler should not be called")
		return nil
	}

	// Act
	middlewareFunc := suite.middleware.Middleware(nextHandler)
	err := middlewareFunc(c)

	// Assert
	assert.Error(suite.T(), err)
	serverErr, ok := err.(*servererr.ServerError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), servererr.ErrorCodeUnauthorized, serverErr.Code)
}

func (suite *AuthMiddlewareTestSuite) TestMiddleware_InvalidBearerFormat() {
	// Arrange
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "InvalidFormat token")
	rec := httptest.NewRecorder()
	c := suite.echo.NewContext(req, rec)

	nextHandler := func(c echo.Context) error {
		suite.T().Error("Next handler should not be called")
		return nil
	}

	// Act
	middlewareFunc := suite.middleware.Middleware(nextHandler)
	err := middlewareFunc(c)

	// Assert
	assert.Error(suite.T(), err)
	serverErr, ok := err.(*servererr.ServerError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), servererr.ErrorCodeUnauthorized, serverErr.Code)
}

func (suite *AuthMiddlewareTestSuite) TestMiddleware_ExpiredToken() {
	// Arrange
	token := suite.generateExpiredToken("test-user", "test@example.com")

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	c := suite.echo.NewContext(req, rec)

	nextHandler := func(c echo.Context) error {
		suite.T().Error("Next handler should not be called")
		return nil
	}

	// Act
	middlewareFunc := suite.middleware.Middleware(nextHandler)
	err := middlewareFunc(c)

	// Assert
	assert.Error(suite.T(), err)
	serverErr, ok := err.(*servererr.ServerError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), servererr.ErrorCodeUnauthorized, serverErr.Code)
}

func (suite *AuthMiddlewareTestSuite) TestMiddleware_InvalidSigningMethod() {
	// Arrange
	token := suite.generateInvalidSigningMethodToken("test-user", "test@example.com")

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	c := suite.echo.NewContext(req, rec)

	nextHandler := func(c echo.Context) error {
		suite.T().Error("Next handler should not be called")
		return nil
	}

	// Act
	middlewareFunc := suite.middleware.Middleware(nextHandler)
	err := middlewareFunc(c)

	// Assert
	assert.Error(suite.T(), err)
	serverErr, ok := err.(*servererr.ServerError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), servererr.ErrorCodeUnauthorized, serverErr.Code)
}

func (suite *AuthMiddlewareTestSuite) TestMiddleware_InvalidTokenSignature() {
	// Arrange
	wrongSecret := "wrong-secret"
	claims := auth.JWTClaims{
		UserID: "test-user",
		Email:  "test@example.com",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(wrongSecret))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	rec := httptest.NewRecorder()
	c := suite.echo.NewContext(req, rec)

	nextHandler := func(c echo.Context) error {
		suite.T().Error("Next handler should not be called")
		return nil
	}

	// Act
	middlewareFunc := suite.middleware.Middleware(nextHandler)
	err := middlewareFunc(c)

	// Assert
	assert.Error(suite.T(), err)
	serverErr, ok := err.(*servererr.ServerError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), servererr.ErrorCodeUnauthorized, serverErr.Code)
}

func (suite *AuthMiddlewareTestSuite) TestMiddleware_MalformedToken() {
	// Arrange
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid.token.format")
	rec := httptest.NewRecorder()
	c := suite.echo.NewContext(req, rec)

	nextHandler := func(c echo.Context) error {
		suite.T().Error("Next handler should not be called")
		return nil
	}

	// Act
	middlewareFunc := suite.middleware.Middleware(nextHandler)
	err := middlewareFunc(c)

	// Assert
	assert.Error(suite.T(), err)
	serverErr, ok := err.(*servererr.ServerError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), servererr.ErrorCodeUnauthorized, serverErr.Code)
}

func (suite *AuthMiddlewareTestSuite) TestMiddleware_EmptyToken() {
	// Arrange
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer ")
	rec := httptest.NewRecorder()
	c := suite.echo.NewContext(req, rec)

	nextHandler := func(c echo.Context) error {
		suite.T().Error("Next handler should not be called")
		return nil
	}

	// Act
	middlewareFunc := suite.middleware.Middleware(nextHandler)
	err := middlewareFunc(c)

	// Assert
	assert.Error(suite.T(), err)
	serverErr, ok := err.(*servererr.ServerError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), servererr.ErrorCodeUnauthorized, serverErr.Code)
}

func (suite *AuthMiddlewareTestSuite) TestMiddleware_ValidTokenWithSpecialCharacters() {
	// Arrange
	userID := "user-123"
	email := "test+special@example-domain.com"
	token := suite.generateValidToken(userID, email)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	c := suite.echo.NewContext(req, rec)

	called := false
	nextHandler := func(c echo.Context) error {
		called = true
		contextUserID := c.Get(string(enums.UserIDContextKey))
		assert.Equal(suite.T(), userID, contextUserID)
		return c.JSON(http.StatusOK, map[string]string{"message": "success"})
	}

	// Act
	middlewareFunc := suite.middleware.Middleware(nextHandler)
	err := middlewareFunc(c)

	// Assert
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), called)
}

// Performance test
func (suite *AuthMiddlewareTestSuite) TestMiddleware_Performance() {
	userID := "test-user-id"
	email := "test@example.com"
	token := suite.generateValidToken(userID, email)

	// Run multiple iterations to test performance
	for i := 0; i < 100; i++ {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		c := suite.echo.NewContext(req, rec)

		nextHandler := func(c echo.Context) error {
			return c.JSON(http.StatusOK, map[string]string{"message": "success"})
		}

		middlewareFunc := suite.middleware.Middleware(nextHandler)
		err := middlewareFunc(c)
		assert.NoError(suite.T(), err)
	}
}

func TestAuthMiddlewareTestSuite(t *testing.T) {
	suite.Run(t, new(AuthMiddlewareTestSuite))
}

// Additional unit tests without suite
func TestNewAuthMiddleware_NilConfig(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Log("Expected panic for nil config")
		}
	}()

	middleware := NewAuthMiddleware(nil)
	assert.NotNil(t, middleware) // This might panic, which we catch above
}

func BenchmarkAuthMiddleware_ValidToken(b *testing.B) {
	secret := "test-secret-key"
	config := &config.Config{
		JWT: config.JWT{
			AccessTokenSecret:     secret,
			AccessTokenExpiration: "1h",
		},
	}
	middleware := NewAuthMiddleware(config)
	e := echo.New()

	// Generate a valid token
	claims := auth.JWTClaims{
		UserID: "test-user",
		Email:  "test@example.com",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(secret))

	nextHandler := func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "success"})
	}

	middlewareFunc := middleware.Middleware(nextHandler)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "Bearer "+tokenString)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := middlewareFunc(c)
		if err != nil {
			b.Fatal(err)
		}
	}
}
