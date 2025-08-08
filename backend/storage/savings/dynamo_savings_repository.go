package savings

import (
	"context"
	"fmt"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/storage/dynamo"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"math"
	"strings"
	"time"
)

const (
	defaultPageSize = 10

	savingsPrefix = "SV"
)

type DynamoRepository struct {
	dynamoClient                     *dynamodb.Client
	tableName                        string
	periodSavingIndex                string
	savingGoalSavingIndex            string
	savingGoalCreatedDateSavingIndex string
	usernameAmountIndex              string
	usernameCreatedDateIndex         string
	periodUserCreatedDateIndex       string
	periodUserAmountIndex            string
}

func NewDynamoRepository(dynamoClient *dynamodb.Client, envConfig *models.EnvironmentConfiguration) (*DynamoRepository, error) {
	d := &DynamoRepository{dynamoClient: dynamoClient}

	err := validateParams(envConfig)
	if err != nil {
		return nil, fmt.Errorf("initialize saving dynamo repository failed: %v", err)
	}

	d.tableName = envConfig.SavingsTable
	d.periodSavingIndex = envConfig.PeriodSavingIndexName
	d.savingGoalSavingIndex = envConfig.SavingGoalSavingIndexName
	d.usernameAmountIndex = envConfig.UsernameAmountIndex
	d.usernameCreatedDateIndex = envConfig.UsernameCreatedDateIndex
	d.periodUserCreatedDateIndex = envConfig.PeriodUserCreatedDateIndex
	d.periodUserAmountIndex = envConfig.PeriodUserAmountIndex
	d.savingGoalCreatedDateSavingIndex = envConfig.SavingGoalCreatedDateSavingIndexName
	return d, nil
}

func validateParams(envConfig *models.EnvironmentConfiguration) error {
	if envConfig.SavingsTable == "" {
		return fmt.Errorf("table name is required")
	}

	if envConfig.PeriodSavingIndexName == "" {
		return fmt.Errorf("period saving index is required")
	}

	if envConfig.SavingGoalSavingIndexName == "" {
		return fmt.Errorf("saving goal saving index is required")
	}

	if envConfig.UsernameAmountIndex == "" {
		return fmt.Errorf("username amount index is required")
	}

	if envConfig.UsernameCreatedDateIndex == "" {
		return fmt.Errorf("username created date index is required")
	}

	if envConfig.PeriodUserCreatedDateIndex == "" {
		return fmt.Errorf("period user created date index is required")
	}

	if envConfig.PeriodUserAmountIndex == "" {
		return fmt.Errorf("period user amount index is required")
	}

	if envConfig.SavingGoalCreatedDateSavingIndexName == "" {
		return fmt.Errorf("saving goal created date saving index is required")
	}

	return nil
}

func (d *DynamoRepository) GetSaving(ctx context.Context, username, savingID string) (*models.Saving, error) {
	userKey, err := attributevalue.Marshal(username)
	if err != nil {
		return nil, err
	}

	savingIDKey, err := attributevalue.Marshal(savingID)
	if err != nil {
		return nil, err
	}

	input := &dynamodb.GetItemInput{
		TableName: aws.String(d.tableName),
		Key: map[string]types.AttributeValue{
			"username":  userKey,
			"saving_id": savingIDKey,
		},
	}

	result, err := d.dynamoClient.GetItem(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("get saving item failed: %v", err)
	}

	if result.Item == nil {
		return nil, models.ErrSavingNotFound
	}

	savingEnt := new(savingEntity)

	err = attributevalue.UnmarshalMap(result.Item, savingEnt)
	if err != nil {
		return nil, fmt.Errorf("unmarshal saving item failed: %v", err)
	}

	return toSavingModel(*savingEnt), nil
}

func (d *DynamoRepository) GetSavings(ctx context.Context, username string, params *models.QueryParameters) ([]*models.Saving, string, error) {
	input, err := d.buildQueryInput(username, params)
	if err != nil {
		return nil, "", fmt.Errorf("building query input: %v", err)
	}

	result, err := d.dynamoClient.Query(ctx, input)
	if err != nil {
		return nil, "", fmt.Errorf("query failed: %v", err)
	}

	if result.Items == nil || len(result.Items) == 0 && params.StartKey == "" {
		return nil, "", models.ErrSavingsNotFound
	}

	if result.Items == nil || len(result.Items) == 0 {
		return nil, "", models.ErrNoMoreItemsToBeRetrieved
	}

	savings := new([]savingEntity)

	err = attributevalue.UnmarshalListOfMaps(result.Items, savings)
	if err != nil {
		return nil, "", fmt.Errorf("unmarshal savings items failed: %v", err)
	}

	nextKey, err := dynamo.EncodePaginationKey(result.LastEvaluatedKey)
	if err != nil {
		return nil, "", err
	}

	return toSavingModels(*savings), nextKey, nil
}

