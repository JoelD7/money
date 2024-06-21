package expenses_recurring

import (
	"context"
	"fmt"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/env"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
	"time"
)

var (
	tableName               = env.GetString("EXPENSES_RECURRING_TABLE_NAME", "expenses-recurring")
	batchWriteRetries       = env.GetInt("BATCH_WRITE_RETRIES", 3)
	batchWriteBaseDelay     = env.GetInt("BATCH_WRITE_BASE_DELAY_IN_MS", 300)
	batchWriteBackoffFactor = env.GetInt("BATCH_WRITE_BACKOFF_FACTOR", 2)
	dynamoDBMaxBatchWrite   = 25
)

type DynamoRepository struct {
	dynamoClient *dynamodb.Client
}

func NewExpenseRecurringDynamoRepository(dynamoClient *dynamodb.Client) *DynamoRepository {
	return &DynamoRepository{
		dynamoClient: dynamoClient,
	}
}

func (d *DynamoRepository) CreateExpenseRecurring(ctx context.Context, expenseRecurring *models.ExpenseRecurring) (*models.ExpenseRecurring, error) {
	entity := toExpenseRecurringEntity(expenseRecurring)

	item, err := attributevalue.MarshalMap(entity)
	if err != nil {
		return nil, fmt.Errorf("marshal expense recurring entity failed: %v", err)
	}

	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(tableName),
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

	start := 0
	end := dynamoDBMaxBatchWrite
	entitiesInBatch := entities

	if len(entitiesInBatch) > dynamoDBMaxBatchWrite {
		entitiesInBatch = entities[start:end]
	}

	for {
		input := &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]types.WriteRequest{
				tableName: getBatchWriteRequests(entitiesInBatch, log),
			},
		}

		result, err := d.dynamoClient.BatchWriteItem(ctx, input)
		if err != nil {
			return fmt.Errorf("batch write recurring expenses failed: %v", err)
		}

		if result != nil && len(result.UnprocessedItems) > 0 {
			return d.handleBatchWriteRetries(ctx, result.UnprocessedItems)
		}

		if end >= len(entities) {
			break
		}

		start += dynamoDBMaxBatchWrite
		end += dynamoDBMaxBatchWrite

		if len(entities[start:]) > dynamoDBMaxBatchWrite {
			entitiesInBatch = entities[start:end]
			continue
		}

		entitiesInBatch = entities[start:]
	}

	return nil
}

func getBatchWriteRequests(entities []*ExpenseRecurringEntity, log logger.LogAPI) []types.WriteRequest {
	writeRequests := make([]types.WriteRequest, 0, len(entities))

	for _, entity := range entities {
		item, err := attributevalue.MarshalMap(entity)
		if err != nil {
			log.Warning("marshal_recurring_expense_failed", err, []models.LoggerObject{entity})
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

func (d *DynamoRepository) handleBatchWriteRetries(ctx context.Context, unprocessedItems map[string][]types.WriteRequest) error {
	var result *dynamodb.BatchWriteItemOutput
	var err error

	delay := time.Duration(batchWriteBaseDelay) * time.Millisecond

	for i := 0; i < batchWriteRetries; i++ {
		time.Sleep(delay)

		input := &dynamodb.BatchWriteItemInput{
			RequestItems: unprocessedItems,
		}

		result, err = d.dynamoClient.BatchWriteItem(ctx, input)
		if err != nil {
			return fmt.Errorf("batch write recurring expenses failed: %v", err)
		}

		if result != nil && len(result.UnprocessedItems) == 0 {
			return nil
		}

		unprocessedItems = result.UnprocessedItems
		delay *= time.Duration(batchWriteBackoffFactor)
	}

	return nil
}

func (d *DynamoRepository) ScanExpensesForDay(ctx context.Context, day int) ([]*models.ExpenseRecurring, error) {
	filter := expression.Name("recurring_day").Equal(expression.Value(day))
	expr, err := expression.NewBuilder().WithFilter(filter).Build()
	if err != nil {
		return nil, err
	}

	input := &dynamodb.ScanInput{
		TableName:                 aws.String(tableName),
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
