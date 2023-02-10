package storage

import "github.com/gusaul/go-dynamock"

var DynamoMock *dynamock.DynaMock

func InitDynamoMock() {
	Dynamo = new(DynamoDB)
	Dynamo.Db, DynamoMock = dynamock.New()
}
