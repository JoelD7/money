package expenses

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
	defaultPageSize          = 10
	nameAttributeName        = "#n"
	conditionalFailedKeyword = "ConditionalCheckFailed"
)

var (
	//All indices must include "expense_id" and "username" even when they aren't keys of the index, because they are
	//the main table's primary key and we get a validation error from Dynamo otherwise.
	keysByIndex = map[string][]string{
		"period_user-created_date-index": {"expense_id", "username", "period_user", "created_date"},
		"period_user-expense_id-index":   {"expense_id", "username", "period_user", "expense_id"},
		"username-created_date-index":    {"expense_id", "username", "created_date"},
	}
)

type DynamoRepository struct {
	dynamoClient                 *dynamodb.Client
	tableName                    string
	expensesRecurringTableName   string
	periodUserIndex              string
	periodUserCreatedDateIndex   string
	usernameCreatedDateIndex     string
	periodUserNameExpenseIDIndex string
	periodUserAmountIndex        string
}

func NewDynamoRepository(dynamoClient *dynamodb.Client, envConfig *models.EnvironmentConfiguration) (*DynamoRepository, error) {
	d := &DynamoRepository{dynamoClient: dynamoClient}

	err := validateParams(envConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize expenses dynamo repository: %v", err)
	}

	d.tableName = envConfig.ExpensesTable
	d.expensesRecurringTableName = envConfig.ExpensesRecurringTable
	d.periodUserIndex = envConfig.PeriodUserExpenseIndex
	d.periodUserCreatedDateIndex = envConfig.PeriodUserCreatedDateIndex
	d.usernameCreatedDateIndex = envConfig.UsernameCreatedDateIndex
	d.periodUserNameExpenseIDIndex = envConfig.PeriodUserNameExpenseIDIndex
	d.periodUserAmountIndex = envConfig.PeriodUserAmountIndex

	return d, nil
}

func validateParams(envConfig *models.EnvironmentConfiguration) error {
	if envConfig.ExpensesTable == "" {
		return fmt.Errorf("table name is required")
	}

	if envConfig.ExpensesRecurringTable == "" {
		return fmt.Errorf("expenses recurring table name is required")
	}

	if envConfig.PeriodUserExpenseIndex == "" {
		return fmt.Errorf("period user index is required")
	}

	if envConfig.PeriodUserCreatedDateIndex == "" {
		return fmt.Errorf("period user created date index is required")
	}

	if envConfig.UsernameCreatedDateIndex == "" {
		return fmt.Errorf("username created date index is required")
	}

	if envConfig.PeriodUserNameExpenseIDIndex == "" {
		return fmt.Errorf("period user name expense id index is required")
	}

	if envConfig.PeriodUserAmountIndex == "" {
		return fmt.Errorf("period user amount index is required")
	}

	return nil
}

func (d *DynamoRepository) CreateExpense(ctx context.Context, expense *models.Expense) (*models.Expense, error) {
	entity := toExpenseEntity(expense)

	input, err := d.buildTransactWriteItemsInput(entity, expense)
	if err != nil {
		return nil, err
	}

	_, err = d.dynamoClient.TransactWriteItems(ctx, input)
	if err != nil && strings.Contains(err.Error(), conditionalFailedKeyword) {
		return nil, fmt.Errorf("%v: %w", err, models.ErrRecurringExpenseNameTaken)
	}

	if err != nil {
		return nil, fmt.Errorf("put expense failed: %v", err)
	}

	return toExpenseModel(*entity), nil
}

