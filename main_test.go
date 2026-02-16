package main

import (
	"testing"
	"time"
)

func TestFormatBranchName(t *testing.T) {
	tests := []struct {
		key      string
		summary  string
		expected string
	}{
		{
			key:      "PRJ-2134",
			summary:  "Adjust foo and bar to be ready for prod",
			expected: "PRJ-2134-adjust-foo-bar-ready",
		},
		{
			key:      "PROJ-1",
			summary:  "The quick brown fox jumps over the lazy dog",
			expected: "PROJ-1-quick-brown-fox-jumps",
		},
		{
			key:      "TASK-123",
			summary:  "fix: use-correct-token or something",
			expected: "TASK-123-fix-use-correct-token",
		},
		{
			key:      "A-1",
			summary:  "to and or be for with",
			expected: "A-1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			got := formatBranchName(tt.key, tt.summary)
			if got != tt.expected {
				t.Errorf("formatBranchName() = %v, want %v", got, tt.expected)
			}
		})
	}
}

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
			got, err := parseJiraTime(tt.input)
			if err != nil {
				t.Fatalf("parseJiraTime(%q) error: %v", tt.input, err)
			}
			gotStr := got.Format(time.RFC3339)
			if gotStr != tt.expected {
				t.Errorf("parseJiraTime(%q) = %v, want %v", tt.input, gotStr, tt.expected)
			}
		})
	}
}
