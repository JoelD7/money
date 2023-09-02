package income

import (
	"context"
	"errors"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/env"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go/aws"
)

const (
	splitter = ":"
)

var (
	TableName               = env.GetString("INCOME_TABLE_NAME", "income")
	periodUserIncomeIDIndex = "period_user-income_id-index"

	ErrNotFound    = errors.New("income not found")
	ErrEmptyUserID = errors.New("empty userID")
	ErrEmptyPeriod = errors.New("empty period")
)

type DynamoRepository struct {
	dynamoClient *dynamodb.Client
}

func NewDynamoRepository(dynamoClient *dynamodb.Client) *DynamoRepository {
	return &DynamoRepository{dynamoClient: dynamoClient}
}

func (d *DynamoRepository) getIncomeByPeriod(ctx context.Context, userID, periodID string) ([]*models.Income, error) {
	if userID == "" {
		return nil, ErrEmptyUserID
	}

	if periodID == "" {
		return nil, ErrEmptyPeriod
	}

	periodUser := periodID + splitter + userID

	nameEx := expression.Name("period_user").Equal(expression.Value(periodUser))

	expr, err := expression.NewBuilder().WithCondition(nameEx).Build()
	if err != nil {
		return nil, err
	}

	input := &dynamodb.QueryInput{
		TableName:                 aws.String(TableName),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.Condition(),
		IndexName:                 aws.String(periodUserIncomeIDIndex),
	}

	result, err := d.dynamoClient.Query(ctx, input)
	if err != nil {
		return nil, err
	}

	if result.Items == nil || len(result.Items) == 0 {
		return nil, ErrNotFound
	}

	income := make([]*models.Income, 0)

	err = attributevalue.UnmarshalListOfMaps(result.Items, &income)
	if err != nil {
		return nil, err
	}

	return income, nil
}

//TODO: el create income debe requerir un period
