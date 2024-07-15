package income

import (
	"context"
	"fmt"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/env"
	"github.com/JoelD7/money/backend/storage/dynamo"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"strings"
)

const (
	defaultPageSize   = 10
	conditionFailedEx = "ConditionalCheckFailedException"
)

var (
	periodUserIncomeIDIndex = "period_user-income_id-index"
)

type DynamoRepository struct {
	dynamoClient *dynamodb.Client
	tableName    string
}

func NewDynamoRepository(dynamoClient *dynamodb.Client, tableName string) (*DynamoRepository, error) {
	d := &DynamoRepository{dynamoClient: dynamoClient}

	tableNameEnv := env.GetString("INCOME_TABLE_NAME", "")
	if tableName == "" && tableNameEnv == "" {
		return nil, fmt.Errorf("initialize income recurring dynamo repository failed: table name is required")
	}

	d.tableName = tableName
	if d.tableName == "" {
		d.tableName = tableNameEnv
	}

	return d, nil
}

func (d *DynamoRepository) CreateIncome(ctx context.Context, income *models.Income) (*models.Income, error) {
	incomeEnt := toIncomeEntity(income)
	incomeEnt.PeriodUser = dynamo.BuildPeriodUser(income.Username, *income.Period)

	incomeAv, err := attributevalue.MarshalMap(incomeEnt)
	if err != nil {
		return nil, fmt.Errorf("marshal income attribute value failed: %v", err)
	}

	cond := expression.Name("income_id").AttributeNotExists()

	expr, err := expression.NewBuilder().WithCondition(cond).Build()
	if err != nil {
		return nil, fmt.Errorf("build expression failed: %v", err)
	}

	input := &dynamodb.PutItemInput{
		TableName:                aws.String(d.tableName),
		Item:                     incomeAv,
		ConditionExpression:      expr.Condition(),
		ExpressionAttributeNames: expr.Names(),
	}

	_, err = d.dynamoClient.PutItem(ctx, input)
	if err != nil && strings.Contains(err.Error(), conditionFailedEx) {
		return nil, fmt.Errorf("%v: %w", err, models.ErrExistingIncome)
	}

	if err != nil {
		return nil, err
	}

	return income, nil
}

func (d *DynamoRepository) GetIncome(ctx context.Context, username, incomeID string) (*models.Income, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(d.tableName),
		Key: map[string]types.AttributeValue{
			"income_id": &types.AttributeValueMemberS{Value: incomeID},
			"username":  &types.AttributeValueMemberS{Value: username},
		},
	}

	result, err := d.dynamoClient.GetItem(ctx, input)
	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		return nil, models.ErrIncomeNotFound
	}

	incomeEnt := new(incomeEntity)

	err = attributevalue.UnmarshalMap(result.Item, incomeEnt)
	if err != nil {
		return nil, err
	}

	return toIncomeModel(incomeEnt), nil
}

func (d *DynamoRepository) GetIncomeByPeriod(ctx context.Context, username, periodID, startKey string, pageSize int) ([]*models.Income, string, error) {
	var decodedStartKey map[string]types.AttributeValue
	var err error

	if startKey != "" {
		decodedStartKey, err = dynamo.DecodePaginationKey(startKey, &keysPeriodUserIndex{})
		if err != nil {
			return nil, "", fmt.Errorf("%v: %w", err, models.ErrInvalidStartKey)
		}
	}

	periodUser := dynamo.BuildPeriodUser(username, periodID)

	nameEx := expression.Name("period_user").Equal(expression.Value(periodUser))

	expr, err := expression.NewBuilder().WithCondition(nameEx).Build()
	if err != nil {
		return nil, "", err
	}

	input := &dynamodb.QueryInput{
		TableName:                 aws.String(d.tableName),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.Condition(),
		IndexName:                 aws.String(periodUserIncomeIDIndex),
		Limit:                     getPageSize(pageSize),
		ExclusiveStartKey:         decodedStartKey,
	}

	result, err := d.dynamoClient.Query(ctx, input)
	if err != nil {
		return nil, "", err
	}

	if (result.Items == nil || len(result.Items) == 0) && startKey == "" {
		return nil, "", models.ErrIncomeNotFound
	}

	if result.Items == nil || len(result.Items) == 0 {
		return nil, "", models.ErrNoMoreItemsToBeRetrieved
	}

	incomeEntities := new([]*incomeEntity)

	err = attributevalue.UnmarshalListOfMaps(result.Items, &incomeEntities)
	if err != nil {
		return nil, "", err
	}

	nextKey, err := dynamo.EncodePaginationKey(result.LastEvaluatedKey, &keysPeriodUserIndex{})
	if err != nil {
		return nil, "", err
	}

	return toIncomeModels(*incomeEntities), nextKey, nil
}

func (d *DynamoRepository) GetAllIncome(ctx context.Context, username, startKey string, pageSize int) ([]*models.Income, string, error) {
	var decodedStartKey map[string]types.AttributeValue
	var err error

	if startKey != "" {
		decodedStartKey, err = dynamo.DecodePaginationKey(startKey, &keys{})
		if err != nil {
			return nil, "", fmt.Errorf("%v: %w", err, models.ErrInvalidStartKey)
		}
	}

	nameEx := expression.Name("username").Equal(expression.Value(username))

	expr, err := expression.NewBuilder().WithCondition(nameEx).Build()
	if err != nil {
		return nil, "", err
	}

	input := &dynamodb.QueryInput{
		TableName:                 aws.String(d.tableName),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.Condition(),
		Limit:                     getPageSize(pageSize),
		ExclusiveStartKey:         decodedStartKey,
	}

	result, err := d.dynamoClient.Query(ctx, input)
	if err != nil {
		return nil, "", err
	}

	if result.Items == nil || len(result.Items) == 0 && startKey == "" {
		return nil, "", models.ErrIncomeNotFound
	}

	if result.Items == nil || len(result.Items) == 0 {
		return nil, "", models.ErrNoMoreItemsToBeRetrieved
	}

	incomeEntities := make([]*incomeEntity, 0, len(result.Items))

	err = attributevalue.UnmarshalListOfMaps(result.Items, &incomeEntities)
	if err != nil {
		return nil, "", err
	}

	nextKey, err := dynamo.EncodePaginationKey(result.LastEvaluatedKey, &keys{})
	if err != nil {
		return nil, "", err
	}

	return toIncomeModels(incomeEntities), nextKey, nil
}

func getPageSize(pageSize int) *int32 {
	if pageSize == 0 {
		return aws.Int32(defaultPageSize)
	}

	return aws.Int32(int32(pageSize))
}
