package period

import (
	"context"
	"fmt"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/env"
	"github.com/JoelD7/money/backend/storage/dynamo"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"regexp"
	"strings"
	"time"
	"unicode"
)

const (
	defaultPageSize          = 10
	conditionalFailedKeyword = "ConditionalCheckFailed"
	periodPrefix             = "PRD"
)

var (
	cancelReasonRegex = regexp.MustCompile("\\[[a-zA-Z,\\s]+\\]")
)

type DynamoRepository struct {
	dynamoClient              *dynamodb.Client
	periodTableName           string
	uniquePeriodNameTableName string
}

func NewDynamoRepository(dynamoClient *dynamodb.Client, periodTableName, uniquePeriodTableName string) (*DynamoRepository, error) {
	d := &DynamoRepository{dynamoClient: dynamoClient}

	periodTableNameEnv := env.GetString("PERIOD_TABLE_NAME", "")
	uniquePeriodTableNameEnv := env.GetString("UNIQUE_PERIOD_TABLE_NAME", "")

	err := validateParams(periodTableName, uniquePeriodTableName, periodTableNameEnv, uniquePeriodTableNameEnv)
	if err != nil {
		return nil, fmt.Errorf("initialize period dynamo repository failed: %v", err)
	}

	d.periodTableName = periodTableName
	if d.periodTableName == "" {
		d.periodTableName = periodTableNameEnv
	}

	d.uniquePeriodNameTableName = uniquePeriodTableName
	if d.uniquePeriodNameTableName == "" {
		d.uniquePeriodNameTableName = uniquePeriodTableNameEnv
	}

	return d, nil
}

func validateParams(periodTableName, uniquePeriodTableName, periodTableNameEnv, uniquePeriodTableNameEnv string) error {
	if periodTableName == "" && periodTableNameEnv == "" {
		return fmt.Errorf("period table name is required")
	}

	if uniquePeriodTableName == "" && uniquePeriodTableNameEnv == "" {
		return fmt.Errorf("unique period table name is required")
	}

	return nil
}

func (d *DynamoRepository) CreatePeriod(ctx context.Context, period *models.Period) (*models.Period, error) {
	period.ID = *period.Name
	periodEnt := toPeriodEntity(*period)

	attrValue, err := attributevalue.MarshalMap(periodEnt)
	if err != nil {
		return nil, fmt.Errorf("marshal period item failed: %v", err)
	}

	cond := expression.Name("period").AttributeNotExists()

	expr, err := expression.NewBuilder().WithCondition(cond).Build()
	if err != nil {
		return nil, fmt.Errorf("build expression failed: %v", err)
	}

	input := &dynamodb.PutItemInput{
		TableName:                aws.String(d.periodTableName),
		Item:                     attrValue,
		ConditionExpression:      expr.Condition(),
		ExpressionAttributeNames: expr.Names(),
	}

	_, err = d.dynamoClient.PutItem(ctx, input)
	if err != nil && strings.Contains(err.Error(), conditionalFailedKeyword) {
		return nil, fmt.Errorf("%v: %w", err, models.ErrPeriodNameIsTaken)
	}

	if err != nil {
		return nil, err
	}

	return period, nil
}

func (d *DynamoRepository) UpdatePeriod(ctx context.Context, period *models.Period) error {
	periodEnt := toPeriodEntity(*period)

	periodEnt.UpdatedDate = time.Now()

	uPeriodName := &uniquePeriodNameEntity{
		Name:     *periodEnt.Name,
		Username: periodEnt.Username,
	}

	periodAv, err := attributevalue.MarshalMap(periodEnt)
	if err != nil {
		return fmt.Errorf("marshaling period to attribute value: %v", err)
	}

	uPeriodNameAv, err := attributevalue.MarshalMap(uPeriodName)
	if err != nil {
		return fmt.Errorf("marshaling unique period name to attribute value failed: %v", err)
	}

	periodExistsCond := expression.Name("period").AttributeExists()
	periodNameNotTakenCond := expression.Name("name").AttributeNotExists().And(expression.Name("username").AttributeNotExists())

	periodTableExpr, err := expression.NewBuilder().WithCondition(periodExistsCond).Build()
	if err != nil {
		return fmt.Errorf("building period table expression failed: %v", err)
	}

	uniquePeriodNameTableExpr, err := expression.NewBuilder().WithCondition(periodNameNotTakenCond).Build()
	if err != nil {
		return fmt.Errorf("building unique period name table expression failed: %v", err)
	}

	errByCondition := map[string]error{
		*periodTableExpr.Condition():           models.ErrUpdatePeriodNotFound,
		*uniquePeriodNameTableExpr.Condition(): models.ErrPeriodNameIsTaken,
	}

	transactItems := []types.TransactWriteItem{
		{
			Put: &types.Put{
				TableName:                aws.String(d.periodTableName),
				ConditionExpression:      periodTableExpr.Condition(),
				ExpressionAttributeNames: periodTableExpr.Names(),
				Item:                     periodAv,
			},
		},
		{
			Put: &types.Put{
				TableName:                aws.String(d.uniquePeriodNameTableName),
				ConditionExpression:      uniquePeriodNameTableExpr.Condition(),
				ExpressionAttributeNames: uniquePeriodNameTableExpr.Names(),
				Item:                     uPeriodNameAv,
			},
		},
	}

	input := &dynamodb.TransactWriteItemsInput{
		TransactItems: transactItems,
	}

	_, err = d.dynamoClient.TransactWriteItems(ctx, input)
	if err != nil {
		return handleUpdatePeriodError(transactItems, errByCondition, err)
	}

	return nil
}

