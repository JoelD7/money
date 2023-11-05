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
	"time"
)

const (
	periodTableName           = "period"
	uniquePeriodNameTableName = "unique-period-name"
	defaultPageSize           = 10
	conditionalFailedKeyword  = "ConditionalCheckFailed"
)

type DynamoRepository struct {
	dynamoClient *dynamodb.Client
}

func NewDynamoRepository(dynamoClient *dynamodb.Client) *DynamoRepository {
	return &DynamoRepository{dynamoClient: dynamoClient}
}

func (d *DynamoRepository) CreatePeriod(ctx context.Context, period *models.Period) (*models.Period, error) {
	periodEnt := toPeriodEntity(*period)
	uPeriodName := &uniquePeriodNameEntity{
		Name:     *periodEnt.Name,
		Username: periodEnt.Username,
	}

	attrValue, err := attributevalue.MarshalMap(periodEnt)
	if err != nil {
		return nil, fmt.Errorf("marshal period item failed: %v", err)
	}

	uPeriodNameAttrValue, err := attributevalue.MarshalMap(uPeriodName)
	if err != nil {
		return nil, fmt.Errorf("marshal unique period name item failed: %v", err)
	}

	condExpr := expression.Name("name").AttributeNotExists().And(expression.Name("username").AttributeNotExists())

	expr, err := expression.NewBuilder().WithCondition(condExpr).Build()
	if err != nil {
		return nil, fmt.Errorf("build expression failed: %v", err)
	}

	input := &dynamodb.TransactWriteItemsInput{
		TransactItems: []types.TransactWriteItem{
			{
				Put: &types.Put{
					Item:      attrValue,
					TableName: aws.String(periodTableName),
				},
			},
			{
				Put: &types.Put{
					Item:                     uPeriodNameAttrValue,
					TableName:                aws.String(uniquePeriodNameTableName),
					ConditionExpression:      expr.Condition(),
					ExpressionAttributeNames: expr.Names(),
				},
			},
		},
	}

	_, err = d.dynamoClient.TransactWriteItems(ctx, input)
	if err != nil && strings.Contains(err.Error(), conditionalFailedKeyword) {
		return nil, fmt.Errorf("%v: %w", err, models.ErrPeriodNameIsTaken)
	}

	if err != nil {
		return nil, err
	}

	return period, nil
}

func (d *DynamoRepository) UpdatePeriod(ctx context.Context, period *models.Period) error {
	periodEnt := toPeriodEntity(*period)

	periodEnt.UpdatedDate = time.Now()

	uPeriodName := &uniquePeriodNameEntity{
		Name:     *periodEnt.Name,
		Username: periodEnt.Username,
	}

	periodAv, err := attributevalue.MarshalMap(periodEnt)
	if err != nil {
		return fmt.Errorf("marshaling period to attribute value: %v", err)
	}

	uPeriodNameAv, err := attributevalue.MarshalMap(uPeriodName)
	if err != nil {
		return fmt.Errorf("marshaling unique period name to attribute value failed: %v", err)
	}

	periodExistsCond := expression.Name("period").AttributeExists()
	periodNameNotTakenCond := expression.Name("name").AttributeNotExists().And(expression.Name("username").AttributeNotExists())

	periodTableExpr, err := expression.NewBuilder().WithCondition(periodExistsCond).Build()
	if err != nil {
		return fmt.Errorf("building period table expression failed: %v", err)
	}

	uniquePeriodNameTableExpr, err := expression.NewBuilder().WithCondition(periodNameNotTakenCond).Build()
	if err != nil {
		return fmt.Errorf("building unique period name table expression failed: %v", err)
	}

	input := &dynamodb.TransactWriteItemsInput{
		TransactItems: []types.TransactWriteItem{
			{
				Put: &types.Put{
					TableName:                aws.String(periodTableName),
					ConditionExpression:      periodTableExpr.Condition(),
					ExpressionAttributeNames: periodTableExpr.Names(),
					Item:                     periodAv,
				},
			},
			{
				Put: &types.Put{
					TableName:                aws.String(uniquePeriodNameTableName),
					ConditionExpression:      uniquePeriodNameTableExpr.Condition(),
					ExpressionAttributeNames: uniquePeriodNameTableExpr.Names(),
					Item:                     uPeriodNameAv,
				},
			},
		},
	}

	_, err = d.dynamoClient.TransactWriteItems(ctx, input)
	if err != nil {
		return handleUpdatePeriodError(err)
	}

	return nil
}

func handleUpdatePeriodError(err error) error {
	periodTableConditionFailed := fmt.Sprintf("[%s, None]", conditionalFailedKeyword)
	uniquePeriodNameTableConditionFailed := fmt.Sprintf("[None, %s]", conditionalFailedKeyword)

	if strings.Contains(err.Error(), periodTableConditionFailed) {
		return fmt.Errorf("%v: %w", err, models.ErrUpdatePeriodNotFound)
	}

	if strings.Contains(err.Error(), uniquePeriodNameTableConditionFailed) {
		return fmt.Errorf("%v: %w", err, models.ErrPeriodNameIsTaken)
	}

	return fmt.Errorf("updating period item: %v", err)
}

func (d *DynamoRepository) GetPeriod(ctx context.Context, username, period string) (*models.Period, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(periodTableName),
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

func (d *DynamoRepository) GetPeriods(ctx context.Context, username, startKey string, pageSize int) ([]*models.Period, string, error) {
	keyConditionExpression := expression.Key("username").Equal(expression.Value(username))

	conditionBuilder := expression.NewBuilder().WithKeyCondition(keyConditionExpression)

	expr, err := conditionBuilder.Build()
	if err != nil {
		return nil, "", fmt.Errorf("build expression failed: %v", err)
	}

	var decodedStartKey map[string]types.AttributeValue

	if startKey != "" {
		decodedStartKey, err = decodeStartKey(startKey)
		if err != nil {
			return nil, "", fmt.Errorf("%v: %w", err, models.ErrInvalidStartKey)
		}
	}

	input := &dynamodb.QueryInput{
		TableName:                 aws.String(periodTableName),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ExclusiveStartKey:         decodedStartKey,
		Limit:                     getPageSize(pageSize),
	}

	result, err := d.dynamoClient.Query(ctx, input)
	if err != nil {
		return nil, "", err
	}

	if (result.Items == nil || len(result.Items) == 0) && startKey == "" {
		return nil, "", models.ErrPeriodsNotFound
	}

	periods := make([]periodEntity, 0, len(result.Items))

	err = attributevalue.UnmarshalListOfMaps(result.Items, &periods)
	if err != nil {
		return nil, "", fmt.Errorf("unmarshal periods failed: %v", err)
	}

	nextKey, err := encodeLastKey(result.LastEvaluatedKey)
	if err != nil {
		return nil, "", err
	}

	return toPeriodModels(periods), nextKey, nil
}

func (d *DynamoRepository) GetLastPeriod(ctx context.Context, username string) (*models.Period, error) {
	keyConditionExpression := expression.Key("username").Equal(expression.Value(username))

	conditionBuilder := expression.NewBuilder().WithKeyCondition(keyConditionExpression)

	expr, err := conditionBuilder.Build()
	if err != nil {
		return nil, fmt.Errorf("build expression failed: %v", err)
	}

	input := &dynamodb.QueryInput{
		TableName:                 aws.String(periodTableName),
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
