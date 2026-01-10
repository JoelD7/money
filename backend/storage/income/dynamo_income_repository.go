package income

import (
	"context"
	"fmt"
	"github.com/JoelD7/money/backend/models"
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

type DynamoRepository struct {
	dynamoClient                *dynamodb.Client
	tableName                   string
	periodUserIncomeIndex       string
	periodUserCreatedDateIndex  string
	usernameCreatedDateIndex    string
	periodUserNameIncomeIDIndex string
	periodUserAmountIndex       string
}

func NewDynamoRepository(dynamoClient *dynamodb.Client, envConfig *models.EnvironmentConfiguration) (*DynamoRepository, error) {
	d := &DynamoRepository{dynamoClient: dynamoClient}

	err := validateParams(envConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize income dynamo repository: %v", err)
	}

	d.tableName = envConfig.IncomeTable
	d.periodUserIncomeIndex = envConfig.PeriodUserIncomeIndex
	d.periodUserCreatedDateIndex = envConfig.PeriodUserCreatedDateIndex
	d.usernameCreatedDateIndex = envConfig.UsernameCreatedDateIndex
	d.periodUserNameIncomeIDIndex = envConfig.PeriodUserNameIncomeIDIndex
	d.periodUserAmountIndex = envConfig.PeriodUserAmountIndex

	return d, nil
}

func validateParams(envConfig *models.EnvironmentConfiguration) error {
	if envConfig.IncomeTable == "" {
		return fmt.Errorf("income table name is required")
	}

	if envConfig.PeriodUserIncomeIndex == "" {
		return fmt.Errorf("period user income index is required")
	}

	if envConfig.PeriodUserCreatedDateIndex == "" {
		return fmt.Errorf("period user created date index is required")
	}

	if envConfig.UsernameCreatedDateIndex == "" {
		return fmt.Errorf("username created date index is required")
	}

	if envConfig.PeriodUserNameIncomeIDIndex == "" {
		return fmt.Errorf("period user name income id index is required")
	}

	if envConfig.PeriodUserAmountIndex == "" {
		return fmt.Errorf("period user amount index is required")
	}

	return nil
}

func (d *DynamoRepository) CreateIncome(ctx context.Context, income *models.Income) (*models.Income, error) {
	incomeEnt := toIncomeEntity(income)
	incomeEnt.PeriodUser = dynamo.BuildPeriodUser(income.Username, *income.PeriodID)

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
		incomeEnt.PeriodUser = dynamo.BuildPeriodUser(income.Username, *income.PeriodID)
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

func (d *DynamoRepository) GetIncomeByPeriod(ctx context.Context, username string, params *models.IncomeQueryParameters) ([]*models.Income, string, error) {
	input, err := d.buildQueryInput(username, params, nil)
	if err != nil {
		return nil, "", err
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

func (d *DynamoRepository) GetAllIncome(ctx context.Context, username string, params *models.IncomeQueryParameters) ([]*models.Income, string, error) {
	input, err := d.buildQueryInput(username, params, nil)
	if err != nil {
		return nil, "", err
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

func (d *DynamoRepository) GetAllIncomeByPeriod(ctx context.Context, username string, params *models.IncomeQueryParameters) ([]*models.Income, error) {
	input, err := d.buildQueryInput(username, params, nil)
	if err != nil {
		return nil, err
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
	input, err := d.buildQueryInput(username, nil, nil)
	if err != nil {
		return nil, err
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
		if entity.PeriodID == nil {
			continue
		}

		_, exists = existsMap[*entity.PeriodID]
		if !exists {
			periods = append(periods, *entity.PeriodID)
			existsMap[*entity.PeriodID] = struct{}{}
		}
	}

	return periods, nil
}

func (d *DynamoRepository) buildQueryInput(username string, params *models.IncomeQueryParameters, projection *expression.ProjectionBuilder) (*dynamodb.QueryInput, error) {
	if params == nil {
		params = &models.IncomeQueryParameters{}
	}

	input := &dynamodb.QueryInput{
		TableName: aws.String(d.tableName),
		Limit:     getPageSize(params.PageSize),
	}

	if params.SortType == string(models.SortOrderDescending) {
		input.ScanIndexForward = aws.Bool(false)
	}

	err := dynamo.SetExclusiveStartKey(params.StartKey, input)
	if err != nil {
		return nil, err
	}

	keyCondition := d.setQueryIndex(input, username, params)
	conditionBuilder := expression.NewBuilder().WithKeyCondition(keyCondition)

	if projection != nil {
		conditionBuilder.WithProjection(*projection)
	}

	expr, err := conditionBuilder.Build()
	if err != nil {
		return nil, err
	}

	input.ExpressionAttributeNames = expr.Names()
	input.ExpressionAttributeValues = expr.Values()
	input.KeyConditionExpression = expr.KeyCondition()
	input.ProjectionExpression = expr.Projection()

	return input, nil
}

// setQueryIndex sets the index to be used in the query based on the sorting and filter parameters. Returns a key
// condition expression formed with the index's primary key.
func (d *DynamoRepository) setQueryIndex(input *dynamodb.QueryInput, username string, params *models.IncomeQueryParameters) expression.KeyConditionBuilder {
	keyConditionEx := expression.Key("username").Equal(expression.Value(username))

	if params.Period != "" {
		input.IndexName = aws.String(d.periodUserIncomeIndex)

		periodUser := dynamo.BuildPeriodUser(username, params.Period)
		keyConditionEx = expression.Key("period_user").Equal(expression.Value(periodUser))
	}

	if params.SortBy == string(models.SortParamCreatedDate) {
		input.IndexName = aws.String(d.usernameCreatedDateIndex)
	}

	if params.Period != "" && params.SortBy == string(models.SortParamCreatedDate) {
		input.IndexName = aws.String(d.periodUserCreatedDateIndex)
	}

	if params.Period != "" && params.SortBy == string(models.SortParamAmount) {
		input.IndexName = aws.String(d.periodUserAmountIndex)
	}

	if params.Period != "" && params.SortBy == string(models.SortParamName) {
		input.IndexName = aws.String(d.periodUserNameIncomeIDIndex)
	}

	return keyConditionEx
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
