package wall

import (
	"errors"
	"github.com/safetorun/PromptDefender/utils"
	"testing"
	"time"
)

func TestRetry(t *testing.T) {
	var attemptsMade int
	mockFunc := func() (*float64, error) {
		attemptsMade++
		if attemptsMade < 3 {
			return nil, errors.New("mock error")
		}
		val := 1.0
		return &val, nil
	}

	val, err := utils.Retry(5, 1*time.Millisecond, mockFunc)
	if err != nil {
		t.Errorf("Retry function failed: %v", err)
	}
	if *val != 1.0 {
		t.Errorf("Retry function returned incorrect value: got %v want %v", *val, 1.0)
	}
	if attemptsMade != 3 {
		t.Errorf("Retry function did not make the correct number of attempts: got %v want %v", attemptsMade, 3)
	}
}
