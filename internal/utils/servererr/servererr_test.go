package servererr

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ServerErrorTestSuite struct {
	suite.Suite
}

// Test ErrorCode constants
func (suite *ServerErrorTestSuite) TestErrorCode_Constants() {
	assert.Equal(suite.T(), ErrorCode("INTERNAL_SERVER_ERROR"), ErrorCodeInternalServerError)
	assert.Equal(suite.T(), ErrorCode("BAD_REQUEST"), ErrorCodeBadRequest)
	assert.Equal(suite.T(), ErrorCode("NOT_FOUND"), ErrorCodeNotFound)
	assert.Equal(suite.T(), ErrorCode("UNAUTHORIZED"), ErrorCodeUnauthorized)
	assert.Equal(suite.T(), ErrorCode("FORBIDDEN"), ErrorCodeForbidden)
	assert.Equal(suite.T(), ErrorCode("CONFLICT"), ErrorCodeConflict)
	assert.Equal(suite.T(), ErrorCode("TOO_MANY_REQUESTS"), ErrorCodeTooManyRequests)
	assert.Equal(suite.T(), ErrorCode("SERVICE_UNAVAILABLE"), ErrorCodeServiceUnavailable)
}

// Test ErrorCode.String() method
func (suite *ServerErrorTestSuite) TestErrorCode_String() {
	assert.Equal(suite.T(), "INTERNAL_SERVER_ERROR", ErrorCodeInternalServerError.String())
	assert.Equal(suite.T(), "BAD_REQUEST", ErrorCodeBadRequest.String())
	assert.Equal(suite.T(), "NOT_FOUND", ErrorCodeNotFound.String())
	assert.Equal(suite.T(), "UNAUTHORIZED", ErrorCodeUnauthorized.String())
	assert.Equal(suite.T(), "FORBIDDEN", ErrorCodeForbidden.String())
	assert.Equal(suite.T(), "CONFLICT", ErrorCodeConflict.String())
	assert.Equal(suite.T(), "TOO_MANY_REQUESTS", ErrorCodeTooManyRequests.String())
	assert.Equal(suite.T(), "SERVICE_UNAVAILABLE", ErrorCodeServiceUnavailable.String())
}

// Test ErrorCode.HTTPStatus() method
func (suite *ServerErrorTestSuite) TestErrorCode_HTTPStatus() {
	// Test known error codes
	assert.Equal(suite.T(), 500, ErrorCodeInternalServerError.HTTPStatus())
	assert.Equal(suite.T(), 400, ErrorCodeBadRequest.HTTPStatus())
	assert.Equal(suite.T(), 404, ErrorCodeNotFound.HTTPStatus())
	assert.Equal(suite.T(), 401, ErrorCodeUnauthorized.HTTPStatus())
	assert.Equal(suite.T(), 403, ErrorCodeForbidden.HTTPStatus())
	assert.Equal(suite.T(), 409, ErrorCodeConflict.HTTPStatus())
	assert.Equal(suite.T(), 429, ErrorCodeTooManyRequests.HTTPStatus())
	assert.Equal(suite.T(), 503, ErrorCodeServiceUnavailable.HTTPStatus())
}

func (suite *ServerErrorTestSuite) TestErrorCode_HTTPStatus_UnknownCode() {
	// Test unknown error code defaults to 500
	unknownCode := ErrorCode("UNKNOWN_ERROR")
	assert.Equal(suite.T(), 500, unknownCode.HTTPStatus())
}

// Test ServerError struct
func (suite *ServerErrorTestSuite) TestNewError() {
	// Arrange
	code := ErrorCodeBadRequest
	message := "Invalid input provided"

	// Act
	err := NewError(code, message)

	// Assert
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), code, err.Code)
	assert.Equal(suite.T(), message, err.Message)
}

func (suite *ServerErrorTestSuite) TestServerError_Error() {
	// Arrange
	code := ErrorCodeUnauthorized
	message := "Access denied"
	err := NewError(code, message)

	// Act
	errorMessage := err.Error()

	// Assert
	assert.Equal(suite.T(), message, errorMessage)
}

func (suite *ServerErrorTestSuite) TestServerError_ErrorInterface() {
	// Test that ServerError implements the error interface
	var err error = NewError(ErrorCodeInternalServerError, "Test error")
	
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), "Test error", err.Error())
}

// Test edge cases
func (suite *ServerErrorTestSuite) TestNewError_EmptyMessage() {
	// Arrange
	code := ErrorCodeNotFound
	message := ""

	// Act
	err := NewError(code, message)

	// Assert
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), code, err.Code)
	assert.Equal(suite.T(), message, err.Message)
	assert.Equal(suite.T(), "", err.Error())
}

