package savingoal

import (
	"context"
	"fmt"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/env"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoRepository struct {
	dynamoClient *dynamodb.Client
	tableName    string
}

func NewDynamoRepository(dynamoClient *dynamodb.Client, tableName string) (*DynamoRepository, error) {
	d := &DynamoRepository{dynamoClient: dynamoClient}
	tableNameEnv := env.GetString("SAVING_GOALS_TABLE_NAME", "")

	if tableNameEnv == "" && tableName == "" {
		return nil, fmt.Errorf("initialize saving goal dynamo repository failed: table name is required")
	}

	d.tableName = tableName
	if d.tableName == "" {
		d.tableName = tableNameEnv
	}

	return d, nil
}

func (d *DynamoRepository) GetSavingGoal(ctx context.Context, username, savingGoalID string) (*models.SavingGoal, error) {
	userKey, err := attributevalue.Marshal(username)
	if err != nil {
		return nil, err
	}

	savingGoalIDKey, err := attributevalue.Marshal(savingGoalID)
	if err != nil {
		return nil, err
	}

	input := &dynamodb.GetItemInput{
		TableName: aws.String(d.tableName),
		Key: map[string]types.AttributeValue{
			"username":       userKey,
			"saving_goal_id": savingGoalIDKey,
		},
	}

	result, err := d.dynamoClient.GetItem(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("get saving goal item failed: %v", err)
	}

	savingGoal := new(savingGoalEntity)

	err = attributevalue.UnmarshalMap(result.Item, savingGoal)
	if err != nil {
		return nil, fmt.Errorf("unmarshal saving goal item failed: %v", err)
	}

	return toSavingGoalModel(savingGoal), nil
}

func (d *DynamoRepository) GetSavingGoals(ctx context.Context, username string) ([]*models.SavingGoal, error) {
	nameEx := expression.Name("username").Equal(expression.Value(username))

	expr, err := expression.NewBuilder().WithCondition(nameEx).Build()
	if err != nil {
		return nil, err
	}

	input := &dynamodb.QueryInput{
		TableName:                 aws.String(d.tableName),
		KeyConditionExpression:    expr.Condition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	}

	result, err := d.dynamoClient.Query(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("get saving goals failed: %v", err)
	}

	savingGoalsEntities := new([]*savingGoalEntity)

	err = attributevalue.UnmarshalListOfMaps(result.Items, savingGoalsEntities)
	if err != nil {
		return nil, fmt.Errorf("unmarshal saving goal items failed: %v", err)
	}

	return toSavingGoalModels(*savingGoalsEntities), nil
}
