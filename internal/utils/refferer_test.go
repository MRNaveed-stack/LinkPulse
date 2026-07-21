package utils_test

import (
	"testing"

	"github.com/MRNaveed-stack/LinkPulse/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeReferrer(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Empty Referrer",
			input:    "",
			expected: "Direct",
		},
		{
			name:     "Invalid URL format",
			input:    "://invalid-url",
			expected: "Unknown",
		},
		{
			name:     "Instagram Referrer",
			input:    "https://www.instagram.com/p/abc",
			expected: "Instagram",
		},
		{
			name:     "Facebook Referrer",
			input:    "https://facebook.com/profile.php",
			expected: "Facebook",
		},
		{
			name:     "Twitter Referrer",
			input:    "http://twitter.com/status/123",
			expected: "Twitter",
		},
		{
			name:     "LinkedIn Referrer",
			input:    "https://www.linkedin.com/feed/",
			expected: "LinkedIn",
		},
		{
			name:     "TikTok Referrer",
			input:    "https://tiktok.com/@user",
			expected: "TikTok",
		},
		{
			name:     "Reddit Referrer",
			input:    "https://www.reddit.com/r/golang",
			expected: "Reddit",
		},
		{
			name:     "GitHub Referrer",
			input:    "https://github.com/MRNaveed-stack",
			expected: "GitHub",
		},
		{
			name:     "Other Domain",
			input:    "https://example.com/some/path",
			expected: "Other",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := utils.NormalizeReferrer(tt.input)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