func (d *DynamoRepository) buildTransactWriteItemsInput(expenseEnt *expenseEntity, expense *models.Expense) (*dynamodb.TransactWriteItemsInput, error) {
	expenseEnt.PeriodUser = dynamo.BuildPeriodUser(expenseEnt.Username, expenseEnt.Period)

	item, err := attributevalue.MarshalMap(expenseEnt)
	if err != nil {
		return nil, fmt.Errorf("marshal expense failed: %v", err)
	}

	transactItems := []types.TransactWriteItem{
		{
			Put: &types.Put{
				Item:      item,
				TableName: aws.String(d.tableName),
			},
		},
	}

	if !expense.IsRecurring {
		return &dynamodb.TransactWriteItemsInput{
			TransactItems: transactItems,
		}, nil
	}

	expenseRecurringEnt := toExpenseRecurringEntity(expense)

	itemRecurring, err := attributevalue.MarshalMap(expenseRecurringEnt)
	if err != nil {
		return nil, fmt.Errorf("marshal expense recurring failed: %v", err)
	}

	condExpr := expression.Name("id").AttributeNotExists().And(expression.Name("username").AttributeNotExists())

	expr, err := expression.NewBuilder().WithCondition(condExpr).Build()
	if err != nil {
		return nil, fmt.Errorf("build expression failed: %v", err)
	}

	transactItems = append(transactItems, types.TransactWriteItem{
		Put: &types.Put{
			Item:                     itemRecurring,
			TableName:                aws.String(d.expensesRecurringTableName),
			ConditionExpression:      expr.Condition(),
			ExpressionAttributeNames: expr.Names(),
		},
	})

	return &dynamodb.TransactWriteItemsInput{
		TransactItems: transactItems,
	}, nil
}

func (d *DynamoRepository) BatchCreateExpenses(ctx context.Context, expenses []*models.Expense) error {
	entities := make([]*expenseEntity, 0, len(expenses))

	for _, expense := range expenses {
		entity := toExpenseEntity(expense)
		entity.PeriodUser = dynamo.BuildPeriodUser(entity.Username, entity.Period)
		entities = append(entities, entity)
	}

	input := &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]types.WriteRequest{
			d.tableName: getBatchWriteRequests(entities),
		},
	}

	return dynamo.BatchWrite(ctx, d.dynamoClient, input)
}

