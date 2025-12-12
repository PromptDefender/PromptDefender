package pii_aws

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/comprehend"
	pii_scanner "github.com/safetorun/PromptDefender/pii"
)

type AwsPIIScanner struct {
}

func New() pii_scanner.Scanner {
	return AwsPIIScanner{}
}

func (s AwsPIIScanner) Scan(textToScan string) (*pii_scanner.ScanResult, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1"),
	})

	if err != nil {
		fmt.Println("Error creating session:", err)
		return nil, err
	}

	svc := comprehend.New(sess)
	input := &comprehend.DetectPiiEntitiesInput{
		Text:         aws.String(textToScan),
		LanguageCode: aws.String("en"),
	}

	output, err := svc.DetectPiiEntities(input)

	if err != nil {
		fmt.Println("Error detecting PII:", err)
		return nil, err
	}

	for _, entity := range output.Entities {
		fmt.Println("Entity:", *entity.Type, "Score:", *entity.Score)
	}

	return &pii_scanner.ScanResult{ContainingPii: len(output.Entities) > 0}, nil
}
