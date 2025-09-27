package validator

import (
	"testing"

	"github.com/graphzc/sdd-task-management-example/internal/utils/servererr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ValidatorTestSuite struct {
	suite.Suite
	validator *Validator
}

func (suite *ValidatorTestSuite) SetupTest() {
	suite.validator = NewValidator()
}

// Test structs for validation
type TestValidStruct struct {
	Name  string `validate:"required,min=2,max=50"`
	Email string `validate:"required,email"`
	Age   int    `validate:"min=0,max=150"`
}

type TestInvalidStruct struct {
	Name  string `validate:"required,min=2,max=50"`
	Email string `validate:"required,email"`
	Age   int    `validate:"min=0,max=150"`
}

type TestNestedStruct struct {
	User    TestValidStruct `validate:"required"`
	Address string          `validate:"required,min=5"`
}

type TestComplexStruct struct {
	Username string   `validate:"required,alphanum,min=3,max=20"`
	Password string   `validate:"required,min=8"`
	Tags     []string `validate:"required,dive,required"`
	IsActive bool     // No validation
}

// Test NewValidator
func (suite *ValidatorTestSuite) TestNewValidator() {
	validator := NewValidator()
	assert.NotNil(suite.T(), validator)
	assert.NotNil(suite.T(), validator.validate)
}

// Test ValidateStruct with valid data
func (suite *ValidatorTestSuite) TestValidateStruct_ValidData() {
	// Arrange
	validStruct := TestValidStruct{
		Name:  "John Doe",
		Email: "john@example.com",
		Age:   25,
	}

	// Act
	err := suite.validator.ValidateStruct(validStruct)

	// Assert
	assert.NoError(suite.T(), err)
}

// Test ValidateStruct with missing required fields
func (suite *ValidatorTestSuite) TestValidateStruct_MissingRequiredFields() {
	// Arrange
	invalidStruct := TestValidStruct{
		Name:  "", // Empty required field
		Email: "john@example.com",
		Age:   25,
	}

	// Act
	err := suite.validator.ValidateStruct(invalidStruct)

	// Assert
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "name is required")
}

// Test ValidateStruct with invalid email
func (suite *ValidatorTestSuite) TestValidateStruct_InvalidEmail() {
	// Arrange
	invalidStruct := TestValidStruct{
		Name:  "John Doe",
		Email: "invalid-email", // Invalid email format
		Age:   25,
	}

	// Act
	err := suite.validator.ValidateStruct(invalidStruct)

	// Assert
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "email is email")
}

// Test ValidateStruct with multiple validation errors
func (suite *ValidatorTestSuite) TestValidateStruct_MultipleErrors() {
	// Arrange
	invalidStruct := TestValidStruct{
		Name:  "A",             // Too short
		Email: "invalid-email", // Invalid format
		Age:   200,             // Too high
	}

	// Act
	err := suite.validator.ValidateStruct(invalidStruct)

	// Assert
	assert.Error(suite.T(), err)
	errorMsg := err.Error()
	assert.Contains(suite.T(), errorMsg, "name is min")
	assert.Contains(suite.T(), errorMsg, "email is email")
	assert.Contains(suite.T(), errorMsg, "age is max")
}

// Test ValidateStruct with nested structures
func (suite *ValidatorTestSuite) TestValidateStruct_NestedStruct_Valid() {
	// Arrange
	validNested := TestNestedStruct{
		User: TestValidStruct{
			Name:  "John Doe",
			Email: "john@example.com",
			Age:   25,
		},
		Address: "123 Main Street",
	}

	// Act
	err := suite.validator.ValidateStruct(validNested)

	// Assert
	assert.NoError(suite.T(), err)
}