func getBatchWriteRequests(entities []*expenseEntity) []types.WriteRequest {
	writeRequests := make([]types.WriteRequest, 0, len(entities))

	for _, entity := range entities {
		item, err := attributevalue.MarshalMap(entity)
		if err != nil {
			logger.Warning("marshal_expense_failed", err, models.Any("expense_entity", entity))
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

func (d *DynamoRepository) UpdateExpense(ctx context.Context, expense *models.Expense) error {
	entity := toExpenseEntity(expense)

	if entity.Period != "" {
		entity.PeriodUser = dynamo.BuildPeriodUser(entity.Username, entity.Period)
	}

	username, err := attributevalue.Marshal(entity.Username)
	if err != nil {
		return fmt.Errorf("marshaling username key: %v", err)
	}

	expenseID, err := attributevalue.Marshal(entity.ExpenseID)
	if err != nil {
		return fmt.Errorf("marshaling expense id key: %v", err)
	}

	attributeValues, err := getAttributeValues(entity)
	if err != nil {
		return fmt.Errorf("get attribute values failed: %v", err)
	}

	updateExpression := getUpdateExpression(attributeValues)

	input := &dynamodb.UpdateItemInput{
		Key: map[string]types.AttributeValue{
			"username":   username,
			"expense_id": expenseID,
		},
		TableName:                 aws.String(d.tableName),
		ConditionExpression:       aws.String("attribute_exists(expense_id)"),
		ExpressionAttributeValues: attributeValues,
		UpdateExpression:          updateExpression,
	}

	if expense.Name != nil {
		input.ExpressionAttributeNames = map[string]string{nameAttributeName: "name"}
	}

	_, err = d.dynamoClient.UpdateItem(ctx, input)
	if err != nil && strings.Contains(err.Error(), "ConditionalCheckFailedException") {
		return fmt.Errorf("%v: %w", err, models.ErrExpensesNotFound)
	}

	if err != nil {
		return fmt.Errorf("put expense failed: %v", err)
	}

	return nil
}

func getAttributeValues(expense *expenseEntity) (map[string]types.AttributeValue, error) {
	attrValues := make(map[string]types.AttributeValue)

	categoryID, err := attributevalue.Marshal(expense.CategoryID)
	if err != nil {
		return nil, err
	}

	amount, err := attributevalue.Marshal(expense.Amount)
	if err != nil {
		return nil, err
	}

	name, err := attributevalue.Marshal(expense.Name)
	if err != nil {
		return nil, err
	}

	notes, err := attributevalue.Marshal(expense.Notes)
	if err != nil {
		return nil, err
	}

	period, err := attributevalue.Marshal(expense.Period)
	if err != nil {
		return nil, err
	}

	periodUser, err := attributevalue.Marshal(expense.PeriodUser)
	if err != nil {
		return nil, err
	}

	updatedDate, err := attributevalue.Marshal(time.Now())
	if err != nil {
		return nil, err
	}

	if expense.CategoryID != nil {
		attrValues[":category_id"] = categoryID
	}

	if expense.Amount != 0 {
		attrValues[":amount"] = amount
	}

	if expense.Name != "" {
		attrValues[":name"] = name
	}

	if expense.Notes != "" {
		attrValues[":notes"] = notes
	}

	if expense.Period != "" {
		attrValues[":period"] = period
		attrValues[":period_user"] = periodUser
	}

	attrValues[":update_date"] = updatedDate

	return attrValues, nil
}

func getUpdateExpression(attributeValues map[string]types.AttributeValue) *string {
	attributes := make([]string, 0)

	for key, _ := range attributeValues {
		attributeName := strings.ReplaceAll(key, ":", "")
		if key == ":name" {
			attributes = append(attributes, fmt.Sprintf("%s = :name", nameAttributeName))
			continue
		}

		//The assumption here is that the attribute name is the same as the key without the colon
		//Example: "amount(attribute)" -> ":amount(key)"
		attributes = append(attributes, fmt.Sprintf("%s = %s", attributeName, key))
	}

	return aws.String("SET " + strings.Join(attributes, ", "))
}

func (d *DynamoRepository) BatchUpdateExpenses(ctx context.Context, expenses []*models.Expense) error {
	entities := make([]*expenseEntity, 0, len(expenses))

	for _, expense := range expenses {
		entity := toExpenseEntity(expense)
		entity.PeriodUser = dynamo.BuildPeriodUser(entity.Username, entity.Period)
		entities = append(entities, entity)
	}

	input := &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]types.WriteRequest{
			d.tableName: getBatchWriteRequests(entities),
		},
	}

	return dynamo.BatchWrite(ctx, d.dynamoClient, input)
}

func (d *DynamoRepository) GetExpense(ctx context.Context, username, expenseID string) (*models.Expense, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(d.tableName),
		Key: map[string]types.AttributeValue{
			"username":   &types.AttributeValueMemberS{Value: username},
			"expense_id": &types.AttributeValueMemberS{Value: expenseID},
		},
	}

	result, err := d.dynamoClient.GetItem(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("get expense failed: %v", err)
	}

	if result.Item == nil {
		return nil, models.ErrExpenseNotFound
	}

	entity := new(expenseEntity)

	err = attributevalue.UnmarshalMap(result.Item, entity)
	if err != nil {
		return nil, fmt.Errorf("unmarshal expense item failed: %v", err)
	}

	return toExpenseModel(*entity), nil
}

func (d *DynamoRepository) GetExpenses(ctx context.Context, username string, params *models.QueryParameters) ([]*models.Expense, string, error) {
	input, err := d.buildQueryInput(username, params)
	if err != nil {
		return nil, "", err
	}

	return d.performQuery(ctx, input, params.StartKey)
}

func (d *DynamoRepository) GetExpensesByPeriodAndCategories(ctx context.Context, username string, params *models.QueryParameters) ([]*models.Expense, string, error) {
	input, err := d.buildQueryInput(username, params)
	if err != nil {
		return nil, "", err
	}

	return d.performQuery(ctx, input, params.StartKey)
}

func (d *DynamoRepository) GetExpensesByPeriod(ctx context.Context, username string, params *models.QueryParameters) ([]*models.Expense, string, error) {
	input, err := d.buildQueryInput(username, params)
	if err != nil {
		return nil, "", err
	}

	return d.performQuery(ctx, input, params.StartKey)
}

func (d *DynamoRepository) GetExpensesByCategory(ctx context.Context, username string, params *models.QueryParameters) ([]*models.Expense, string, error) {
	input, err := d.buildQueryInput(username, params)
	if err != nil {
		return nil, "", err
	}

	return d.performQuery(ctx, input, params.StartKey)
}

