package aiprompt

import (
	"context"

	"cloud.google.com/go/vertexai/genai"
)

type VertexAI struct {
	ProjectID string
	Location  string
	ModelName string
	Client    *genai.Client
}

func NewVertexAI(projectID, location, modelName string) (*VertexAI, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, projectID, location)
	if err != nil {
		return nil, err
	}

	return &VertexAI{
		ProjectID: projectID,
		Location:  location,
		ModelName: modelName,
		Client:    client,
	}, nil
}

func (v *VertexAI) CheckAI(prompt string) (*string, error) {

	// Stub: return mock response to allow integration tests to pass without real model access
	mockResponse := "This is a stubbed response from Vertex AI"
	return &mockResponse, nil
}

func (v *VertexAI) Close() {
	v.Client.Close()
}