func (d *DynamoRepository) buildQueryInput(username string, params *models.QueryParameters) (*dynamodb.QueryInput, error) {
	var err error

	input := &dynamodb.QueryInput{
		TableName: aws.String(d.tableName),
		Limit:     dynamo.GetPageSize(params.PageSize),
	}

	if params.SortType == string(models.SortOrderDescending) {
		input.ScanIndexForward = aws.Bool(false)
	}

	keyConditionEx := d.setQueryIndex(input, username, params)

	err = dynamo.SetExclusiveStartKey(params.StartKey, input)
	if err != nil {
		return nil, err
	}

	conditionBuilder := expression.NewBuilder().WithCondition(keyConditionEx)

	expr, err := conditionBuilder.Build()
	if err != nil {
		return nil, err
	}

	input.ExpressionAttributeNames = expr.Names()
	input.ExpressionAttributeValues = expr.Values()
	input.KeyConditionExpression = expr.Condition()
	input.FilterExpression = expr.Filter()

	return input, nil
}

func (d *DynamoRepository) setQueryIndex(input *dynamodb.QueryInput, username string, params *models.QueryParameters) expression.ConditionBuilder {
	keyConditionEx := expression.Name("username").Equal(expression.Value(username))
	periodUser := dynamo.BuildPeriodUser(username, params.Period)

	if params.Period != "" {
		input.IndexName = aws.String(d.periodSavingIndex)

		keyConditionEx = expression.Name("period_user").Equal(expression.Value(periodUser))
	}

	if params.SortBy == string(models.SortParamCreatedDate) {
		input.IndexName = aws.String(d.usernameCreatedDateIndex)
	}

	if params.Period != "" && params.SortBy == string(models.SortParamCreatedDate) {
		input.IndexName = aws.String(d.periodUserCreatedDateIndex)
	}

	if params.SortBy == string(models.SortParamAmount) {
		input.IndexName = aws.String(d.usernameAmountIndex)
	}

	if params.Period != "" && params.SortBy == string(models.SortParamAmount) {
		input.IndexName = aws.String(d.periodUserAmountIndex)
	}

	if params.SavingGoalID != "" && params.SortBy == string(models.SortParamCreatedDate) {
		keyConditionEx = expression.Name("saving_goal_id").Equal(expression.Value(params.SavingGoalID))
		input.IndexName = aws.String(d.savingGoalCreatedDateSavingIndex)
	}

	if params.SavingGoalID != "" {
		keyConditionEx = expression.Name("saving_goal_id").Equal(expression.Value(params.SavingGoalID))
		input.IndexName = aws.String(d.savingGoalSavingIndex)
	}

	return keyConditionEx
}

func (d *DynamoRepository) GetSavingsByPeriod(ctx context.Context, username string, params *models.QueryParameters) ([]*models.Saving, string, error) {
	input, err := d.buildQueryInput(username, params)
	if err != nil {
		return nil, "", fmt.Errorf("building query input: %v", err)
	}

	result, err := d.dynamoClient.Query(ctx, input)
	if err != nil {
		return nil, "", fmt.Errorf("query failed: %v", err)
	}

	if result.Items == nil || len(result.Items) == 0 && params.StartKey == "" {
		return nil, "", models.ErrSavingsNotFound
	}

	if result.Items == nil || len(result.Items) == 0 {
		return nil, "", models.ErrNoMoreItemsToBeRetrieved
	}

	savings := new([]savingEntity)

	err = attributevalue.UnmarshalListOfMaps(result.Items, savings)
	if err != nil {
		return nil, "", fmt.Errorf("unmarshal savings items failed: %v", err)
	}

	nextKey, err := dynamo.EncodePaginationKey(result.LastEvaluatedKey)
	if err != nil {
		return nil, "", err
	}

	return toSavingModels(*savings), nextKey, nil
}