func handleUpdatePeriodError(transactItems []types.TransactWriteItem, errByCondition map[string]error, err error) error {
	defaultErr := fmt.Errorf("updating period item: %v", err)

	if !strings.Contains(err.Error(), conditionalFailedKeyword) {
		return defaultErr
	}

	cancelReasons := extractCancellationReason(err)
	if len(cancelReasons) == 0 {
		return defaultErr
	}

	conditionByPos := make(map[int]string)

	for i, item := range transactItems {
		conditionByPos[i] = getConditionExpression(item)
	}

	errPosition := -1

	for i, reason := range cancelReasons {
		if reason == conditionalFailedKeyword {
			errPosition = i
			break
		}
	}

	failedCondition := conditionByPos[errPosition]

	conditionErr, ok := errByCondition[failedCondition]
	if !ok {
		return defaultErr
	}

	return fmt.Errorf("%v: %w", err, conditionErr)
}

// extractCancellationReason extracts the cancellation reason array from the error.
//
// When a transaction fails, the error message contains an ordered list of errors for each item in the request which
// caused the transaction to get cancelled in the form of "[Reason, None, ...]". If a transact item did not fail, the
// error in the list will contain "None" instead of a reason.
// Read more about this here: https://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_CancellationReason.html
func extractCancellationReason(err error) []string {
	condReason := cancelReasonRegex.FindAllString(err.Error(), -1)

	splitFunc := func(c rune) bool {
		return unicode.IsSpace(c) || c == ','
	}

	for _, part := range condReason {
		// part has the form "[Reason, None, ...]"
		if strings.Contains(part, conditionalFailedKeyword) {
			return strings.FieldsFunc(strings.Trim(part, "[]"), splitFunc)
		}
	}

	return nil
}

func getConditionExpression(item types.TransactWriteItem) string {
	if item.Put != nil {
		return *item.Put.ConditionExpression
	}

	if item.Delete != nil {
		return *item.Delete.ConditionExpression
	}

	if item.Update != nil {
		return *item.Update.ConditionExpression
	}

	return ""
}

func (d *DynamoRepository) GetPeriod(ctx context.Context, username, period string) (*models.Period, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(d.periodTableName),
		Key: map[string]types.AttributeValue{
			"username": &types.AttributeValueMemberS{Value: username},
			"period":   &types.AttributeValueMemberS{Value: period},
		},
	}

	result, err := d.dynamoClient.GetItem(ctx, input)
	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		return nil, models.ErrPeriodNotFound
	}

	periodStruct := periodEntity{}

	err = attributevalue.UnmarshalMap(result.Item, &periodStruct)
	if err != nil {
		return nil, fmt.Errorf("unmarshal period item failed: %v", err)
	}

	return toPeriodModel(periodStruct), nil
}

func (d *DynamoRepository) GetPeriods(ctx context.Context, username, startKey string, pageSize int) ([]*models.Period, string, error) {
	keyConditionExpression := expression.Key("username").Equal(expression.Value(username))

	conditionBuilder := expression.NewBuilder().WithKeyCondition(keyConditionExpression)

	expr, err := conditionBuilder.Build()
	if err != nil {
		return nil, "", fmt.Errorf("build expression failed: %v", err)
	}

	var decodedStartKey map[string]types.AttributeValue

	if startKey != "" {
		decodedStartKey, err = dynamo.DecodePaginationKey(startKey)
		if err != nil {
			return nil, "", fmt.Errorf("%v: %w", err, models.ErrInvalidStartKey)
		}
	}

	input := &dynamodb.QueryInput{
		TableName:                 aws.String(d.periodTableName),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ExclusiveStartKey:         decodedStartKey,
		Limit:                     getPageSize(pageSize),
	}

	result, err := d.dynamoClient.Query(ctx, input)
	if err != nil {
		return nil, "", err
	}

	if (result.Items == nil || len(result.Items) == 0) && startKey == "" {
		return nil, "", models.ErrPeriodsNotFound
	}

	periods := make([]periodEntity, 0, len(result.Items))

	err = attributevalue.UnmarshalListOfMaps(result.Items, &periods)
	if err != nil {
		return nil, "", fmt.Errorf("unmarshal periods failed: %v", err)
	}

	nextKey, err := dynamo.EncodePaginationKey(result.LastEvaluatedKey)
	if err != nil {
		return nil, "", err
	}

	return toPeriodModels(periods), nextKey, nil
}

