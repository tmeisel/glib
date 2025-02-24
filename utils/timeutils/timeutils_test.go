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