func (suite *ServerErrorTestSuite) TestNewError_LongMessage() {
	// Arrange
	code := ErrorCodeInternalServerError
	message := "This is a very long error message that contains detailed information about what went wrong in the system and provides extensive context for debugging purposes"

	// Act
	err := NewError(code, message)

	// Assert
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), code, err.Code)
	assert.Equal(suite.T(), message, err.Message)
	assert.Equal(suite.T(), message, err.Error())
}

func (suite *ServerErrorTestSuite) TestNewError_WithSpecialCharacters() {
	// Arrange
	code := ErrorCodeBadRequest
	message := "Error with special characters: àáâãäå æç èéêë ìíîï ñ òóôõö ùúûü ý"

	// Act
	err := NewError(code, message)

	// Assert
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), code, err.Code)
	assert.Equal(suite.T(), message, err.Message)
	assert.Equal(suite.T(), message, err.Error())
}

// Test JSON serialization behavior (struct tags)
func (suite *ServerErrorTestSuite) TestServerError_StructFields() {
	// Test that the struct has the expected field names and JSON tags
	err := NewError(ErrorCodeForbidden, "Access forbidden")
	
	// Check field values are accessible
	assert.Equal(suite.T(), ErrorCodeForbidden, err.Code)
	assert.Equal(suite.T(), "Access forbidden", err.Message)
	
	// Note: JSON tag testing would require reflection or actual JSON marshaling
	// which is typically tested in integration tests
}

func TestServerErrorTestSuite(t *testing.T) {
	suite.Run(t, new(ServerErrorTestSuite))
}

// Additional unit tests without suite structure
func TestErrorCode_HTTPStatus_TableDriven(t *testing.T) {
	testCases := []struct {
		name           string
		code           ErrorCode
		expectedStatus int
	}{
		{"Internal Server Error", ErrorCodeInternalServerError, 500},
		{"Bad Request", ErrorCodeBadRequest, 400},
		{"Not Found", ErrorCodeNotFound, 404},
		{"Unauthorized", ErrorCodeUnauthorized, 401},
		{"Forbidden", ErrorCodeForbidden, 403},
		{"Conflict", ErrorCodeConflict, 409},
		{"Too Many Requests", ErrorCodeTooManyRequests, 429},
		{"Service Unavailable", ErrorCodeServiceUnavailable, 503},
		{"Unknown Code", ErrorCode("UNKNOWN"), 500},
		{"Empty Code", ErrorCode(""), 500},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			status := tc.code.HTTPStatus()
			assert.Equal(t, tc.expectedStatus, status)
		})
	}
}

func TestNewError_TableDriven(t *testing.T) {
	testCases := []struct {
		name    string
		code    ErrorCode
		message string
	}{
		{"Standard error", ErrorCodeBadRequest, "Standard error message"},
		{"Empty message", ErrorCodeUnauthorized, ""},
		{"Long message", ErrorCodeInternalServerError, "This is a very long error message with lots of details"},
		{"Special characters", ErrorCodeNotFound, "Error with émojis 🚨 and symbols @#$%^&*()"},
		{"Newlines in message", ErrorCodeForbidden, "Error\nwith\nnewlines"},
		{"Unicode message", ErrorCodeConflict, "错误信息 エラーメッセージ"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := NewError(tc.code, tc.message)
			
			assert.NotNil(t, err)
			assert.Equal(t, tc.code, err.Code)
			assert.Equal(t, tc.message, err.Message)
			assert.Equal(t, tc.message, err.Error())
		})
	}
}

// Benchmark tests
func BenchmarkNewError(b *testing.B) {
	code := ErrorCodeBadRequest
	message := "Test error message"

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = NewError(code, message)
	}
}

func BenchmarkErrorCode_HTTPStatus(b *testing.B) {
	code := ErrorCodeBadRequest

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = code.HTTPStatus()
	}
}

func BenchmarkServerError_Error(b *testing.B) {
	err := NewError(ErrorCodeInternalServerError, "Benchmark error message")

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = err.Error()
	}
}

// Test compatibility with standard error interface
func TestServerError_ErrorInterface_Compatibility(t *testing.T) {
	// Test that our error can be used where standard error is expected
	serverErr := NewError(ErrorCodeBadRequest, "Test error")
	
	// Function that accepts error interface
	processError := func(err error) string {
		if err != nil {
			return err.Error()
		}
		return "no error"
	}
	
	result := processError(serverErr)
	assert.Equal(t, "Test error", result)
}

// Test error code string representation
func TestErrorCode_StringRepresentation(t *testing.T) {
	codes := []ErrorCode{
		ErrorCodeInternalServerError,
		ErrorCodeBadRequest,
		ErrorCodeNotFound,
		ErrorCodeUnauthorized,
		ErrorCodeForbidden,
		ErrorCodeConflict,
		ErrorCodeTooManyRequests,
		ErrorCodeServiceUnavailable,
	}

	for _, code := range codes {
		// String method should return the same as string conversion
		assert.Equal(t, string(code), code.String())
		
		// Should not be empty
		assert.NotEmpty(t, code.String())
	}
}
