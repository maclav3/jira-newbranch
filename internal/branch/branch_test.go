package branch

import (
	"testing"
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
			got := FormatBranchName(tt.key, tt.summary)
			if got != tt.expected {
				t.Errorf("FormatBranchName() = %v, want %v", got, tt.expected)
			}
		})
	}
}
