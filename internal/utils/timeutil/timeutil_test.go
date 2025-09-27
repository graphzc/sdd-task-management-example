package timeutil

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TimeUtilTestSuite struct {
	suite.Suite
}

func (suite *TimeUtilTestSuite) TestGetBangkokLocation_Success() {
	// Act
	location, err := GetBangkokLocation()

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), location)
	assert.Equal(suite.T(), "Asia/Bangkok", location.String())
}

func (suite *TimeUtilTestSuite) TestGetBangkokLocation_Consistency() {
	// Act - Call multiple times
	location1, err1 := GetBangkokLocation()
	location2, err2 := GetBangkokLocation()

	// Assert - Should return consistent results
	assert.NoError(suite.T(), err1)
	assert.NoError(suite.T(), err2)
	assert.NotNil(suite.T(), location1)
	assert.NotNil(suite.T(), location2)
	assert.Equal(suite.T(), location1.String(), location2.String())
}

func (suite *TimeUtilTestSuite) TestBangkokNow_ReturnsValidTime() {
	// Act
	bangkokTime := BangkokNow()

	// Assert
	assert.NotNil(suite.T(), bangkokTime)
	assert.False(suite.T(), bangkokTime.IsZero())

	// Check timezone
	zone, _ := bangkokTime.Zone()
	assert.Contains(suite.T(), zone, "+07") // Bangkok is UTC+7
}

func (suite *TimeUtilTestSuite) TestBangkokNow_IsInBangkokTimezone() {
	// Act
	bangkokTime := BangkokNow()

	// Assert
	expectedLocation, _ := time.LoadLocation("Asia/Bangkok")
	assert.Equal(suite.T(), expectedLocation.String(), bangkokTime.Location().String())
}

func (suite *TimeUtilTestSuite) TestBangkokNow_IsRecentTime() {
	// Arrange
	before := time.Now().UTC()

	// Act
	bangkokTime := BangkokNow()

	// Arrange
	after := time.Now().UTC()

	// Assert - Bangkok time should be between before and after when converted to UTC
	bangkokUTC := bangkokTime.UTC()
	assert.True(suite.T(), bangkokUTC.After(before.Add(-time.Second)) || bangkokUTC.Equal(before.Add(-time.Second)))
	assert.True(suite.T(), bangkokUTC.Before(after.Add(time.Second)) || bangkokUTC.Equal(after.Add(time.Second)))
}

func (suite *TimeUtilTestSuite) TestBangkokNow_DifferentFromUTC() {
	// Act
	bangkokTime := BangkokNow()
	utcTime := time.Now().UTC()

	// Convert Bangkok time to UTC for comparison
	bangkokInUTC := bangkokTime.UTC()

	// Assert - The times should be very close (within a few seconds)
	// but the locations should be different
	timeDiff := bangkokInUTC.Sub(utcTime)
	assert.True(suite.T(), timeDiff < time.Second*5 && timeDiff > -time.Second*5)

	// Location should be different
	assert.NotEqual(suite.T(), utcTime.Location().String(), bangkokTime.Location().String())
}

func (suite *TimeUtilTestSuite) TestBangkokNow_ConsistentOffset() {
	// Act
	bangkokTime := BangkokNow()

	// Assert - Bangkok should be UTC+7
	_, offset := bangkokTime.Zone()
	expectedOffset := 7 * 60 * 60 // 7 hours in seconds
	assert.Equal(suite.T(), expectedOffset, offset)
}

// Test multiple calls in sequence
func (suite *TimeUtilTestSuite) TestBangkokNow_Sequential() {
	// Act
	time1 := BangkokNow()
	time.Sleep(10 * time.Millisecond) // Small delay
	time2 := BangkokNow()

	// Assert
	assert.True(suite.T(), time2.After(time1) || time2.Equal(time1))
	assert.Equal(suite.T(), time1.Location().String(), time2.Location().String())
}

func TestTimeUtilTestSuite(t *testing.T) {
	suite.Run(t, new(TimeUtilTestSuite))
}

