package savingoal

import (
	"context"
	"fmt"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/storage/dynamo"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const savingGoalIDPrefix = "SVG"

type DynamoRepository struct {
	dynamoClient          *dynamodb.Client
	tableName             string
	usernameDeadlineIndex string
	usernameTargetIndex   string
}

func NewDynamoRepository(dynamoClient *dynamodb.Client, envConfig *models.EnvironmentConfiguration) (*DynamoRepository, error) {
	d := &DynamoRepository{dynamoClient: dynamoClient}

	err := validateParams(envConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize saving goals dynamo repository: %v", err)
	}

	d.tableName = envConfig.SavingGoalsTable
	d.usernameDeadlineIndex = envConfig.UsernameDeadlineIndex
	d.usernameTargetIndex = envConfig.UsernameTargetIndex

	return d, nil
}

func validateParams(envConfig *models.EnvironmentConfiguration) error {
	if envConfig.SavingGoalsTable == "" {
		return fmt.Errorf("table name is required")
	}

	if envConfig.UsernameDeadlineIndex == "" {
		return fmt.Errorf("username deadline index is required")
	}

	if envConfig.UsernameTargetIndex == "" {
		return fmt.Errorf("username target index is required")
	}

	return nil
}

func (d *DynamoRepository) CreateSavingGoal(ctx context.Context, savingGoal *models.SavingGoal) (*models.SavingGoal, error) {
	savingGoal.SavingGoalID = dynamo.GenerateID(savingGoalIDPrefix)
	entity := toSavingGoalEntity(savingGoal)

	av, err := attributevalue.MarshalMap(entity)
	if err != nil {
		return nil, fmt.Errorf("marshal saving goal item failed: %v", err)
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(d.tableName),
		Item:      av,
	}

	_, err = d.dynamoClient.PutItem(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("create saving goal item failed: %v", err)
	}

	return toSavingGoalModel(entity), nil
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

	if result.Item == nil {
		return nil, models.ErrSavingGoalNotFound
	}

	savingGoal := new(savingGoalEntity)

	err = attributevalue.UnmarshalMap(result.Item, savingGoal)
	if err != nil {
		return nil, fmt.Errorf("unmarshal saving goal item failed: %v", err)
	}

	return toSavingGoalModel(savingGoal), nil
}

func (d *DynamoRepository) GetSavingGoals(ctx context.Context, username string, params *models.QueryParameters) ([]*models.SavingGoal, string, error) {
	input, err := d.buildQueryInput(username, params)
	if err != nil {
		return nil, "", err
	}

	result, err := d.dynamoClient.Query(ctx, input)
	if err != nil {
		return nil, "", fmt.Errorf("get saving goals failed: %v", err)
	}

	if len(result.Items) == 0 {
		return nil, "", models.ErrSavingGoalsNotFound
	}

	savingGoalsEntities := new([]*savingGoalEntity)

	err = attributevalue.UnmarshalListOfMaps(result.Items, savingGoalsEntities)
	if err != nil {
		return nil, "", fmt.Errorf("unmarshal saving goal items failed: %v", err)
	}

	nextKey, err := dynamo.EncodePaginationKey(result.LastEvaluatedKey)
	if err != nil {
		return nil, "", err
	}

	return toSavingGoalModels(*savingGoalsEntities), nextKey, nil
}

func (d *DynamoRepository) buildQueryInput(username string, params *models.QueryParameters) (*dynamodb.QueryInput, error) {
	var err error

	input := &dynamodb.QueryInput{
		TableName: aws.String(d.tableName),
		Limit:     dynamo.GetPageSize(params.PageSize),
	}

	if params.SortType == string(models.SortOrderDescending) {
		input.ScanIndexForward = aws.Bool(false)
	}

	keyConditionEx := d.setQueryIndex(input, username, params)

	err = dynamo.SetExclusiveStartKey(params.StartKey, input)
	if err != nil {
		return nil, err
	}

	conditionBuilder := expression.NewBuilder().WithCondition(keyConditionEx)

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

func (d *DynamoRepository) setQueryIndex(input *dynamodb.QueryInput, username string, params *models.QueryParameters) expression.ConditionBuilder {
	keyConditionEx := expression.Name("username").Equal(expression.Value(username))

	if params.SortBy == string(models.SortParamDeadline) {
		input.IndexName = aws.String(d.usernameDeadlineIndex)
	}

	if params.SortBy == string(models.SortParamTarget) {
		input.IndexName = aws.String(d.usernameTargetIndex)
	}

	return keyConditionEx
}

func (d *DynamoRepository) DeleteSavingGoal(ctx context.Context, username, savingGoalID string) error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(d.tableName),
		Key: map[string]types.AttributeValue{
			"username":       &types.AttributeValueMemberS{Value: username},
			"saving_goal_id": &types.AttributeValueMemberS{Value: savingGoalID},
		},
	}

	_, err := d.dynamoClient.DeleteItem(ctx, input)
	if err != nil {
		return fmt.Errorf("delete saving goal item failed: %v", err)
	}

	return nil
}
