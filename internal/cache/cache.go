package cache

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type Cache interface {
	Get(key string) (*string, error)
	Set(key string, value string) error
}

type DdbCache struct {
	TableName string
	ddb       *dynamodb.DynamoDB
}

func New(tableName string) Cache {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1"),
	}))

	return DdbCache{
		TableName: tableName,
		ddb:       dynamodb.New(sess),
	}
}

func (c DdbCache) Get(key string) (*string, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(c.TableName),
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: aws.String(key),
			},
		},
	}

	result, err := c.ddb.GetItem(input)
	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		return nil, nil
	}

	return result.Item["Value"].S, nil
}

func (c DdbCache) Set(key string, value string) error {
	input := &dynamodb.PutItemInput{
		TableName: aws.String(c.TableName),
		Item: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: aws.String(key),
			},
			"Value": {
				S: aws.String(value),
			},
		},
	}

	_, err := c.ddb.PutItem(input)
	if err != nil {
		return err
	}

	return nil
}
