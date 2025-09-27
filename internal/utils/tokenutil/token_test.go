package tokenutil

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TokenUtilTestSuite struct {
	suite.Suite
	echo *echo.Echo
}

func (suite *TokenUtilTestSuite) SetupTest() {
	suite.echo = echo.New()
}

// Test SplitBearerToken function
func (suite *TokenUtilTestSuite) TestSplitBearerToken_ValidToken() {
	// Arrange
	bearer := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.test-token"
	expectedToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.test-token"

	// Act
	token, err := SplitBearerToken(bearer)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedToken, token)
}

func (suite *TokenUtilTestSuite) TestSplitBearerToken_ValidTokenWithExtraSpaces() {
	// Arrange
	bearer := "  Bearer   eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.test-token  "
	expectedToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.test-token"

	// Act
	token, err := SplitBearerToken(bearer)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedToken, token)
}

func (suite *TokenUtilTestSuite) TestSplitBearerToken_EmptyString() {
	// Arrange
	bearer := ""

	// Act
	token, err := SplitBearerToken(bearer)

	// Assert
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), ErrInvalidToken, err)
	assert.Empty(suite.T(), token)
}

func (suite *TokenUtilTestSuite) TestSplitBearerToken_InvalidFormat_NoBearer() {
	// Arrange
	bearer := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.test-token"

	// Act
	token, err := SplitBearerToken(bearer)

	// Assert
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), ErrInvalidToken, err)
	assert.Empty(suite.T(), token)
}

func (suite *TokenUtilTestSuite) TestSplitBearerToken_InvalidFormat_WrongPrefix() {
	// Arrange
	bearer := "Basic eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.test-token"

	// Act
	token, err := SplitBearerToken(bearer)

	// Assert
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), ErrInvalidToken, err)
	assert.Empty(suite.T(), token)
}

func (suite *TokenUtilTestSuite) TestSplitBearerToken_InvalidFormat_MultipleBearers() {
	// Arrange
	bearer := "Bearer token1 Bearer token2"

	// Act
	token, err := SplitBearerToken(bearer)

	// Assert
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), ErrInvalidToken, err)
	assert.Empty(suite.T(), token)
}

func (suite *TokenUtilTestSuite) TestSplitBearerToken_OnlyBearer() {
	// Arrange
	bearer := "Bearer "

	// Act
	token, err := SplitBearerToken(bearer)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Empty(suite.T(), token)
}

func (suite *TokenUtilTestSuite) TestSplitBearerToken_OnlyBearerWithSpaces() {
	// Arrange
	bearer := "  Bearer   "

	// Act
	token, err := SplitBearerToken(bearer)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Empty(suite.T(), token)
}

func (suite *TokenUtilTestSuite) TestSplitBearerToken_CaseSensitive() {
	// Arrange
	bearer := "bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.test-token"

	// Act
	token, err := SplitBearerToken(bearer)

	// Assert
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), ErrInvalidToken, err)
	assert.Empty(suite.T(), token)
}

// Test GetTokenFromEchoHeader function
func (suite *TokenUtilTestSuite) TestGetTokenFromEchoHeader_ValidAuthorizationHeader() {
	// Arrange
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.test-token"
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	c := suite.echo.NewContext(req, rec)

	// Act
	result, err := GetTokenFromEchoHeader(c)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), token, result)
}

func (suite *TokenUtilTestSuite) TestGetTokenFromEchoHeader_ValidAuthorizationHeaderWithSpaces() {
	// Arrange
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.test-token"
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "  Bearer  "+token+"  ")
	rec := httptest.NewRecorder()
	c := suite.echo.NewContext(req, rec)

	// Act
	result, err := GetTokenFromEchoHeader(c)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), token+"  ", result) // Note: SplitBearerToken doesn't trim the token part
}

func (suite *TokenUtilTestSuite) TestGetTokenFromEchoHeader_MissingAuthorizationHeader() {
	// Arrange
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	c := suite.echo.NewContext(req, rec)

	// Act
	result, err := GetTokenFromEchoHeader(c)

	// Assert
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), ErrInvalidToken, err)
	assert.Empty(suite.T(), result)
}

