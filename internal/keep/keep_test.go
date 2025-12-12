package keep

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

type MockAIChecker struct{}

func (_ *MockAIChecker) CheckAI(_ string) (*string, error) {
	response := "response from AI"
	return &response, nil
}

type MockAICheckerError struct{}

func (_ *MockAICheckerError) CheckAI(_ string) (*string, error) {
	return nil, errors.New("AI error")
}

func TestBuildKeep(t *testing.T) {
	mockAIChecker := &MockAIChecker{}
	mockAICheckerError := &MockAICheckerError{}

	t.Run("BuildKeep_EmptyPrompt", func(t *testing.T) {
		k := New(mockAIChecker)
		startingPrompt := StartingPrompt{Prompt: ""}

		_, err := k.BuildKeep(startingPrompt)
		assert.Error(t, err)
		assert.Equal(t, "prompt cannot be empty", err.Error())
		assert.True(t, IsPromptRequiredError(err))
	})

	t.Run("BuildKeep_Success", func(t *testing.T) {
		k := New(mockAIChecker)
		startingPrompt := StartingPrompt{Prompt: "test prompt"}

		newPrompt, err := k.BuildKeep(startingPrompt)
		assert.NoError(t, err)
		assert.Equal(t, "response from AI", newPrompt.NewPrompt)
	})

	t.Run("BuildKeep_Error", func(t *testing.T) {
		k := New(mockAICheckerError)
		startingPrompt := StartingPrompt{Prompt: "test prompt"}

		_, err := k.BuildKeep(startingPrompt)
		assert.Error(t, err)
	})
}
