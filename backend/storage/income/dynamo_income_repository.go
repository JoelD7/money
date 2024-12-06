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
	dynamoClient          *dynamodb.Client
	tableName             string
	periodUserIncomeIndex string
}

func NewDynamoRepository(dynamoClient *dynamodb.Client, tableName string, periodUserIndex string) (*DynamoRepository, error) {
	d := &DynamoRepository{dynamoClient: dynamoClient}

	tableNameEnv := env.GetString("INCOME_TABLE_NAME", "")
	if tableName == "" && tableNameEnv == "" {
		return nil, fmt.Errorf("initialize income dynamo repository failed: table name is required")
	}

	periodUserEnv := env.GetString("PERIOD_USER_INCOME_INDEX", "")
	if periodUserIndex == "" && periodUserEnv == "" {
		return nil, fmt.Errorf("initialize income dynamo repository failed: period user index is required")
	}

	d.tableName = tableName
	if d.tableName == "" {
		d.tableName = tableNameEnv
	}

	d.periodUserIncomeIndex = periodUserIndex
	if d.periodUserIncomeIndex == "" {
		d.periodUserIncomeIndex = periodUserEnv
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

func (d *DynamoRepository) BatchCreateIncome(ctx context.Context, incomes []*models.Income) error {
	incomeEntities := make([]incomeEntity, 0, len(incomes))

	for _, income := range incomes {
		incomeEnt := toIncomeEntity(income)
		incomeEnt.PeriodUser = dynamo.BuildPeriodUser(income.Username, *income.Period)
		incomeEntities = append(incomeEntities, *incomeEnt)
	}

	writeRequests := make([]types.WriteRequest, 0, len(incomeEntities))

	for _, incomeEnt := range incomeEntities {
		incomeAv, err := attributevalue.MarshalMap(incomeEnt)
		if err != nil {
			return fmt.Errorf("marshal income attribute value failed: %v", err)
		}

		writeRequests = append(writeRequests, types.WriteRequest{
			PutRequest: &types.PutRequest{
				Item: incomeAv,
			},
		})
	}

	input := &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]types.WriteRequest{
			d.tableName: writeRequests,
		},
	}

	err := dynamo.BatchWrite(ctx, d.dynamoClient, input)
	if err != nil {
		return err
	}

	return nil
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

	incomeEnt := incomeEntity{}

	err = attributevalue.UnmarshalMap(result.Item, incomeEnt)
	if err != nil {
		return nil, err
	}

	return toIncomeModel(incomeEnt), nil
}

func (d *DynamoRepository) GetIncomeByPeriod(ctx context.Context, username string, params *models.QueryParameters) ([]*models.Income, string, error) {
	var decodedStartKey map[string]types.AttributeValue
	var err error

	if params.StartKey != "" {
		decodedStartKey, err = dynamo.DecodePaginationKey(params.StartKey)
		if err != nil {
			return nil, "", fmt.Errorf("%v: %w", err, models.ErrInvalidStartKey)
		}
	}

	periodUser := dynamo.BuildPeriodUser(username, params.Period)

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
		Limit:                     getPageSize(params.PageSize),
		ExclusiveStartKey:         decodedStartKey,
	}

	result, err := d.dynamoClient.Query(ctx, input)
	if err != nil {
		return nil, "", err
	}

	if (result.Items == nil || len(result.Items) == 0) && params.StartKey == "" {
		return nil, "", models.ErrIncomeNotFound
	}

	if result.Items == nil || len(result.Items) == 0 {
		return nil, "", models.ErrNoMoreItemsToBeRetrieved
	}

	incomeEntities := make([]incomeEntity, 0, len(result.Items))

	err = attributevalue.UnmarshalListOfMaps(result.Items, &incomeEntities)
	if err != nil {
		return nil, "", err
	}

	nextKey, err := dynamo.EncodePaginationKey(result.LastEvaluatedKey)
	if err != nil {
		return nil, "", err
	}

	return toIncomeModels(incomeEntities), nextKey, nil
}

func (d *DynamoRepository) GetAllIncome(ctx context.Context, username string, params *models.QueryParameters) ([]*models.Income, string, error) {
	var decodedStartKey map[string]types.AttributeValue
	var err error

	if params.StartKey != "" {
		decodedStartKey, err = dynamo.DecodePaginationKey(params.StartKey)
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
		Limit:                     getPageSize(params.PageSize),
		ExclusiveStartKey:         decodedStartKey,
	}

	result, err := d.dynamoClient.Query(ctx, input)
	if err != nil {
		return nil, "", err
	}

	if result.Items == nil || len(result.Items) == 0 && params.StartKey == "" {
		return nil, "", models.ErrIncomeNotFound
	}

	if result.Items == nil || len(result.Items) == 0 {
		return nil, "", models.ErrNoMoreItemsToBeRetrieved
	}

	incomeEntities := make([]incomeEntity, 0, len(result.Items))

	err = attributevalue.UnmarshalListOfMaps(result.Items, &incomeEntities)
	if err != nil {
		return nil, "", err
	}

	nextKey, err := dynamo.EncodePaginationKey(result.LastEvaluatedKey)
	if err != nil {
		return nil, "", err
	}

	return toIncomeModels(incomeEntities), nextKey, nil
}

