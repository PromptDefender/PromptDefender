package aiprompt

import (
	"context"
	"fmt"

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
	ctx := context.Background()
	model := v.Client.GenerativeModel(v.ModelName)

	// Match OpenAI implementation prompts if needed, but here we just send the prompt.
	// OpenAI impl had system prompt "You are a helpful assistant." + user prompt.
	// We can simulate that or just send the prompt.
	// Let's stick to the prompt.

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, err
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("no response from Vertex AI")
	}

	// Assuming text response
	part := resp.Candidates[0].Content.Parts[0]
	if txt, ok := part.(genai.Text); ok {
		s := string(txt)
		return &s, nil
	}

	return nil, fmt.Errorf("unexpected response format")
}

func (v *VertexAI) Close() {
	v.Client.Close()
}