func (suite *TokenUtilTestSuite) TestGetTokenFromEchoHeader_EmptyAuthorizationHeader() {
	// Arrange
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "")
	rec := httptest.NewRecorder()
	c := suite.echo.NewContext(req, rec)

	// Act
	result, err := GetTokenFromEchoHeader(c)

	// Assert
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), ErrInvalidToken, err)
	assert.Empty(suite.T(), result)
}

func (suite *TokenUtilTestSuite) TestGetTokenFromEchoHeader_InvalidBearerFormat() {
	// Arrange
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Basic dGVzdA==")
	rec := httptest.NewRecorder()
	c := suite.echo.NewContext(req, rec)

	// Act
	result, err := GetTokenFromEchoHeader(c)

	// Assert
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), ErrInvalidToken, err)
	assert.Empty(suite.T(), result)
}

func (suite *TokenUtilTestSuite) TestGetTokenFromEchoHeader_OnlyBearerNoToken() {
	// Arrange
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer")
	rec := httptest.NewRecorder()
	c := suite.echo.NewContext(req, rec)

	// Act
	result, err := GetTokenFromEchoHeader(c)

	// Assert
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), ErrInvalidToken, err)
	assert.Empty(suite.T(), result)
}

func (suite *TokenUtilTestSuite) TestGetTokenFromEchoHeader_OnlySpaces() {
	// Arrange
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "   ")
	rec := httptest.NewRecorder()
	c := suite.echo.NewContext(req, rec)

	// Act
	result, err := GetTokenFromEchoHeader(c)

	// Assert
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), ErrInvalidToken, err)
	assert.Empty(suite.T(), result)
}

// Edge case tests
func (suite *TokenUtilTestSuite) TestGetTokenFromEchoHeader_MultipleBearerWords() {
	// Arrange
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer Bearer token")
	rec := httptest.NewRecorder()
	c := suite.echo.NewContext(req, rec)

	// Act
	result, err := GetTokenFromEchoHeader(c)

	// Assert
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), ErrInvalidToken, err)
	assert.Empty(suite.T(), result)
}

func (suite *TokenUtilTestSuite) TestGetTokenFromEchoHeader_LongValidToken() {
	// Arrange
	longToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+longToken)
	rec := httptest.NewRecorder()
	c := suite.echo.NewContext(req, rec)

	// Act
	result, err := GetTokenFromEchoHeader(c)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), longToken, result)
}

func TestTokenUtilTestSuite(t *testing.T) {
	suite.Run(t, new(TokenUtilTestSuite))
}

// Additional unit tests for error constants and edge cases
func TestErrInvalidToken_ErrorMessage(t *testing.T) {
	assert.Equal(t, "invalid token", ErrInvalidToken.Error())
}

// Benchmark tests
func BenchmarkSplitBearerToken(b *testing.B) {
	bearer := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := SplitBearerToken(bearer)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGetTokenFromEchoHeader(b *testing.B) {
	e := echo.New()
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		_, err := GetTokenFromEchoHeader(c)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Table-driven tests
func TestSplitBearerToken_TableDriven(t *testing.T) {
	testCases := []struct {
		name        string
		input       string
		expectedOut string
		expectError bool
	}{
		{
			name:        "Valid Bearer token",
			input:       "Bearer valid-token-123",
			expectedOut: "valid-token-123",
			expectError: false,
		},
		{
			name:        "Bearer with extra spaces",
			input:       "  Bearer   spaced-token  ",
			expectedOut: "spaced-token  ",
			expectError: false,
		},
		{
			name:        "Empty input",
			input:       "",
			expectedOut: "",
			expectError: true,
		},
		{
			name:        "No Bearer prefix",
			input:       "just-a-token",
			expectedOut: "",
			expectError: true,
		},
		{
			name:        "Wrong case",
			input:       "bearer lowercase-token",
			expectedOut: "",
			expectError: true,
		},
		{
			name:        "Only Bearer",
			input:       "Bearer",
			expectedOut: "",
			expectError: true,
		},
		{
			name:        "Multiple Bearer words",
			input:       "Bearer Bearer double-bearer",
			expectedOut: "",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := SplitBearerToken(tc.input)

			if tc.expectError {
				assert.Error(t, err)
				assert.Equal(t, ErrInvalidToken, err)
				assert.Empty(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedOut, result)
			}
		})
	}
}