func (d *DynamoRepository) GetAllIncomeByPeriod(ctx context.Context, username string, params *models.QueryParameters) ([]*models.Income, error) {
	periodUser := dynamo.BuildPeriodUser(username, params.Period)
	periodUserCond := expression.Key("period_user").Equal(expression.Value(periodUser))

	expr, err := expression.NewBuilder().WithKeyCondition(periodUserCond).Build()
	if err != nil {
		return nil, err
	}

	input := &dynamodb.QueryInput{
		TableName:                 aws.String(d.tableName),
		IndexName:                 aws.String(d.periodUserIncomeIndex),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
	}

	var result *dynamodb.QueryOutput
	entities := make([]incomeEntity, 0)
	var itemsInQuery []incomeEntity

	for {
		itemsInQuery = make([]incomeEntity, 0)
		result, err = d.dynamoClient.Query(ctx, input)
		if err != nil {
			return nil, err
		}

		if (result.Items == nil || len(result.Items) == 0) && result.LastEvaluatedKey == nil {
			break
		}

		err = attributevalue.UnmarshalListOfMaps(result.Items, &itemsInQuery)
		if err != nil {
			return nil, fmt.Errorf("unmarshal income items failed: %v", err)
		}

		entities = append(entities, itemsInQuery...)
		input.ExclusiveStartKey = result.LastEvaluatedKey

		if result.LastEvaluatedKey == nil {
			break
		}
	}

	if len(entities) == 0 {
		return nil, models.ErrIncomeNotFound
	}

	return toIncomeModels(entities), nil
}

func (d *DynamoRepository) GetAllIncomePeriods(ctx context.Context, username string) ([]string, error) {
	keyCond := expression.Key("username").Equal(expression.Value(username))
	projection := expression.NamesList(expression.Name("period"))

	expr, err := expression.NewBuilder().WithKeyCondition(keyCond).WithProjection(projection).Build()
	if err != nil {
		return nil, fmt.Errorf("build expression failed: %v", err)
	}

	input := &dynamodb.QueryInput{
		TableName:                 aws.String(d.tableName),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		ProjectionExpression:      expr.Projection(),
	}

	var result *dynamodb.QueryOutput
	entities := make([]incomeEntity, 0)
	var itemsInQuery []incomeEntity

	for {
		itemsInQuery = make([]incomeEntity, 0)
		result, err = d.dynamoClient.Query(ctx, input)
		if err != nil {
			return nil, fmt.Errorf("query income periods failed: %v", err)
		}

		if (result.Items == nil || len(result.Items) == 0) && result.LastEvaluatedKey == nil {
			break
		}

		err = attributevalue.UnmarshalListOfMaps(result.Items, &itemsInQuery)
		if err != nil {
			return nil, fmt.Errorf("unmarshal income items failed: %v", err)
		}

		entities = append(entities, itemsInQuery...)
		input.ExclusiveStartKey = result.LastEvaluatedKey

		if result.LastEvaluatedKey == nil {
			break
		}
	}

	if len(entities) == 0 {
		return nil, models.ErrIncomeNotFound
	}

	periods := make([]string, 0, len(entities))
	existsMap := make(map[string]struct{})
	var exists bool

	for _, entity := range entities {
		if entity.Period == nil {
			continue
		}

		_, exists = existsMap[*entity.Period]
		if !exists {
			periods = append(periods, *entity.Period)
			existsMap[*entity.Period] = struct{}{}
		}
	}

	return periods, nil
}

func (d *DynamoRepository) BatchDeleteIncome(ctx context.Context, income []*models.Income) error {
	writeRequests := make([]types.WriteRequest, 0, len(income))

	var usernameAttrValue types.AttributeValue
	var incomeIDAttrValue types.AttributeValue
	var err error

	for _, in := range income {
		usernameAttrValue, err = attributevalue.Marshal(in.Username)
		if err != nil {
			return fmt.Errorf("marshal username key failed: %v", err)
		}

		incomeIDAttrValue, err = attributevalue.Marshal(in.IncomeID)
		if err != nil {
			return fmt.Errorf("marshal id key failed: %v", err)
		}

		writeRequests = append(writeRequests, types.WriteRequest{
			DeleteRequest: &types.DeleteRequest{
				Key: map[string]types.AttributeValue{
					"username":  usernameAttrValue,
					"income_id": incomeIDAttrValue,
				},
			},
		})
	}

	input := &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]types.WriteRequest{
			d.tableName: writeRequests,
		},
	}

	return dynamo.BatchWrite(ctx, d.dynamoClient, input)
}

func getPageSize(pageSize int) *int32 {
	if pageSize == 0 {
		return aws.Int32(defaultPageSize)
	}

	return aws.Int32(int32(pageSize))
}
