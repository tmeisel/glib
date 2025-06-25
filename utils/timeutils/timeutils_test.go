package timeutils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddDate(t *testing.T) {
	type testCase struct {
		timestamp   time.Time
		durString   string
		expected    time.Time
		expectError bool
	}

	reference, err := time.Parse("2006-01-02, 15:04", "1983-12-19, 10:00")
	if err != nil {
		t.Fatalf("failed to parse reference date: %v", err)
	}

	for name, tc := range map[string]testCase{
		"empty duration": {
			expectError: true,
		},
		"1 year": {
			timestamp: reference,
			durString: "1y",
			expected:  reference.AddDate(1, 0, 0),
		},
	} {
		t.Run(name, func(t *testing.T) {
			out, err := AddDate(tc.timestamp, tc.durString)

			if tc.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "invalid duration string")
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.expected, out)
		})
	}
}

func TestParseDate(t *testing.T) {
	type testCase struct {
		Input         string
		ExpectedYear  int
		ExpectedMonth time.Month
		ExpectedDay   int
		ExpectError   bool
	}

	for name, tc := range map[string]testCase{
		"19831219": {
			Input:         "19831219",
			ExpectedYear:  1983,
			ExpectedMonth: time.Month(12),
			ExpectedDay:   19,
			ExpectError:   false,
		},
		"too short": {
			Input:       "2012051",
			ExpectError: true,
		},
		"wrong month": {
			Input:       "20121301",
			ExpectError: true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			y, m, d, err := ParseDate(tc.Input)
			if tc.ExpectError {
				require.Error(t, err)

				assert.Equal(t, 1, y)
				assert.Equal(t, time.January, m)
				assert.Equal(t, 1, d)

				return
			}

			require.NoError(t, err)

			assert.Equal(t, tc.ExpectedYear, y)
			assert.Equal(t, tc.ExpectedMonth, m)
			assert.Equal(t, tc.ExpectedDay, d)
		})
	}
}