// Additional unit tests
func TestGetBangkokLocation_ErrorHandling(t *testing.T) {
	// This test verifies the function handles the timezone correctly
	// In most systems, Asia/Bangkok should be available
	location, err := GetBangkokLocation()

	// Asia/Bangkok should be available in standard timezone databases
	assert.NoError(t, err)
	assert.NotNil(t, location)

	if location != nil {
		assert.Equal(t, "Asia/Bangkok", location.String())
	}
}

func TestBangkokNow_Precision(t *testing.T) {
	// Test that we get nanosecond precision
	time1 := BangkokNow()
	time2 := BangkokNow()

	// Even if called immediately after each other, they should be different
	// (or at least one should not be significantly less precise than the other)
	assert.True(t, time2.UnixNano() >= time1.UnixNano())
}

// Benchmark tests
func BenchmarkGetBangkokLocation(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := GetBangkokLocation()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBangkokNow(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = BangkokNow()
	}
}

// Table-driven tests for edge cases
func TestTimeUtil_EdgeCases(t *testing.T) {
	testCases := []struct {
		name     string
		testFunc func(t *testing.T)
	}{
		{
			name: "Bangkok location is valid timezone",
			testFunc: func(t *testing.T) {
				location, err := GetBangkokLocation()
				assert.NoError(t, err)

				// Should be able to create a time in this location
				testTime := time.Date(2023, 1, 1, 12, 0, 0, 0, location)
				assert.Equal(t, location, testTime.Location())
			},
		},
		{
			name: "Bangkok time has correct offset",
			testFunc: func(t *testing.T) {
				bangkokTime := BangkokNow()

				// Get the offset
				_, offset := bangkokTime.Zone()

				// Bangkok is UTC+7, so offset should be 7*3600 = 25200 seconds
				assert.Equal(t, 25200, offset)
			},
		},
		{
			name: "Bangkok time zone name",
			testFunc: func(t *testing.T) {
				bangkokTime := BangkokNow()

				// Get zone name - it might be "+07" or "ICT" depending on system
				zoneName, _ := bangkokTime.Zone()

				// Should contain either +07 or ICT (Indochina Time)
				assert.True(t, zoneName == "+07" || zoneName == "ICT" || zoneName == "Asia/Bangkok")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, tc.testFunc)
	}
}

// Test daylight saving time considerations
func TestBangkokNow_NoDaylightSaving(t *testing.T) {
	// Bangkok doesn't observe daylight saving time
	// So the offset should be consistent throughout the year

	location, err := GetBangkokLocation()
	assert.NoError(t, err)

	// Test different times of year
	dates := []time.Time{
		time.Date(2023, 1, 15, 12, 0, 0, 0, location),  // January
		time.Date(2023, 6, 15, 12, 0, 0, 0, location),  // June
		time.Date(2023, 12, 15, 12, 0, 0, 0, location), // December
	}

	var offsets []int
	for _, date := range dates {
		_, offset := date.Zone()
		offsets = append(offsets, offset)
	}

	// All offsets should be the same (no DST)
	for i := 1; i < len(offsets); i++ {
		assert.Equal(t, offsets[0], offsets[i], "Bangkok should not observe daylight saving time")
	}
}

// Test concurrent access
func TestBangkokNow_Concurrent(t *testing.T) {
	const numGoroutines = 100
	results := make(chan time.Time, numGoroutines)

	// Launch concurrent calls
	for i := 0; i < numGoroutines; i++ {
		go func() {
			results <- BangkokNow()
		}()
	}

	// Collect results
	times := make([]time.Time, 0, numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		times = append(times, <-results)
	}

	// All should have Bangkok timezone
	for _, timeVal := range times {
		assert.Equal(t, "Asia/Bangkok", timeVal.Location().String())
	}

	// Times should be reasonable (within a few seconds of each other)
	minTime := times[0]
	maxTime := times[0]

	for _, timeVal := range times[1:] {
		if timeVal.Before(minTime) {
			minTime = timeVal
		}
		if timeVal.After(maxTime) {
			maxTime = timeVal
		}
	}

	// All times should be within a few seconds
	assert.True(t, maxTime.Sub(minTime) < 5*time.Second)
}