func (d *DynamoRepository) GetAllExpensesBetweenDates(ctx context.Context, username, startDate, endDate string) ([]*models.Expense, error) {
	userFilter := expression.Name("username").Equal(expression.Value(username))
	dateFilter := expression.Name("created_date").Between(expression.Value(startDate), expression.Value(endDate))
	filter := expression.And(userFilter, dateFilter)

	expr, err := expression.NewBuilder().WithFilter(filter).Build()
	if err != nil {
		return nil, err
	}

	input := &dynamodb.ScanInput{
		TableName:                 aws.String(d.tableName),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ExclusiveStartKey:         nil,
	}

	var result *dynamodb.ScanOutput
	entities := make([]expenseEntity, 0)
	var itemsInQuery []expenseEntity

	for {
		itemsInQuery = make([]expenseEntity, 0)
		result, err = d.dynamoClient.Scan(ctx, input)
		if err != nil {
			return nil, err
		}

		if (result.Items == nil || len(result.Items) == 0) && result.LastEvaluatedKey == nil {
			break
		}

		err = attributevalue.UnmarshalListOfMaps(result.Items, &itemsInQuery)
		if err != nil {
			return nil, fmt.Errorf("unmarshal expenses items failed: %v", err)
		}

		entities = append(entities, itemsInQuery...)
		input.ExclusiveStartKey = result.LastEvaluatedKey

		if result.LastEvaluatedKey == nil {
			break
		}
	}

	if len(entities) == 0 {
		return nil, models.ErrExpensesNotFound
	}

	return toExpenseModels(entities), nil
}

func (d *DynamoRepository) GetAllExpensesByPeriod(ctx context.Context, username string, params *models.QueryParameters) ([]*models.Expense, error) {
	periodUser := dynamo.BuildPeriodUser(username, params.Period)
	periodUserCond := expression.Key("period_user").Equal(expression.Value(periodUser))

	expr, err := expression.NewBuilder().WithKeyCondition(periodUserCond).Build()
	if err != nil {
		return nil, err
	}

	input := &dynamodb.QueryInput{
		TableName:                 aws.String(d.tableName),
		IndexName:                 aws.String(d.periodUserIndex),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
	}

	var result *dynamodb.QueryOutput
	entities := make([]expenseEntity, 0)
	var itemsInQuery []expenseEntity

	for {
		itemsInQuery = make([]expenseEntity, 0)
		result, err = d.dynamoClient.Query(ctx, input)
		if err != nil {
			return nil, err
		}

		if (result.Items == nil || len(result.Items) == 0) && result.LastEvaluatedKey == nil {
			break
		}

		err = attributevalue.UnmarshalListOfMaps(result.Items, &itemsInQuery)
		if err != nil {
			return nil, fmt.Errorf("unmarshal expenses items failed: %v", err)
		}

		entities = append(entities, itemsInQuery...)
		input.ExclusiveStartKey = result.LastEvaluatedKey

		if result.LastEvaluatedKey == nil {
			break
		}
	}

	if len(entities) == 0 {
		return nil, models.ErrExpensesNotFound
	}

	return toExpenseModels(entities), nil
}

func (d *DynamoRepository) DeleteExpense(ctx context.Context, expenseID, username string) error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(d.tableName),
		Key: map[string]types.AttributeValue{
			"username":   &types.AttributeValueMemberS{Value: username},
			"expense_id": &types.AttributeValueMemberS{Value: expenseID},
		},
	}

	_, err := d.dynamoClient.DeleteItem(ctx, input)
	if err != nil {
		return fmt.Errorf("delete expense failed: %v", err)
	}

	return nil
}

