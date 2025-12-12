package aiprompt

import (
	"strings"
	"testing"
)

func TestNormalizeString(t *testing.T) {
	testCases := []struct {
		input string
		want  string
	}{
		{" Test String!  ", "test string"},
		{"ANOTHER TEST_ STRING!!", "another test string"},
		{" ", ""},
		{"", ""},
	}

	for _, tc := range testCases {
		got := normalizeString(tc.input)
		if got != tc.want {
			t.Errorf("normalizeString(%q) = %q; want %q", tc.input, got, tc.want)
		}
	}
}

func TestRenderPromptForPiDetection(t *testing.T) {
	userInput := "test"
	got := RenderPromptForPiDetection(userInput)
	if !strings.Contains(got, userInput) {
		t.Errorf("Expected renderPromptForPiDetection to contain %q, but it did not.", userInput)
	}
}
