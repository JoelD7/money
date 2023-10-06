package savingoal

import (
	"context"
	"fmt"
	"github.com/JoelD7/money/backend/models"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
)

const (
	tableName = "saving-goals"
)

type DynamoRepository struct {
	dynamoClient *dynamodb.Client
}

func NewDynamoRepository(dynamoClient *dynamodb.Client) *DynamoRepository {
	return &DynamoRepository{dynamoClient: dynamoClient}
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
		TableName: aws.String(tableName),
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
