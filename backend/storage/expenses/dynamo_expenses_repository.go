package expenses

import (
	"context"
	"fmt"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/env"
	"github.com/JoelD7/money/backend/storage/shared"
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
	defaultPageSize   = 10
	nameAttributeName = "#n"
)

var (
	tableName                = env.GetString("EXPENSES_TABLE_NAME", "expenses")
	periodUserExpenseIDIndex = "period_user-expense_id-index"
)

type keys struct {
	ExpenseID string `json:"expense_id" dynamodbav:"expense_id"`
	Username  string `json:"username" dynamodbav:"username"`
}

type keysPeriodUserIndex struct {
	ExpenseID  string `json:"expense_id" dynamodbav:"expense_id"`
	Username   string `json:"username,omitempty" dynamodbav:"username"`
	PeriodUser string `json:"period_user,omitempty" dynamodbav:"period_user"`
}

type DynamoRepository struct {
	dynamoClient *dynamodb.Client
}

func NewDynamoRepository(dynamoClient *dynamodb.Client) *DynamoRepository {
	return &DynamoRepository{dynamoClient: dynamoClient}
}

func (d *DynamoRepository) CreateExpense(ctx context.Context, expense *models.Expense) (*models.Expense, error) {
	entity := toExpenseEntity(expense)

	entity.PeriodUser = shared.BuildPeriodUser(entity.Username, *entity.Period)

	item, err := attributevalue.MarshalMap(entity)
	if err != nil {
		return nil, fmt.Errorf("marshal expense failed: %v", err)
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      item,
	}

	_, err = d.dynamoClient.PutItem(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("put expense failed: %v", err)
	}

	return toExpenseModel(*entity), nil
}

