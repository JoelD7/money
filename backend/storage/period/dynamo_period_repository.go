package period

import (
	"context"
	"fmt"
	"github.com/JoelD7/money/backend/models"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
)

const (
	tableName = "period"
)

type DynamoRepository struct {
	dynamoClient *dynamodb.Client
}

func NewDynamoRepository(dynamoClient *dynamodb.Client) *DynamoRepository {
	return &DynamoRepository{dynamoClient: dynamoClient}
}

func (d *DynamoRepository) CreatePeriod(ctx context.Context, period *models.Period) (*models.Period, error) {
	periodStruct := toPeriodEntity(*period)

	attrValue, err := attributevalue.MarshalMap(periodStruct)
	if err != nil {
		return nil, fmt.Errorf("marshal period item failed: %v", err)
	}

	input := &dynamodb.PutItemInput{
		Item:      attrValue,
		TableName: aws.String(tableName),
	}

	_, err = d.dynamoClient.PutItem(ctx, input)
	if err != nil {
		return nil, err
	}

	return period, nil
}

func (d *DynamoRepository) GetPeriod(ctx context.Context, username, period string) (*models.Period, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"username": &types.AttributeValueMemberS{Value: username},
			"period":   &types.AttributeValueMemberS{Value: period},
		},
	}

	result, err := d.dynamoClient.GetItem(ctx, input)
	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		return nil, models.ErrPeriodNotFound
	}

	periodStruct := periodEntity{}

	err = attributevalue.UnmarshalMap(result.Item, &periodStruct)
	if err != nil {
		return nil, fmt.Errorf("unmarshal period item failed: %v", err)
	}

	return toPeriodModel(periodStruct), nil
}

func (d *DynamoRepository) GetPeriods(ctx context.Context, username string) ([]*models.Period, error) {
	keyConditionExpression := expression.Key("username").Equal(expression.Value(username))

	conditionBuilder := expression.NewBuilder().WithKeyCondition(keyConditionExpression)

	expr, err := conditionBuilder.Build()
	if err != nil {
		return nil, fmt.Errorf("build expression failed: %v", err)
	}

	input := &dynamodb.QueryInput{
		TableName:                 aws.String(tableName),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeValues: expr.Values(),
	}

	result, err := d.dynamoClient.Query(ctx, input)
	if err != nil {
		return nil, err
	}

	if result.Items == nil || len(result.Items) == 0 {
		return nil, models.ErrPeriodsNotFound
	}

	periods := make([]periodEntity, 0, len(result.Items))

	err = attributevalue.UnmarshalListOfMaps(result.Items, &periods)
	if err != nil {
		return nil, fmt.Errorf("unmarshal periods failed: %v", err)
	}

	return toPeriodModels(periods), nil
}

func (d *DynamoRepository) GetLastPeriod(ctx context.Context, username string) (*models.Period, error) {
	keyConditionExpression := expression.Key("username").Equal(expression.Value(username))

	conditionBuilder := expression.NewBuilder().WithKeyCondition(keyConditionExpression)

	expr, err := conditionBuilder.Build()
	if err != nil {
		return nil, fmt.Errorf("build expression failed: %v", err)
	}

	input := &dynamodb.QueryInput{
		TableName:                 aws.String(tableName),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ScanIndexForward:          aws.Bool(false),
	}

	result, err := d.dynamoClient.Query(ctx, input)
	if err != nil {
		return nil, err
	}

	if result.Items == nil || len(result.Items) == 0 {
		return nil, models.ErrPeriodsNotFound
	}

	periodStruct := periodEntity{}

	err = attributevalue.UnmarshalMap(result.Items[0], &periodStruct)
	if err != nil {
		return nil, fmt.Errorf("unmarshal period item failed: %v", err)
	}

	return toPeriodModel(periodStruct), nil
}
