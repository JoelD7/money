package income

import (
	"context"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/env"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

const (
	splitter = ":"
)

var (
	TableName               = env.GetString("INCOME_TABLE_NAME", "income")
	periodUserIncomeIDIndex = "period_user-income_id-index"
)

type DynamoRepository struct {
	dynamoClient *dynamodb.Client
}

func NewDynamoRepository(dynamoClient *dynamodb.Client) *DynamoRepository {
	return &DynamoRepository{dynamoClient: dynamoClient}
}

func (d *DynamoRepository) GetIncomeByPeriod(ctx context.Context, username, periodID string) ([]*models.Income, error) {
	if username == "" {
		return nil, models.ErrMissingUsername
	}

	if periodID == "" {
		return nil, models.ErrMissingPeriod
	}

	periodUser := periodID + splitter + username

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
		return nil, models.ErrIncomeNotFound
	}

	incomeEntities := new([]*incomeEntity)

	err = attributevalue.UnmarshalListOfMaps(result.Items, &incomeEntities)
	if err != nil {
		return nil, err
	}

	return toIncomeModels(*incomeEntities), nil
}

//TODO: el create income debe requerir un period
