package savings

import (
	"context"
	"fmt"
	"github.com/JoelD7/money/backend/models"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go/aws"
)

const (
	tableName = "savings"
)

type DynamoRepository struct {
	dynamoClient *dynamodb.Client
}

func NewDynamoRepository(dynamoClient *dynamodb.Client) *DynamoRepository {
	return &DynamoRepository{dynamoClient: dynamoClient}
}

func (d *DynamoRepository) GetSavings(ctx context.Context, email string) ([]*models.Saving, error) {
	nameEx := expression.Name("email").Equal(expression.Value(email))

	expr, err := expression.NewBuilder().WithCondition(nameEx).Build()
	if err != nil {
		return nil, err
	}

	input := &dynamodb.QueryInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.Condition(),
		TableName:                 aws.String(tableName),
	}

	result, err := d.dynamoClient.Query(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("query failed: %v", err)
	}

	if result.Items == nil || len(result.Items) == 0 {
		return nil, models.ErrSavingsNotFound
	}

	savings := new([]*models.Saving)

	err = attributevalue.UnmarshalListOfMaps(result.Items, savings)
	if err != nil {
		return nil, fmt.Errorf("unmarshal savings items failed: %v", err)
	}

	return *savings, nil
}

func (d *DynamoRepository) CreateSaving(ctx context.Context, saving *models.Saving) error {
	item, err := attributevalue.MarshalMap(saving)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(tableName),
	}

	_, err = d.dynamoClient.PutItem(ctx, input)
	if err != nil {
		return fmt.Errorf("put saving item failed: %v", err)
	}

	return nil
}