// Test ValidateStruct with invalid nested structure
func (suite *ValidatorTestSuite) TestValidateStruct_NestedStruct_Invalid() {
	// Arrange
	invalidNested := TestNestedStruct{
		User: TestValidStruct{
			Name:  "", // Invalid nested field
			Email: "john@example.com",
			Age:   25,
		},
		Address: "123 Main Street",
	}

	// Act
	err := suite.validator.ValidateStruct(invalidNested)

	// Assert
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "name is required")
}

// Test ValidateStruct with complex validation rules
func (suite *ValidatorTestSuite) TestValidateStruct_ComplexRules() {
	// Arrange
	validComplex := TestComplexStruct{
		Username: "johndoe123",
		Password: "strongpassword123",
		Tags:     []string{"tag1", "tag2", "tag3"},
		IsActive: true,
	}

	// Act
	err := suite.validator.ValidateStruct(validComplex)

	// Assert
	assert.NoError(suite.T(), err)
}

// Test ValidateStruct with invalid complex rules
func (suite *ValidatorTestSuite) TestValidateStruct_ComplexRules_Invalid() {
	// Arrange
	invalidComplex := TestComplexStruct{
		Username: "jo",                         // Too short
		Password: "weak",                       // Too short
		Tags:     []string{"tag1", "", "tag3"}, // Empty tag in slice
		IsActive: false,
	}

	// Act
	err := suite.validator.ValidateStruct(invalidComplex)

	// Assert
	assert.Error(suite.T(), err)
	errorMsg := err.Error()
	assert.Contains(suite.T(), errorMsg, "username is min")
	assert.Contains(suite.T(), errorMsg, "password is min")
}

// Test ValidateStruct with nil input
func (suite *ValidatorTestSuite) TestValidateStruct_NilInput() {
	// Act
	err := suite.validator.ValidateStruct(nil)

	// Assert
	assert.Error(suite.T(), err)
}

// Test ValidateStruct with non-struct input
func (suite *ValidatorTestSuite) TestValidateStruct_NonStructInput() {
	// Arrange
	nonStruct := "this is a string"

	// Act
	err := suite.validator.ValidateStruct(nonStruct)

	// Assert
	assert.Error(suite.T(), err)
}

// Test echo validator wrapper
func (suite *ValidatorTestSuite) TestValidator_EchoValidate_ValidData() {
	// Arrange
	validStruct := TestValidStruct{
		Name:  "John Doe",
		Email: "john@example.com",
		Age:   25,
	}

	// Act
	err := suite.validator.Validate(validStruct)

	// Assert
	assert.NoError(suite.T(), err)
}

// Test echo validator wrapper with invalid data
func (suite *ValidatorTestSuite) TestValidator_EchoValidate_InvalidData() {
	// Arrange
	invalidStruct := TestValidStruct{
		Name:  "", // Required field is empty
		Email: "john@example.com",
		Age:   25,
	}

	// Act
	err := suite.validator.Validate(invalidStruct)

	// Assert
	assert.Error(suite.T(), err)

	// Check if it's a ServerError
	serverErr, ok := err.(*servererr.ServerError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), servererr.ErrorCodeBadRequest, serverErr.Code)
	assert.Contains(suite.T(), serverErr.Message, "name is required")
}

func TestValidatorTestSuite(t *testing.T) {
	suite.Run(t, new(ValidatorTestSuite))
}

