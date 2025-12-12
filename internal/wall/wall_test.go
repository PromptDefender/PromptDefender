package wall

import (
	"errors"
	"github.com/safetorun/PromptDefender/tracer"
	"testing"

	"github.com/safetorun/PromptDefender/badwords"
	"github.com/safetorun/PromptDefender/pii"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPiiScanner struct {
	mock.Mock
}

func (m *MockPiiScanner) Scan(text string) (*pii.ScanResult, error) {
	args := m.Called(text)
	if args.Get(0) != nil {
		return args.Get(0).(*pii.ScanResult), args.Error(1)
	}

	return nil, args.Error(1)
}

type MockBadWords struct {
	mock.Mock
}

var emptyTracer = &tracer.EmptyTracer{}

func (m *MockBadWords) CheckPromptContainsBadWords(text string) (*bool, error) {
	args := m.Called(text)

	if args.Get(0) != nil {
		return args.Get(0).(*bool), args.Error(1)
	}

	return args.Get(0).(*bool), args.Error(1)
}

type MockClosestMatcher struct {
	mock.Mock
}

func (m *MockClosestMatcher) GetClosestMatch(text string) (*badwords.ClosestMatchScore, error) {
	args := m.Called(text)
	if args.Get(0) != nil {
		return args.Get(0).(*badwords.ClosestMatchScore), args.Error(1)
	}

	return nil, args.Error(1)
}

func newWall(scanner pii.Scanner, badwordsCheck *badwords.BadWords) WallOpts {
	addAllConfigurations := func(c *Wall) error {
		c.PiiScanner = scanner
		c.BadWordsCheck = badwordsCheck

		return nil
	}

	return addAllConfigurations
}
func TestCheckWall(t *testing.T) {

	t.Run("Check with bad words", func(t *testing.T) {

		scanPii := false
		piiScanner := new(MockPiiScanner)
		closestMatcher := new(MockClosestMatcher)
		badWordsCheck := badwords.New(closestMatcher)
		mt := newWall(piiScanner, badWordsCheck)
		m, _ := New(mt)

		closestMatcher.On("GetClosestMatch", "some text").Return(&badwords.ClosestMatchScore{Score: badwords.Medium}, nil)
		check := PromptToCheck{Prompt: "some text", ScanPii: &scanPii}
		result, err := m.CheckWall(check, emptyTracer)
		assert.Nil(t, err)
		assert.NotNil(t, result.ContainsBadWords)
	})

	t.Run("Check with PII", func(t *testing.T) {

		scanPii := true
		piiScanner := new(MockPiiScanner)
		closestMatcher := new(MockClosestMatcher)
		badWordsCheck := badwords.New(closestMatcher)
		mt := newWall(piiScanner, badWordsCheck)
		m, _ := New(mt)

		piiScanner.On("Scan", "some text").Return(&pii.ScanResult{ContainingPii: true}, nil)
		closestMatcher.On("GetClosestMatch", "some text").Return(&badwords.ClosestMatchScore{Score: badwords.Medium}, nil)
		check := PromptToCheck{Prompt: "some text", ScanPii: &scanPii}
		result, err := m.CheckWall(check, emptyTracer)
		assert.Nil(t, err)
		assert.True(t, result.PiiResult.ContainsPii)
	})

	t.Run("Error while checking bad words", func(t *testing.T) {
		scanPii := false
		piiScanner := new(MockPiiScanner)
		closestMatcher := new(MockClosestMatcher)
		badWordsCheck := badwords.New(closestMatcher)

		mt := newWall(piiScanner, badWordsCheck)
		m, _ := New(mt)

		closestMatcher.On("GetClosestMatch", "some text").Return(nil, errors.New("an error"))
		check := PromptToCheck{Prompt: "some text", ScanPii: &scanPii}
		result, err := m.CheckWall(check, emptyTracer)
		assert.Nil(t, result)
		assert.Error(t, err)
	})

	t.Run("Error while scanning PII", func(t *testing.T) {
		scanPii := true
		piiScanner := new(MockPiiScanner)
		closestMatcher := new(MockClosestMatcher)
		badWordsCheck := badwords.New(closestMatcher)

		mt := newWall(piiScanner, badWordsCheck)
		m, _ := New(mt)

		piiScanner.On("Scan", "some text").Return(nil, errors.New("an error"))
		closestMatcher.On("GetClosestMatch", "some text").Return(&badwords.ClosestMatchScore{Score: badwords.Medium}, nil)
		check := PromptToCheck{Prompt: "some text", ScanPii: &scanPii}
		result, err := m.CheckWall(check, emptyTracer)
		assert.Nil(t, result)
		assert.Error(t, err)
	})
}
