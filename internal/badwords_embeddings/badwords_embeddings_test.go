//go:build ignore
// +build ignore

package badwords_embeddings

import (
	"testing"

	"github.com/safetorun/PromptDefender/embeddings"
	"github.com/stretchr/testify/mock"
)

type MockEmbeddings struct {
	mock.Mock
}

func (m *MockEmbeddings) CreateEmbeddings(prompt string) (*embeddings.EmbeddingValue, error) {
	args := m.Called(prompt)
	return args.Get(0).(*embeddings.EmbeddingValue), args.Error(1)
}

func (m *MockEmbeddings) RetrieveBadwordsEmbeddings() (*[]embeddings.EmbeddingValue, error) {
	args := m.Called()
	return args.Get(0).(*[]embeddings.EmbeddingValue), args.Error(1)
}

func TestGetClosestMatch(t *testing.T) {
	mockEmbeddings := new(MockEmbeddings)
	samplePromptEmbeddings := &embeddings.EmbeddingValue{EmbeddingValue: []float64{0.1, 0.2, 0.3}}
	sampleBadwordEmbeddings := &[]embeddings.EmbeddingValue{
		{EmbeddingValue: []float64{0.4, 0.5, 0.6}},
		{EmbeddingValue: []float64{0.7, 0.8, 0.9}},
	}
	mockEmbeddings.On("CreateEmbeddings", "test prompt").Return(samplePromptEmbeddings, nil)
	mockEmbeddings.On("RetrieveBadwordsEmbeddings").Return(sampleBadwordEmbeddings, nil)
	bw := New(mockEmbeddings)
	result, err := bw.GetClosestMatch("test prompt")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result.Score >= 1.0 {
		t.Errorf("Score is too high: %v", result.Score)
	}
}
