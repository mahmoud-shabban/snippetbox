package main

import (
	"testing"
	"time"

	"github.com/mahmoud-shabban/snippetbox/internal/assert"
)

func TestHumanDate(t *testing.T) {

	cases := []struct {
		name     string
		input    time.Time
		expected string
	}{
		{
			name:     "zeroTime",
			input:    time.Time{},
			expected: "",
		},
		{
			name:     "localZone",
			input:    time.Date(2025, time.September, 1, 11, 0, 0, 0, time.FixedZone("EEST", 3*60*60)),
			expected: "01 Sep 2025 at 11:00",
		},
		{
			name:     "EETZoneConversion",
			input:    time.Date(2025, time.September, 1, 8, 0, 0, 0, time.UTC),
			expected: "01 Sep 2025 at 11:00",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			hd := humanDate(tt.input)

			assert.Equal(t, hd, tt.expected)
		})
	}
}