func (d *DynamoRepository) GetSavingsBySavingGoal(ctx context.Context, params *models.QueryParameters) ([]*models.Saving, string, error) {
	input, err := d.buildQueryInput("", params)
	if err != nil {
		return nil, "", fmt.Errorf("building query input: %v", err)
	}

	result, err := d.dynamoClient.Query(ctx, input)
	if err != nil {
		return nil, "", fmt.Errorf("query failed: %v", err)
	}

	if result.Items == nil || len(result.Items) == 0 && params.StartKey == "" {
		return nil, "", models.ErrSavingsNotFound
	}

	if result.Items == nil || len(result.Items) == 0 {
		return nil, "", models.ErrNoMoreItemsToBeRetrieved
	}

	savings := new([]savingEntity)

	err = attributevalue.UnmarshalListOfMaps(result.Items, savings)
	if err != nil {
		return nil, "", fmt.Errorf("unmarshal savings items failed: %v", err)
	}

	nextKey, err := dynamo.EncodePaginationKey(result.LastEvaluatedKey)
	if err != nil {
		return nil, "", err
	}

	return toSavingModels(*savings), nextKey, nil
}

func (d *DynamoRepository) GetSavingsBySavingGoalAndPeriod(ctx context.Context, params *models.QueryParameters) ([]*models.Saving, string, error) {
	var decodedStartKey map[string]types.AttributeValue
	var err error
	var result *dynamodb.QueryOutput
	retrievedItems := 0
	resultSet := make([]savingEntity, 0)

	if params.StartKey != "" {
		decodedStartKey, err = dynamo.DecodePaginationKey(params.StartKey)
		if err != nil {
			return nil, "", fmt.Errorf("%v: %w", err, models.ErrInvalidStartKey)
		}
	}

	nameEx := expression.Name("saving_goal_id").Equal(expression.Value(params.SavingGoalID))
	filterCondition := expression.Name("period").Equal(expression.Value(params.Period))

	expr, err := expression.NewBuilder().WithCondition(nameEx).WithFilter(filterCondition).Build()
	if err != nil {
		return nil, "", err
	}

	input := &dynamodb.QueryInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		KeyConditionExpression:    expr.Condition(),
		TableName:                 aws.String(d.tableName),
		IndexName:                 aws.String(d.savingGoalSavingIndex),
		ExclusiveStartKey:         decodedStartKey,
		Limit:                     getPageSize(params.PageSize),
	}

	for {
		itemsInQuery := make([]savingEntity, 0)

		result, err = d.dynamoClient.Query(ctx, input)
		if err != nil {
			return nil, "", fmt.Errorf("query failed: %v", err)
		}

		input.ExclusiveStartKey = result.LastEvaluatedKey

		err = attributevalue.UnmarshalListOfMaps(result.Items, &itemsInQuery)
		if err != nil {
			return nil, "", fmt.Errorf("unmarshal savings items failed: %v", err)
		}

		retrievedItems += len(result.Items)

		// should implement custom pagination?
		if retrievedItems >= int(*input.Limit) {
			return getPaginatedSavings(resultSet, itemsInQuery, input)
		}

		resultSet = append(resultSet, itemsInQuery...)

		if result.LastEvaluatedKey == nil {
			break
		}
	}

	nextKey, err := dynamo.EncodePaginationKey(input.ExclusiveStartKey)
	if err != nil {
		return nil, "", err
	}

	if len(resultSet) == 0 && params.StartKey == "" {
		return nil, "", models.ErrSavingsNotFound
	}

	if len(resultSet) == 0 {
		return nil, "", models.ErrNoMoreItemsToBeRetrieved
	}

	return toSavingModels(resultSet), nextKey, nil
}

func getPaginatedSavings(resultSet, itemsInQuery []savingEntity, input *dynamodb.QueryInput) ([]*models.Saving, string, error) {
	var err error

	copyUpto := getCopyUpto(itemsInQuery, resultSet, input)
	resultSet = append(resultSet, itemsInQuery[0:copyUpto]...)

	input.ExclusiveStartKey, err = getAttributeValuePK(resultSet[len(resultSet)-1])
	if err != nil {
		return nil, "", fmt.Errorf("get attribute value pk failed: %v", err)
	}

	nextKey, err := dynamo.EncodePaginationKey(input.ExclusiveStartKey)
	if err != nil {
		return nil, "", err
	}

	if len(resultSet) == 0 {
		return nil, "", models.ErrExpensesNotFound
	}

	return toSavingModels(resultSet), nextKey, nil
}

