package income

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
	"strings"
)

const (
	splitter          = ":"
	conditionFailedEx = "ConditionalCheckFailedException"
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

func (d *DynamoRepository) CreateIncome(ctx context.Context, income *models.Income) (*models.Income, error) {
	incomeEnt := toIncomeEntity(income)
	incomeEnt.PeriodUser = shared.BuildPeriodUser(income.Username, *income.Period)

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
		TableName:                aws.String(TableName),
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
