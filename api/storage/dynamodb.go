package storage

import "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"

type DynamoDB struct {
	Db dynamodbiface.DynamoDBAPI
}