// Additional table-driven tests
func TestValidateStruct_TableDriven(t *testing.T) {
	validator := NewValidator()

	testCases := []struct {
		name      string
		input     TestValidStruct
		expectErr bool
		errMsg    string
	}{
		{
			name: "Valid struct",
			input: TestValidStruct{
				Name:  "John Doe",
				Email: "john@example.com",
				Age:   25,
			},
			expectErr: false,
		},
		{
			name: "Empty name",
			input: TestValidStruct{
				Name:  "",
				Email: "john@example.com",
				Age:   25,
			},
			expectErr: true,
			errMsg:    "name is required",
		},
		{
			name: "Name too short",
			input: TestValidStruct{
				Name:  "J",
				Email: "john@example.com",
				Age:   25,
			},
			expectErr: true,
			errMsg:    "name is min",
		},
		{
			name: "Name too long",
			input: TestValidStruct{
				Name:  "This is a very long name that exceeds the maximum allowed length for the name field",
				Email: "john@example.com",
				Age:   25,
			},
			expectErr: true,
			errMsg:    "name is max",
		},
		{
			name: "Invalid email",
			input: TestValidStruct{
				Name:  "John Doe",
				Email: "not-an-email",
				Age:   25,
			},
			expectErr: true,
			errMsg:    "email is email",
		},
		{
			name: "Empty email",
			input: TestValidStruct{
				Name:  "John Doe",
				Email: "",
				Age:   25,
			},
			expectErr: true,
			errMsg:    "email is required",
		},
		{
			name: "Age too low",
			input: TestValidStruct{
				Name:  "John Doe",
				Email: "john@example.com",
				Age:   -1,
			},
			expectErr: true,
			errMsg:    "age is min",
		},
		{
			name: "Age too high",
			input: TestValidStruct{
				Name:  "John Doe",
				Email: "john@example.com",
				Age:   151,
			},
			expectErr: true,
			errMsg:    "age is max",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validator.ValidateStruct(tc.input)

			if tc.expectErr {
				assert.Error(t, err)
				if tc.errMsg != "" {
					assert.Contains(t, err.Error(), tc.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Benchmark tests
func BenchmarkValidator_ValidateStruct_Valid(b *testing.B) {
	validator := NewValidator()
	validStruct := TestValidStruct{
		Name:  "John Doe",
		Email: "john@example.com",
		Age:   25,
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		err := validator.ValidateStruct(validStruct)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkValidator_ValidateStruct_Invalid(b *testing.B) {
	validator := NewValidator()
	invalidStruct := TestValidStruct{
		Name:  "", // This will cause validation error
		Email: "john@example.com",
		Age:   25,
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = validator.ValidateStruct(invalidStruct)
	}
}

func BenchmarkValidator_EchoValidate(b *testing.B) {
	validator := NewValidator()
	validStruct := TestValidStruct{
		Name:  "John Doe",
		Email: "john@example.com",
		Age:   25,
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		err := validator.Validate(validStruct)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Test edge cases
func TestValidator_EdgeCases(t *testing.T) {
	validator := NewValidator()

	t.Run("Empty struct", func(t *testing.T) {
		type EmptyStruct struct{}
		err := validator.ValidateStruct(EmptyStruct{})
		assert.NoError(t, err)
	})

	t.Run("Struct with no validation tags", func(t *testing.T) {
		type NoValidationStruct struct {
			Name  string
			Email string
			Age   int
		}

		err := validator.ValidateStruct(NoValidationStruct{
			Name:  "",
			Email: "invalid",
			Age:   -1,
		})
		assert.NoError(t, err) // No validation rules, so no errors
	})

	t.Run("Pointer to struct", func(t *testing.T) {
		validStruct := &TestValidStruct{
			Name:  "John Doe",
			Email: "john@example.com",
			Age:   25,
		}
		err := validator.ValidateStruct(validStruct)
		assert.NoError(t, err)
	})

	t.Run("Nil pointer", func(t *testing.T) {
		var nilStruct *TestValidStruct
		err := validator.ValidateStruct(nilStruct)
		assert.Error(t, err)
	})
}

// Test concurrent usage
func TestValidator_Concurrent(t *testing.T) {
	validator := NewValidator()
	const numGoroutines = 100

	// Channel to collect results
	results := make(chan error, numGoroutines)

	// Launch concurrent validations
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			testStruct := TestValidStruct{
				Name:  "John Doe",
				Email: "john@example.com",
				Age:   25,
			}

			err := validator.ValidateStruct(testStruct)
			results <- err
		}(i)
	}

	// Collect results
	for i := 0; i < numGoroutines; i++ {
		err := <-results
		assert.NoError(t, err)
	}
}
