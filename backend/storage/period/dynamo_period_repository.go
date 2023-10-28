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
	"strings"
)

const (
	tableName       = "period"
	defaultPageSize = 10
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
		Item:                attrValue,
		ConditionExpression: aws.String("attribute_not_exists(period)"),
		TableName:           aws.String(tableName),
	}

	_, err = d.dynamoClient.PutItem(ctx, input)
	if err != nil && strings.Contains(err.Error(), "ConditionalCheckFailedException") {
		return nil, fmt.Errorf("%v: %w", err, models.ErrPeriodExists)
	}

	if err != nil {
		return nil, err
	}

	return period, nil
}

func (d *DynamoRepository) UpdatePeriod(ctx context.Context, period *models.Period) error {
	periodEnt := toPeriodEntity(*period)

	periodAv, err := attributevalue.MarshalMap(periodEnt)
	if err != nil {
		return fmt.Errorf("marshaling period to attribute value: %v", err)
	}

	input := &dynamodb.PutItemInput{
		TableName:           aws.String(tableName),
		ConditionExpression: aws.String("attribute_exists(period)"),
		Item:                periodAv,
	}

	_, err = d.dynamoClient.PutItem(ctx, input)
	if err != nil && strings.Contains(err.Error(), "ConditionalCheckFailedException") {
		return fmt.Errorf("%v: %w", err, models.ErrUpdatePeriodNotFound)
	}

	if err != nil {
		return fmt.Errorf("updating period item: %v", err)
	}

	return nil
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

func (d *DynamoRepository) GetPeriods(ctx context.Context, username, startKey string, pageSize int) ([]*models.Period, error) {
	keyConditionExpression := expression.Key("username").Equal(expression.Value(username))

	conditionBuilder := expression.NewBuilder().WithKeyCondition(keyConditionExpression)

	expr, err := conditionBuilder.Build()
	if err != nil {
		return nil, fmt.Errorf("build expression failed: %v", err)
	}

	var decodedStartKey map[string]types.AttributeValue

	if startKey != "" {
		decodedStartKey, err = decodeStartKey(startKey)
		if err != nil {
			return nil, fmt.Errorf("%v: %w", err, models.ErrInvalidStartKey)
		}
	}

	input := &dynamodb.QueryInput{
		TableName:                 aws.String(tableName),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ExclusiveStartKey:         decodedStartKey,
		Limit:                     getPageSize(pageSize),
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

func getPageSize(pageSize int) *int32 {
	if pageSize == 0 {
		return aws.Int32(defaultPageSize)
	}

	return aws.Int32(int32(pageSize))
}