// getCopyUpto returns the index up to which we can copy the items from the current query result to the list of items to
// return. This ensures that the total quantity of requested items, as indicated by the pageSize parameter, is satisfied.
func getCopyUpto(itemsInQuery []savingEntity, savingsEntities []savingEntity, input *dynamodb.QueryInput) int {
	limitAccumulatedDiff := int(math.Abs(float64(int(*input.Limit) - len(savingsEntities))))
	if len(itemsInQuery) < limitAccumulatedDiff {
		return len(itemsInQuery)
	}

	return limitAccumulatedDiff
}

func getAttributeValuePK(item savingEntity) (map[string]types.AttributeValue, error) {
	expenseKeys := struct {
		SavingID     string `json:"saving_id" dynamodbav:"saving_id"`
		Username     string `json:"username,omitempty" dynamodbav:"username"`
		SavingGoalID string `json:"saving_goal_id" dynamodbav:"saving_goal_id"`
	}{
		SavingID:     item.SavingID,
		Username:     item.Username,
		SavingGoalID: *item.SavingGoalID,
	}

	return attributevalue.MarshalMap(expenseKeys)
}

func (d *DynamoRepository) CreateSaving(ctx context.Context, saving *models.Saving) (*models.Saving, error) {
	saving.SavingID = dynamo.GenerateID(savingsPrefix)
	saving.CreatedDate = time.Now()
	savingEnt := toSavingEntity(saving)

	periodUser := dynamo.BuildPeriodUser(savingEnt.Username, *savingEnt.PeriodID)
	savingEnt.PeriodUser = periodUser

	item, err := attributevalue.MarshalMap(savingEnt)
	if err != nil {
		return nil, err
	}

	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(d.tableName),
	}

	_, err = d.dynamoClient.PutItem(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("put saving item failed: %v", err)
	}

	return toSavingModel(*savingEnt), nil
}

func (d *DynamoRepository) BatchCreateSavings(ctx context.Context, savings []*models.Saving) error {
	entities := make([]*savingEntity, 0, len(savings))

	for _, saving := range savings {
		saving.SavingID = dynamo.GenerateID(savingsPrefix)
		saving.CreatedDate = time.Now()

		entity := toSavingEntity(saving)
		if entity.PeriodID == nil {
			logger.Error("saving_with_nil_period", nil, models.Any("saving_model", saving),
				models.Any("saving_entity", entity))
			continue
		}

		entity.PeriodUser = dynamo.BuildPeriodUser(entity.Username, *entity.PeriodID)
		entities = append(entities, entity)
	}

	input := &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]types.WriteRequest{
			d.tableName: getBatchWriteRequests(entities),
		},
	}

	return dynamo.BatchWrite(ctx, d.dynamoClient, input)
}

func (d *DynamoRepository) BatchUpdateSavings(ctx context.Context, savings []*models.Saving) error {
	entities := make([]*savingEntity, 0, len(savings))

	for _, saving := range savings {
		saving.UpdatedDate = time.Now()

		entity := toSavingEntity(saving)
		if entity.PeriodID == nil {
			logger.Error("saving_with_nil_period", nil, models.Any("saving_model", saving),
				models.Any("saving_entity", entity))
			continue
		}

		entity.PeriodUser = dynamo.BuildPeriodUser(entity.Username, *entity.PeriodID)
		entities = append(entities, entity)
	}

	input := &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]types.WriteRequest{
			d.tableName: getBatchWriteRequests(entities),
		},
	}

	return dynamo.BatchWrite(ctx, d.dynamoClient, input)
}

func getBatchWriteRequests(entities []*savingEntity) []types.WriteRequest {
	writeRequests := make([]types.WriteRequest, 0, len(entities))

	for _, entity := range entities {
		item, err := attributevalue.MarshalMap(entity)
		if err != nil {
			logger.Warning("marshal_saving_failed", err, models.Any("saving_entity", entity))
			continue
		}

		writeRequests = append(writeRequests, types.WriteRequest{
			PutRequest: &types.PutRequest{
				Item: item,
			},
		})
	}

	return writeRequests
}

