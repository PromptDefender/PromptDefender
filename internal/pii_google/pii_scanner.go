package pii_google

import (
	"context"
	"fmt"

	dlp "cloud.google.com/go/dlp/apiv2"
	"cloud.google.com/go/dlp/apiv2/dlppb"
	"github.com/safetorun/PromptDefender/pii"
)

type GooglePIIScanner struct {
	Client    *dlp.Client
	ProjectID string
}

func New(projectID string) (*GooglePIIScanner, error) {
	ctx := context.Background()
	client, err := dlp.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	return &GooglePIIScanner{
		Client:    client,
		ProjectID: projectID,
	}, nil
}

func (s *GooglePIIScanner) Scan(input string) (*pii.ScanResult, error) {
	ctx := context.Background()
	req := &dlppb.InspectContentRequest{
		Parent: fmt.Sprintf("projects/%s/locations/global", s.ProjectID),
		Item: &dlppb.ContentItem{
			DataItem: &dlppb.ContentItem_Value{
				Value: input,
			},
		},
		InspectConfig: &dlppb.InspectConfig{
			InfoTypes: []*dlppb.InfoType{
				{Name: "EMAIL_ADDRESS"},
				{Name: "PHONE_NUMBER"},
				{Name: "CREDIT_CARD_NUMBER"},
				{Name: "US_SOCIAL_SECURITY_NUMBER"},
				// Add more as needed matching AWS Comprehend coverage roughly
			},
			MinLikelihood: dlppb.Likelihood_LIKELY,
			IncludeQuote:  true,
		},
	}

	resp, err := s.Client.InspectContent(ctx, req)
	if err != nil {
		return nil, err
	}

	return &pii.ScanResult{
		ContainingPii: len(resp.Result.Findings) > 0,
	}, nil
}

func (s *GooglePIIScanner) Close() error {
	return s.Client.Close()
}
