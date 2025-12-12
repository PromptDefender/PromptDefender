package wall

type BadWordsResult struct {
	ContainsBadWords bool
}

var badWords = []string{
	"Ignore previous",
	"Do anything now",
}

// CheckForBadWords checks for bad words in the input string.
func CheckForBadWords(input string) BadWordsResult {
	// Check if input contains any of the bad words.
	// If it does, return a BadWordsResult with Detected set to true.
	// Otherwise, return a BadWordsResult with Detected set to false.

	return BadWordsResult{
		ContainsBadWords: false,
	}
}
