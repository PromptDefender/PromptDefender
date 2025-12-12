package aiprompt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const OpenAIAPIURL = "https://api.openai.com/v1/chat/completions"

type OpenAI struct {
	ApiKey string
}

func NewOpenAI(apiKey string) RemoteAIChecker {
	return OpenAI{ApiKey: apiKey}
}

// Message represents a message for the OpenAI API.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Payload represents the payload for the OpenAI API request.
type Payload struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type OpenAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func (o OpenAI) CheckAI(prompt string) (*string, error) {
	payloadData := Payload{
		Model: "gpt-4",
		Messages: []Message{
			{
				Role:    "system",
				Content: "You are a helpful assistant.",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	payload, err := json.Marshal(payloadData)
	if err != nil {
		return nil, err
	}

	resp, err := post("application/json", o.ApiKey, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var response OpenAIResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	if len(response.Choices) > 0 {
		return &response.Choices[0].Message.Content, nil
	} else {
		return nil, fmt.Errorf("no response from OpenAI")
	}
}

func post(contentType string, apiKey string, body io.Reader) (*http.Response, error) {
	client := &http.Client{}

	req, err := http.NewRequest("POST", OpenAIAPIURL, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", "Bearer "+apiKey)

	return client.Do(req)
}

//func post(contentType string, apiKey string, body io.Reader) (*http.Response, error) {
//	client := spinhttp.NewClient()
//
//	parsedUrl, err := url.Parse(OpenAIAPIURL)
//
//	if err != nil {
//		return nil, err
//	}
//
//	return client.Do(
//		&http.Request{
//			Method: "POST",
//			URL:    parsedUrl,
//			Header: http.Header{
//				"Content-Type":  []string{contentType},
//				"Authorization": []string{"Bearer " + apiKey},
//			},
//			Body: ioutil.NopCloser(body),
//		})
//
//}