func (d *DynamoRepository) GetLastPeriod(ctx context.Context, username string) (*models.Period, error) {
	keyConditionExpression := expression.Key("username").Equal(expression.Value(username))

	conditionBuilder := expression.NewBuilder().WithKeyCondition(keyConditionExpression)

	expr, err := conditionBuilder.Build()
	if err != nil {
		return nil, fmt.Errorf("build expression failed: %v", err)
	}

	input := &dynamodb.QueryInput{
		TableName:                 aws.String(d.periodTableName),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ScanIndexForward:          aws.Bool(false),
	}

	result, err := d.dynamoClient.Query(ctx, input)
	if err != nil {
		return nil, err
	}

	if result.Items == nil || len(result.Items) == 0 {
		return nil, models.ErrPeriodsNotFound
	}

	periodStruct := periodEntity{}

	err = attributevalue.UnmarshalMap(result.Items[0], &periodStruct)
	if err != nil {
		return nil, fmt.Errorf("unmarshal period item failed: %v", err)
	}

	return toPeriodModel(periodStruct), nil
}

func getPageSize(pageSize int) *int32 {
	if pageSize == 0 {
		return aws.Int32(defaultPageSize)
	}

	return aws.Int32(int32(pageSize))
}

func (d *DynamoRepository) DeletePeriod(ctx context.Context, periodID, username string) error {
	period, err := d.GetPeriod(ctx, username, periodID)
	if err != nil {
		return fmt.Errorf("could not get period to delete: %w", err)
	}

	input := &dynamodb.TransactWriteItemsInput{
		TransactItems: []types.TransactWriteItem{
			{
				Delete: &types.Delete{
					TableName: aws.String(d.periodTableName),
					Key: map[string]types.AttributeValue{
						"username": &types.AttributeValueMemberS{Value: username},
						"period":   &types.AttributeValueMemberS{Value: periodID},
					},
				},
			},
			{
				Delete: &types.Delete{
					TableName: aws.String(d.uniquePeriodNameTableName),
					Key: map[string]types.AttributeValue{
						"name":     &types.AttributeValueMemberS{Value: *period.Name},
						"username": &types.AttributeValueMemberS{Value: username},
					},
				},
			},
		},
	}

	_, err = d.dynamoClient.TransactWriteItems(ctx, input)

	return err
}

func (d *DynamoRepository) BatchDeletePeriods(ctx context.Context, periods []*models.Period) error {
	periodWriteRequests, uniquePeriodWriteRequests, err := getBatchPeriodDeleteRequests(periods)
	if err != nil {
		return err
	}

	input := &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]types.WriteRequest{
			d.periodTableName:           periodWriteRequests,
			d.uniquePeriodNameTableName: uniquePeriodWriteRequests,
		},
	}

	return dynamo.BatchWrite(ctx, d.dynamoClient, input)
}

func getBatchPeriodDeleteRequests(periods []*models.Period) ([]types.WriteRequest, []types.WriteRequest, error) {
	periodWriteRequests := make([]types.WriteRequest, 0, len(periods))
	uniquePeriodWriteRequests := make([]types.WriteRequest, 0, len(periods))

	var periodNameAV types.AttributeValue
	var usernameAV types.AttributeValue
	var err error

	for _, p := range periods {
		periodNameAV, err = attributevalue.Marshal(p.ID)
		if err != nil {
			return nil, nil, fmt.Errorf("marshal id key failed: %v", err)
		}

		usernameAV, err = attributevalue.Marshal(p.Username)
		if err != nil {
			return nil, nil, fmt.Errorf("marshal username key failed: %v", err)
		}

		periodWriteRequests = append(periodWriteRequests, types.WriteRequest{
			DeleteRequest: &types.DeleteRequest{
				Key: map[string]types.AttributeValue{
					"period":   periodNameAV,
					"username": usernameAV,
				},
			},
		})

		uniquePeriodWriteRequests = append(uniquePeriodWriteRequests, types.WriteRequest{
			DeleteRequest: &types.DeleteRequest{
				Key: map[string]types.AttributeValue{
					"name":     periodNameAV,
					"username": usernameAV,
				},
			},
		})
	}

	return periodWriteRequests, uniquePeriodWriteRequests, nil
}
