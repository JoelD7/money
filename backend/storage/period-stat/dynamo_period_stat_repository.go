package period_stat

import (
	"context"
	"fmt"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/storage/dynamo"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoRepository struct {
	dynamoClient      *dynamodb.Client
	tableName         string
	categoryUserIndex string
}

func NewDynamoRepository(dynamoClient *dynamodb.Client, tableName, categoryUserIndex string) (*DynamoRepository, error) {
	if tableName == "" || categoryUserIndex == "" {
		return nil, fmt.Errorf("storage: table name and category user index are required")
	}

	return &DynamoRepository{
		dynamoClient:      dynamoClient,
		tableName:         tableName,
		categoryUserIndex: categoryUserIndex,
	}, nil
}

func (d *DynamoRepository) GetPeriodStat(ctx context.Context, period, username, categoryID string) (*models.PeriodStat, error) {
	periodUser := dynamo.BuildPeriodUser(username, period)
	if periodUser == nil {
		return nil, fmt.Errorf("storage: unable to build period user")
	}

	input := dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"period_user": &types.AttributeValueMemberS{Value: *periodUser},
			"category_id": &types.AttributeValueMemberS{Value: categoryID},
		},
		TableName: aws.String(d.tableName),
	}

	result, err := d.dynamoClient.GetItem(ctx, &input)
	if err != nil {
		return nil, fmt.Errorf("storage: get period stat: %v", err)
	}

	var periodStat periodStatEntity
	err = attributevalue.UnmarshalMap(result.Item, &periodStat)
	if err != nil {
		return nil, fmt.Errorf("storage: unmarshalling period stat: %v", err)
	}

	return toPeriodStatModel(periodStat), nil
}