func (d *DynamoRepository) UpdateExpense(ctx context.Context, expense *models.Expense) error {
	entity := toExpenseEntity(expense)

	if entity.Period != nil {
		entity.PeriodUser = shared.BuildPeriodUser(entity.Username, *entity.Period)
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
		TableName:                 aws.String(tableName),
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

	if expense.Period != nil {
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

func (d *DynamoRepository) GetExpense(ctx context.Context, username, expenseID string) (*models.Expense, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
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

func (d *DynamoRepository) GetExpenses(ctx context.Context, username, startKey string, pageSize int) ([]*models.Expense, string, error) {
	input, err := buildQueryInput(username, "", startKey, nil, pageSize)
	if err != nil {
		return nil, "", err
	}

	return d.performQuery(ctx, input, startKey)
}

func (d *DynamoRepository) GetExpensesByPeriodAndCategories(ctx context.Context, username, periodID, startKey string, categories []string, pageSize int) ([]*models.Expense, string, error) {
	input, err := buildQueryInput(username, periodID, startKey, categories, pageSize)
	if err != nil {
		return nil, "", err
	}

	return d.performQuery(ctx, input, startKey)
}

func (d *DynamoRepository) GetExpensesByPeriod(ctx context.Context, username, periodID, startKey string, pageSize int) ([]*models.Expense, string, error) {
	input, err := buildQueryInput(username, periodID, startKey, nil, pageSize)
	if err != nil {
		return nil, "", err
	}

	return d.performQuery(ctx, input, startKey)
}

func (d *DynamoRepository) GetExpensesByCategory(ctx context.Context, username, startKey string, categories []string, pageSize int) ([]*models.Expense, string, error) {
	input, err := buildQueryInput(username, "", startKey, categories, pageSize)
	if err != nil {
		return nil, "", err
	}

	return d.performQuery(ctx, input, startKey)
}

func (d *DynamoRepository) DeleteExpense(ctx context.Context, expenseID, username string) error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
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

func buildQueryInput(username, periodID, startKey string, categories []string, pageSize int) (*dynamodb.QueryInput, error) {
	var err error

	input := &dynamodb.QueryInput{
		TableName: aws.String(tableName),
		Limit:     getPageSize(pageSize),
	}

	keyConditionEx := expression.Name("username").Equal(expression.Value(username))

	// Query the period_user-expense_id-index
	if periodID != "" {
		input.IndexName = aws.String(periodUserExpenseIDIndex)

		periodUser := shared.BuildPeriodUser(username, periodID)
		keyConditionEx = expression.Name("period_user").Equal(expression.Value(periodUser))
	}

	err = setExclusiveStartKey(startKey, input)
	if err != nil {
		return nil, err
	}

	conditionBuilder := expression.NewBuilder().WithCondition(keyConditionEx)

	if categories != nil || len(categories) > 0 {
		filterCondition := buildCategoriesConditionFilter(categories)
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

func setExclusiveStartKey(startKey string, input *dynamodb.QueryInput) error {
	if startKey == "" {
		return nil
	}

	decodedStartKey, err := shared.DecodePaginationKey(startKey, getPaginationKeyType(input))
	if err != nil {
		return fmt.Errorf("%v: %w", err, models.ErrInvalidStartKey)
	}

	input.ExclusiveStartKey = decodedStartKey

	return nil
}

// getPaginationKeyType returns the type of the key to be used in the pagination according to the index being queried.
// If there isn't an index being queried, the key type corresponding to the main table is used.
func getPaginationKeyType(input *dynamodb.QueryInput) interface{} {
	indexName := ""
	if input.IndexName != nil {
		indexName = *input.IndexName
	}

	switch indexName {
	case periodUserExpenseIDIndex:
		return &keysPeriodUserIndex{}
	default:
		return &keys{}
	}
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

	nextKey, err := shared.EncodePaginationKey(result.LastEvaluatedKey, getPaginationKeyType(input))
	if err != nil {
		return nil, "", err
	}

	return toExpenseModels(expensesEntities), nextKey, nil
}

func (d *DynamoRepository) performQueryWithFilter(ctx context.Context, input *dynamodb.QueryInput, startKey string) ([]*models.Expense, string, error) {
	retrievedItems := 0
	resultSet := make([]expenseEntity, 0)
	var result *dynamodb.QueryOutput
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

		// should we implement custom pagination?
		if retrievedItems >= int(*input.Limit) {
			return getPaginatedExpenses(resultSet, itemsInQuery, input)
		}

		resultSet = append(resultSet, itemsInQuery...)

		if result.LastEvaluatedKey == nil {
			break
		}
	}

	nextKey, err := shared.EncodePaginationKey(input.ExclusiveStartKey, getPaginationKeyType(input))
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

func getPaginatedExpenses(resultSet, itemsInQuery []expenseEntity, input *dynamodb.QueryInput) ([]*models.Expense, string, error) {
	var err error

	copyUpto := getCopyUpto(itemsInQuery, resultSet, input)
	resultSet = append(resultSet, itemsInQuery[0:copyUpto]...)

	input.ExclusiveStartKey, err = getAttributeValuePK(resultSet[len(resultSet)-1], input)
	if err != nil {
		return nil, "", fmt.Errorf("get attribute value pk failed: %v", err)
	}

	nextKey, err := shared.EncodePaginationKey(input.ExclusiveStartKey, getPaginationKeyType(input))
	if err != nil {
		return nil, "", err
	}

	if len(resultSet) == 0 {
		return nil, "", models.ErrExpensesNotFound
	}

	return toExpenseModels(resultSet), nextKey, nil
}

// getCopyUpto returns the index up to which we can copy the items from the current query result to the list of items to
// return. This ensures that the total quantity of requested items, as indicated by the pageSize parameter, is satisfied.
func getCopyUpto(itemsInQuery []expenseEntity, expensesEntities []expenseEntity, input *dynamodb.QueryInput) int {
	limitAccumulatedDiff := int(math.Abs(float64(int(*input.Limit) - len(expensesEntities))))
	if len(itemsInQuery) < limitAccumulatedDiff {
		return len(itemsInQuery)
	}

	return limitAccumulatedDiff
}

func getAttributeValuePK(item expenseEntity, input *dynamodb.QueryInput) (map[string]types.AttributeValue, error) {
	if input.IndexName != nil && *input.IndexName == periodUserExpenseIDIndex {
		expenseKeys := struct {
			ExpenseID  string `json:"expense_id" dynamodbav:"expense_id"`
			Username   string `json:"username,omitempty" dynamodbav:"username"`
			PeriodUser string `json:"period_user,omitempty" dynamodbav:"period_user"`
		}{
			ExpenseID:  item.ExpenseID,
			Username:   item.Username,
			PeriodUser: *item.PeriodUser,
		}

		return attributevalue.MarshalMap(expenseKeys)
	}

	expenseKeys := struct {
		ExpenseID string `json:"expense_id" dynamodbav:"expense_id"`
		Username  string `json:"username,omitempty" dynamodbav:"username"`
	}{
		ExpenseID: item.ExpenseID,
		Username:  item.Username,
	}

	return attributevalue.MarshalMap(expenseKeys)
}

func getPageSize(pageSize int) *int32 {
	if pageSize == 0 {
		return aws.Int32(defaultPageSize)
	}

	return aws.Int32(int32(pageSize))
}