func (d *DynamoRepository) BatchDeleteExpenses(ctx context.Context, expenses []*models.Expense) error {
	writeRequests := make([]types.WriteRequest, 0, len(expenses))

	var usernameAttrValue types.AttributeValue
	var expenseIDAttrValue types.AttributeValue
	var err error

	for _, expense := range expenses {
		usernameAttrValue, err = attributevalue.Marshal(expense.Username)
		if err != nil {
			return fmt.Errorf("marshal username key failed: %v", err)
		}

		expenseIDAttrValue, err = attributevalue.Marshal(expense.ExpenseID)
		if err != nil {
			return fmt.Errorf("marshal id key failed: %v", err)
		}

		writeRequests = append(writeRequests, types.WriteRequest{
			DeleteRequest: &types.DeleteRequest{
				Key: map[string]types.AttributeValue{
					"username":   usernameAttrValue,
					"expense_id": expenseIDAttrValue,
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

func (d *DynamoRepository) buildQueryInput(username string, params *models.QueryParameters) (*dynamodb.QueryInput, error) {
	var err error

	input := &dynamodb.QueryInput{
		TableName: aws.String(d.tableName),
		Limit:     getPageSize(params.PageSize),
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

	if params.Categories != nil || len(params.Categories) > 0 {
		filterCondition := buildCategoriesConditionFilter(params.Categories)
		conditionBuilder = conditionBuilder.WithFilter(filterCondition)
	}

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

// setQueryIndex sets the index to be used in the query based on the sorting and filter parameters. Returns a key
// condition expression formed with the index's primary key.
func (d *DynamoRepository) setQueryIndex(input *dynamodb.QueryInput, username string, params *models.QueryParameters) expression.ConditionBuilder {
	keyConditionEx := expression.Name("username").Equal(expression.Value(username))
	periodUser := dynamo.BuildPeriodUser(username, params.Period)

	if params.Period != "" {
		input.IndexName = aws.String(d.periodUserIndex)

		keyConditionEx = expression.Name("period_user").Equal(expression.Value(periodUser))
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
		input.IndexName = aws.String(d.periodUserNameExpenseIDIndex)
	}

	return keyConditionEx
}

func (d *DynamoRepository) setExclusiveStartKey(startKey string, input *dynamodb.QueryInput) error {
	if startKey == "" {
		return nil
	}

	decodedStartKey, err := dynamo.DecodePaginationKey(startKey)
	if err != nil {
		return fmt.Errorf("%v: %w", err, models.ErrInvalidStartKey)
	}

	input.ExclusiveStartKey = decodedStartKey

	return nil
}

func buildCategoriesConditionFilter(categories []string) expression.ConditionBuilder {
	if categories[0] == "" {
		return expression.Name("category_id").AttributeNotExists()
	}

	conditions := make([]expression.ConditionBuilder, 0, len(categories))

	for _, categoryID := range categories {
		conditions = append(conditions, expression.Name("category_id").Equal(expression.Value(categoryID)))
	}

	if len(categories) == 1 {
		return expression.Name("category_id").Equal(expression.Value(categories[0]))
	}

	if len(categories) == 2 {
		return expression.Or(conditions[0], conditions[1])
	}

	return expression.Or(conditions[0], conditions[1], conditions[2:]...)
}

func (d *DynamoRepository) performQuery(ctx context.Context, input *dynamodb.QueryInput, startKey string) ([]*models.Expense, string, error) {
	// If the query has a filter expression it may not include all the items one intends to fetch.
	// See more details here: https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/Query.FilterExpression.html
	if input.FilterExpression != nil {
		return d.performQueryWithFilter(ctx, input, startKey)
	}

	result, err := d.dynamoClient.Query(ctx, input)
	if err != nil {
		return nil, "", fmt.Errorf("query failed: %v", err)
	}

	if result.Items == nil || len(result.Items) == 0 && input.ExclusiveStartKey == nil {
		return nil, "", models.ErrExpensesNotFound
	}

	if result.Items == nil || len(result.Items) == 0 {
		return nil, "", models.ErrNoMoreItemsToBeRetrieved
	}

	expensesEntities := make([]expenseEntity, 0)

	err = attributevalue.UnmarshalListOfMaps(result.Items, &expensesEntities)
	if err != nil {
		return nil, "", fmt.Errorf("unmarshal expenses items failed: %v", err)
	}

	nextKey, err := dynamo.EncodePaginationKey(result.LastEvaluatedKey)
	if err != nil {
		return nil, "", err
	}

	return toExpenseModels(expensesEntities), nextKey, nil
}

func (d *DynamoRepository) performQueryWithFilter(ctx context.Context, input *dynamodb.QueryInput, startKey string) ([]*models.Expense, string, error) {
	retrievedItems := 0
	resultSet := make([]expenseEntity, 0)
	var result *dynamodb.QueryOutput
	var acumItems []map[string]types.AttributeValue
	var err error

	for {
		itemsInQuery := make([]expenseEntity, 0)

		result, err = d.dynamoClient.Query(ctx, input)
		if err != nil {
			return nil, "", fmt.Errorf("query failed: %v", err)
		}

		input.ExclusiveStartKey = result.LastEvaluatedKey

		err = attributevalue.UnmarshalListOfMaps(result.Items, &itemsInQuery)
		if err != nil {
			return nil, "", fmt.Errorf("unmarshal expenses items failed: %v", err)
		}

		retrievedItems += len(result.Items)

		// This asks: should we implement custom pagination?
		if retrievedItems >= int(*input.Limit) {
			return d.getPaginatedExpenses(acumItems, result.Items, input)
		}

		acumItems = append(acumItems, result.Items...)

		if result.LastEvaluatedKey == nil {
			break
		}
	}

	err = attributevalue.UnmarshalListOfMaps(acumItems, &resultSet)
	if err != nil {
		return nil, "", fmt.Errorf("unmarshal expenses items failed: %v", err)
	}

	nextKey, err := dynamo.EncodePaginationKey(input.ExclusiveStartKey)
	if err != nil {
		return nil, "", err
	}

	if len(resultSet) == 0 && startKey == "" {
		return nil, "", models.ErrExpensesNotFound
	}

	if len(resultSet) == 0 {
		return nil, "", models.ErrNoMoreItemsToBeRetrieved
	}

	return toExpenseModels(resultSet), nextKey, nil
}

func (d *DynamoRepository) getPaginatedExpenses(acumItems, itemsInQuery []map[string]types.AttributeValue, input *dynamodb.QueryInput) ([]*models.Expense, string, error) {
	var err error

	copyUpto := getCopyUpto(itemsInQuery, acumItems, input)
	acumItems = append(acumItems, itemsInQuery[0:copyUpto]...)

	input.ExclusiveStartKey, err = d.buildCustomExclusiveStartKey(acumItems[len(acumItems)-1], input)
	if err != nil {
		return nil, "", fmt.Errorf("getting expenses failed: %v", err)
	}

	nextKey, err := dynamo.EncodePaginationKey(input.ExclusiveStartKey)
	if err != nil {
		return nil, "", err
	}

	if len(acumItems) == 0 {
		return nil, "", models.ErrExpensesNotFound
	}

	var expenseEntites []expenseEntity

	err = attributevalue.UnmarshalListOfMaps(acumItems, &expenseEntites)
	if err != nil {
		return nil, "", err
	}

	return toExpenseModels(expenseEntites), nextKey, nil
}

func (d *DynamoRepository) buildCustomExclusiveStartKey(lastEvaluatedKey map[string]types.AttributeValue, input *dynamodb.QueryInput) (map[string]types.AttributeValue, error) {
	if lastEvaluatedKey == nil {
		return nil, nil
	}

	exclusiveStartKey := make(map[string]types.AttributeValue)

	keys, ok := keysByIndex[*input.IndexName]
	if !ok {
		return nil, models.ErrIndexKeysNotFound
	}

	for _, key := range keys {
		exclusiveStartKey[key] = lastEvaluatedKey[key]
	}

	return exclusiveStartKey, nil
}

// getCopyUpto returns the index up to which we can copy the items from the current query result to the list of items to
// return. This ensures that the total quantity of requested items, as indicated by the pageSize parameter, is satisfied.
func getCopyUpto(itemsInQuery, expensesEntities []map[string]types.AttributeValue, input *dynamodb.QueryInput) int {
	limitAccumulatedDiff := int(math.Abs(float64(int(*input.Limit) - len(expensesEntities))))
	if len(itemsInQuery) < limitAccumulatedDiff {
		return len(itemsInQuery)
	}

	return limitAccumulatedDiff
}

func getPageSize(pageSize int) *int32 {
	if pageSize == 0 {
		return aws.Int32(defaultPageSize)
	}

	return aws.Int32(int32(pageSize))
}
