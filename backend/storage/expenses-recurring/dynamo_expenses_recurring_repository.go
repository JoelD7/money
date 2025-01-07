package expenses_recurring

import (
	"context"
	"fmt"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/env"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/storage/dynamo"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
)

type DynamoRepository struct {
	dynamoClient *dynamodb.Client
	tableName    string
}

func NewExpenseRecurringDynamoRepository(dynamoClient *dynamodb.Client, tableName string) (*DynamoRepository, error) {
	d := &DynamoRepository{dynamoClient: dynamoClient}
	tableNameEnv := env.GetString("EXPENSES_RECURRING_TABLE_NAME", "")

	if tableNameEnv == "" && tableName == "" {
		return nil, fmt.Errorf("initialize expenses recurring dynamo repository failed: table name is required")
	}

	d.tableName = tableName
	if d.tableName == "" {
		d.tableName = tableNameEnv
	}

	return d, nil
}

func (d *DynamoRepository) CreateExpenseRecurring(ctx context.Context, expenseRecurring *models.ExpenseRecurring) (*models.ExpenseRecurring, error) {
	entity := toExpenseRecurringEntity(expenseRecurring)

	item, err := attributevalue.MarshalMap(entity)
	if err != nil {
		return nil, fmt.Errorf("marshal expense recurring entity failed: %v", err)
	}

	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(d.tableName),
	}

	_, err = d.dynamoClient.PutItem(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("put expense recurring item failed: %v", err)
	}

	return expenseRecurring, nil
}

func (d *DynamoRepository) BatchCreateExpenseRecurring(ctx context.Context, log logger.LogAPI, expenseRecurring []*models.ExpenseRecurring) error {
	entities := make([]*ExpenseRecurringEntity, 0, len(expenseRecurring))

	for _, expense := range expenseRecurring {
		entity := toExpenseRecurringEntity(expense)
		entities = append(entities, entity)
	}

	input := &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]types.WriteRequest{
			d.tableName: getBatchWriteRequests(entities, log),
		},
	}

	return dynamo.BatchWrite(ctx, d.dynamoClient, input)
}

func getBatchWriteRequests(entities []*ExpenseRecurringEntity, log logger.LogAPI) []types.WriteRequest {
	writeRequests := make([]types.WriteRequest, 0, len(entities))

	for _, entity := range entities {
		item, err := attributevalue.MarshalMap(entity)
		if err != nil {
			log.Warning("marshal_recurring_expense_failed", err, models.Any("expense_recurring_entity", entity))
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

func (d *DynamoRepository) GetExpenseRecurring(ctx context.Context, expenseRecurringID, username string) (*models.ExpenseRecurring, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(d.tableName),
		Key: map[string]types.AttributeValue{
			"id":       &types.AttributeValueMemberS{Value: expenseRecurringID},
			"username": &types.AttributeValueMemberS{Value: username},
		},
	}

	result, err := d.dynamoClient.GetItem(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("get expense recurring item failed: %v", err)
	}

	if result.Item == nil {
		return nil, models.ErrRecurringExpenseNotFound
	}

	entity := new(ExpenseRecurringEntity)
	err = attributevalue.UnmarshalMap(result.Item, entity)
	if err != nil {
		return nil, fmt.Errorf("unmarshal recurring expense item failed: %v", err)
	}

	return toExpenseRecurringModel(*entity), nil
}

func (d *DynamoRepository) ScanExpensesForDay(ctx context.Context, day int) ([]*models.ExpenseRecurring, error) {
	filter := expression.Name("recurring_day").Equal(expression.Value(day))
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
	entities := make([]*ExpenseRecurringEntity, 0)
	var itemsInQuery []*ExpenseRecurringEntity

	for {
		itemsInQuery = make([]*ExpenseRecurringEntity, 0)
		result, err = d.dynamoClient.Scan(ctx, input)
		if err != nil {
			return nil, err
		}

		if (result.Items == nil || len(result.Items) == 0) && result.LastEvaluatedKey == nil {
			break
		}

		err = attributevalue.UnmarshalListOfMaps(result.Items, &itemsInQuery)
		if err != nil {
			return nil, fmt.Errorf("unmarshal reucrring expenses items failed: %v", err)
		}

		entities = append(entities, itemsInQuery...)
		input.ExclusiveStartKey = result.LastEvaluatedKey

		if result.LastEvaluatedKey == nil {
			break
		}
	}

	if len(entities) == 0 {
		return nil, models.ErrRecurringExpensesNotFound
	}

	return toExpensesRecurringModel(entities), nil
}

func (d *DynamoRepository) BatchDeleteExpenseRecurring(ctx context.Context, log logger.LogAPI, expenseRecurring []*models.ExpenseRecurring) error {
	writeRequests := make([]types.WriteRequest, 0, len(expenseRecurring))

	var idAttrValue types.AttributeValue
	var usernameAttrValue types.AttributeValue
	var err error

	for _, expense := range expenseRecurring {
		idAttrValue, err = attributevalue.Marshal(expense.ID)
		if err != nil {
			return fmt.Errorf("marshal id key failed: %v", err)
		}

		usernameAttrValue, err = attributevalue.Marshal(expense.Username)
		if err != nil {
			return fmt.Errorf("marshal username key failed: %v", err)
		}

		writeRequests = append(writeRequests, types.WriteRequest{
			DeleteRequest: &types.DeleteRequest{
				Key: map[string]types.AttributeValue{
					"id":       idAttrValue,
					"username": usernameAttrValue,
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

func (d *DynamoRepository) DeleteExpenseRecurring(ctx context.Context, expenseRecurringID, username string) error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(d.tableName),
		Key: map[string]types.AttributeValue{
			"id":       &types.AttributeValueMemberS{Value: expenseRecurringID},
			"username": &types.AttributeValueMemberS{Value: username},
		},
	}

	_, err := d.dynamoClient.DeleteItem(ctx, input)
	if err != nil {
		return fmt.Errorf("delete expense recurring item failed: %v", err)
	}

	return nil
}
