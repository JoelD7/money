package expenses_recurring

import (
	"context"
	"fmt"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/env"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go/aws"
)

var (
	tableName = env.GetString("EXPENSES_RECURRING_TABLE_NAME", "expenses-recurring")
)

type ExpenseRecurringDynamoRepository struct {
	dynamoClient *dynamodb.Client
}

func NewExpenseRecurringDynamoRepository(dynamoClient *dynamodb.Client) *ExpenseRecurringDynamoRepository {
	return &ExpenseRecurringDynamoRepository{
		dynamoClient: dynamoClient,
	}
}

func (d *ExpenseRecurringDynamoRepository) ScanExpensesForDay(ctx context.Context, day int) ([]*models.ExpenseRecurring, error) {
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
	itemsInQuery := make([]*ExpenseRecurringEntity, 0)

	for {
		result, err = d.dynamoClient.Scan(ctx, input)
		if err != nil {
			return nil, err
		}

		if (result.Items == nil || len(result.Items) == 0) && input.ExclusiveStartKey == nil {
			break
		}

		err = attributevalue.UnmarshalListOfMaps(result.Items, &itemsInQuery)
		if err != nil {
			return nil, fmt.Errorf("unmarshal reucrring expenses items failed: %v", err)
		}

		entities = append(entities, itemsInQuery...)
	}

	if len(entities) == 0 {
		return nil, models.ErrRecurringExpensesNotFound
	}

	return toExpensesRecurringModel(entities), nil
}
