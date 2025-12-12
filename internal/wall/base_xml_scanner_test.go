package wall

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBasicXmlScanner_Scan(t *testing.T) {
	scanner := NewBasicXmlEscapingScaner()

	tests := []struct {
		name           string
		textToScan     string
		tagToScanFor   string
		expectedResult bool
	}{
		{"Positive Test", "Sample text with <tag>", "tag", true},
		{"Negative Test", "Sample text without tag", "tag", false},
		{"Empty Text Test", "", "<tag>", false},
		{"Empty Tag Test", "Sample text with <tag>", "", false},
		{"End tag test", "Sample text with end </tag>", "tag", true},
		{"Both Empty Test", "", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := scanner.Scan(tt.textToScan, tt.tagToScanFor)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedResult, result.ContainsXmlEscaping)
		})
	}
}
