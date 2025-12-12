package main

import (
	"context"
	"log"
	"net/http"
)

type SuspiciousUser interface {
	CheckSuspiciousUser(ctx context.Context, userId string) (*bool, error)
}

type RemoteSuspiciousUser struct {
	url    string
	apiKey string
	logger *log.Logger
}

func NewRemote(url string, apiKey string) RemoteSuspiciousUser {

	return RemoteSuspiciousUser{
		url:    url,
		apiKey: apiKey,
		logger: log.Default(),
	}
}

func (r *RemoteSuspiciousUser) CreateClient() (*ClientWithResponses, error) {

	addApiKey := func(c *Client) error {
		c.RequestEditors = append(c.RequestEditors, func(ctx context.Context, req *http.Request) error {
			req.Header.Add("x-api-key", r.apiKey)
			return nil
		})

		return nil
	}

	return NewClientWithResponses(r.url, addApiKey)
}

// CheckSuspiciousUser checks if a user is suspicious
func (r *RemoteSuspiciousUser) CheckSuspiciousUser(ctx context.Context, userId string) (*bool, error) {
	client, err := r.CreateClient()

	if err != nil {
		r.logger.Println("Error getting user: ", err)
		return nil, err
	}

	response, err := client.GetUserWithResponse(ctx, userId)

	if err != nil {
		r.logger.Println("Error getting user: ", err)
		return nil, err
	}

	r.logger.Println("Response for ID response: ", response.StatusCode(), " Id: ", userId)

	isUserSuspicious := response.StatusCode() != 404

	// 404 is ok - anything else is an error condition
	if response.StatusCode() != 404 && (response.StatusCode() >= 300 || response.StatusCode() < 200) {
		r.logger.Println("Error getting user: ", response.StatusCode())
		return nil, err
	}

	return &isUserSuspicious, nil
}
