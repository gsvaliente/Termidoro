package config

import (
	"testing"
	"time"
)

func TestParseDuration(t *testing.T) {
	testCases := []struct {
		name        string
		durationStr string
		expected    time.Duration
		expectErr   bool
	}{
		{
			name:        "Empty string",
			durationStr: "",
			expected:    0,
			expectErr:   false,
		},
		{
			name:        "Minutes only",
			durationStr: "25m",
			expected:    25 * time.Minute,
			expectErr:   false,
		},
		{
			name:        "Hours and minutes",
			durationStr: "1h30m",
			expected:    90 * time.Minute,
			expectErr:   false,
		},
		{
			name:        "Seconds only",
			durationStr: "90s",
			expected:    90 * time.Second,
			expectErr:   false,
		},
		{
			name:        "Number without unit",
			durationStr: "45",
			expected:    45 * time.Minute,
			expectErr:   false,
		},
		{
			name:        "Invalid format",
			durationStr: "abc",
			expected:    0,
			expectErr:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			duration, err := parseDuration(tc.durationStr)
			if (err != nil) != tc.expectErr {
				t.Errorf("Expected error: %v, got: %v", tc.expectErr, err)
			}
			if duration != tc.expected {
				t.Errorf("Expected duration: %v, got: %v", tc.expected, duration)
			}
		})
	}
}