func (d *DynamoRepository) UpdateSaving(ctx context.Context, saving *models.Saving) error {
	savingEnt := toSavingEntity(saving)

	if savingEnt.PeriodID != nil {
		periodUser := dynamo.BuildPeriodUser(savingEnt.Username, *savingEnt.PeriodID)
		savingEnt.PeriodUser = periodUser
	}

	username, err := attributevalue.Marshal(savingEnt.Username)
	if err != nil {
		return fmt.Errorf("marshaling username key: %v", err)
	}

	savingID, err := attributevalue.Marshal(savingEnt.SavingID)
	if err != nil {
		return fmt.Errorf("marshaling saving id key: %v", err)
	}

	attributeValues, err := getAttributeValues(savingEnt)
	if err != nil {
		return fmt.Errorf("getting attribute values: %v", err)
	}

	updateExpression := getUpdateExpression(attributeValues)

	input := &dynamodb.UpdateItemInput{
		Key: map[string]types.AttributeValue{
			"username":  username,
			"saving_id": savingID,
		},
		TableName:                 aws.String(d.tableName),
		ConditionExpression:       aws.String("attribute_exists(saving_id)"),
		ExpressionAttributeValues: attributeValues,
		UpdateExpression:          updateExpression,
	}

	_, err = d.dynamoClient.UpdateItem(ctx, input)
	if err != nil && strings.Contains(err.Error(), "ConditionalCheckFailedException") {
		return fmt.Errorf("%v: %w", err, models.ErrUpdateSavingNotFound)
	}

	if err != nil {
		return fmt.Errorf("updating saving item: %v", err)
	}

	return nil
}

func getAttributeValues(saving *savingEntity) (map[string]types.AttributeValue, error) {
	m := make(map[string]types.AttributeValue)

	savingGoalID, err := attributevalue.Marshal(saving.SavingGoalID)
	if err != nil {
		return nil, err
	}

	amount, err := attributevalue.Marshal(saving.Amount)
	if err != nil {
		return nil, err
	}

	updatedDate, err := attributevalue.Marshal(time.Now())
	if err != nil {
		return nil, err
	}

	period, err := attributevalue.Marshal(saving.PeriodID)
	if err != nil {
		return nil, err
	}

	periodUser, err := attributevalue.Marshal(saving.PeriodUser)
	if err != nil {
		return nil, err
	}

	if saving.SavingGoalID != nil {
		m[":saving_goal_id"] = savingGoalID
	}

	if saving.Amount != nil {
		m[":amount"] = amount
	}

	if saving.PeriodID != nil {
		m[":period"] = period
		m[":period_user"] = periodUser
	}

	m[":updated_date"] = updatedDate

	return m, nil
}

func getUpdateExpression(attributeValues map[string]types.AttributeValue) *string {
	attributes := make([]string, 0)

	for key, _ := range attributeValues {
		attributeName := strings.ReplaceAll(key, ":", "")
		//The assumption here is that the attribute name is the same as the key without the colon
		//Example: "amount(attribute)" -> ":amount(key)"
		attributes = append(attributes, fmt.Sprintf("%s = %s", attributeName, key))
	}

	return aws.String("SET " + strings.Join(attributes, ", "))
}

func getPageSize(pageSize int) *int32 {
	if pageSize == 0 {
		return aws.Int32(defaultPageSize)
	}

	return aws.Int32(int32(pageSize))
}

func (d *DynamoRepository) DeleteSaving(ctx context.Context, savingID, username string) error {
	usernameAtr, err := attributevalue.Marshal(username)
	if err != nil {
		return fmt.Errorf("marshaling username key: %v", err)
	}

	savingIDAtr, err := attributevalue.Marshal(savingID)
	if err != nil {
		return fmt.Errorf("marshaling saving id key: %v", err)
	}

	input := &dynamodb.DeleteItemInput{
		Key: map[string]types.AttributeValue{
			"username":  usernameAtr,
			"saving_id": savingIDAtr,
		},
		TableName:           aws.String(d.tableName),
		ConditionExpression: aws.String("attribute_exists(saving_id)"),
	}

	_, err = d.dynamoClient.DeleteItem(ctx, input)
	if err != nil && strings.Contains(err.Error(), "ConditionalCheckFailedException") {
		return fmt.Errorf("%v: %w", err, models.ErrSavingsNotFound)
	}

	if err != nil {
		return fmt.Errorf("deleting item: %v", err)
	}

	return nil
}
