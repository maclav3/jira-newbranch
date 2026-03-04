package jira

import (
	"testing"
	"time"
)

func TestParseJiraTime(t *testing.T) {
	tests := []struct {
		input    string
		expected string // We'll compare formatted as RFC3339
	}{
		{
			input:    "2026-02-16T13:51:20.182+0000",
			expected: "2026-02-16T13:51:20Z",
		},
		{
			input:    "2026-02-16T13:51:20Z",
			expected: "2026-02-16T13:51:20Z",
		},
		{
			input:    "2026-02-16T13:51:20.123+0100",
			expected: "2026-02-16T13:51:20+01:00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ParseJiraTime(tt.input)
			if err != nil {
				t.Fatalf("ParseJiraTime(%q) error: %v", tt.input, err)
			}
			gotStr := got.Format(time.RFC3339)
			if gotStr != tt.expected {
				t.Errorf("ParseJiraTime(%q) = %v, want %v", tt.input, gotStr, tt.expected)
			}
		})
	}
}
