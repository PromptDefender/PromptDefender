package sagemaker_jailbreak_model

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sagemakerruntime"
	"github.com/safetorun/PromptDefender/wall"
	"log"
)

type RemoteSagemakerCaller struct {
	endpointUrl string
	logger      *log.Logger
}

type InjectionResponse []struct {
	Label string  `json:"label"`
	Score float64 `json:"score"`
}
type Payload struct {
	Inputs       string `json:"inputs"`
	WaitForModel bool   `json:"wait_for_model"`
}

func New(endpointUrl string) RemoteSagemakerCaller {
	return RemoteSagemakerCaller{
		logger:      log.Default(),
		endpointUrl: endpointUrl,
	}
}

func (c *RemoteSagemakerCaller) Query(payload Payload) (*float64, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:     aws.String("eu-west-1"),
		MaxRetries: aws.Int(3),
	})
	if err != nil {
		return nil, err
	}
	sm := sagemakerruntime.New(sess)

	inputBytes, err := json.Marshal(payload)

	input := &sagemakerruntime.InvokeEndpointInput{
		Body:         inputBytes,
		EndpointName: aws.String(c.endpointUrl),
		ContentType:  aws.String("application/json"),
	}

	result, err := sm.InvokeEndpoint(input)

	if err != nil {
		return nil, err
	}

	var response InjectionResponse
	err = json.Unmarshal(result.Body, &response)
	if err != nil {
		return nil, err
	}

	for _, re := range response {
		if re.Label == "INJECTION" {
			return &re.Score, nil
		} else if re.Label == "LEGIT" {
			score := 1.0 - re.Score
			return &score, nil
		}
	}

	return nil, errors.New(
		fmt.Sprintf("Could not find the label in the response. Response %+v", result),
	)
}

func (r *RemoteSagemakerCaller) CallRemoteApi(prompt string) (wall.MatchLevel, error) {
	response, err := r.Query(Payload{Inputs: prompt, WaitForModel: true})

	if err != nil {
		r.logger.Println("Failed to call the remote API. Error: ", err)
		return -1, err
	}

	r.logger.Println("Prompt is ", prompt, " Injection score is ", *response)

	return matchLevelForScore(*response), nil
}

func matchLevelForScore(score float64) wall.MatchLevel {
	if score == 1.0 {
		return wall.ExactMatch
	}
	if score > 0.9 {
		return wall.VeryClose
	}
	if score > 0.5 {
		return wall.Medium
	}

	return wall.NoMatch
}
