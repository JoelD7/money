package savings

import (
	"context"
	"fmt"
	"github.com/JoelD7/money/backend/models"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
	"math/rand"
	"strings"
	"time"
)

const (
	tableName = "savings"
)

type DynamoRepository struct {
	dynamoClient *dynamodb.Client
}

func NewDynamoRepository(dynamoClient *dynamodb.Client) *DynamoRepository {
	return &DynamoRepository{dynamoClient: dynamoClient}
}

func (d *DynamoRepository) GetSavings(ctx context.Context, email string) ([]*models.Saving, error) {
	nameEx := expression.Name("email").Equal(expression.Value(email))

	expr, err := expression.NewBuilder().WithCondition(nameEx).Build()
	if err != nil {
		return nil, err
	}

	input := &dynamodb.QueryInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.Condition(),
		TableName:                 aws.String(tableName),
	}

	result, err := d.dynamoClient.Query(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("query failed: %v", err)
	}

	if result.Items == nil || len(result.Items) == 0 {
		return nil, models.ErrSavingsNotFound
	}

	savings := new([]*models.Saving)

	err = attributevalue.UnmarshalListOfMaps(result.Items, savings)
	if err != nil {
		return nil, fmt.Errorf("unmarshal savings items failed: %v", err)
	}

	return *savings, nil
}

func (d *DynamoRepository) CreateSaving(ctx context.Context, saving *models.Saving) error {
	saving.SavingID = generateSavingID()
	saving.CreatedDate = time.Now()

	item, err := attributevalue.MarshalMap(saving)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(tableName),
	}

	_, err = d.dynamoClient.PutItem(ctx, input)
	if err != nil {
		return fmt.Errorf("put saving item failed: %v", err)
	}

	return nil
}

func (d *DynamoRepository) UpdateSaving(ctx context.Context, saving *models.Saving) error {
	email, err := attributevalue.Marshal(saving.Email)
	if err != nil {
		return fmt.Errorf("marshaling email key: %v", err)
	}

	savingID, err := attributevalue.Marshal(saving.SavingID)
	if err != nil {
		return fmt.Errorf("marshaling saving id key: %v", err)
	}

	attributeValues, err := getAttributeValues(saving)
	if err != nil {
		return fmt.Errorf("getting attribute values: %v", err)
	}

	updateExpression := getUpdateExpression(attributeValues)

	input := &dynamodb.UpdateItemInput{
		Key: map[string]types.AttributeValue{
			"email":     email,
			"saving_id": savingID,
		},
		TableName:                 aws.String(tableName),
		ConditionExpression:       aws.String("attribute_exists(saving_id)"),
		ExpressionAttributeValues: attributeValues,
		UpdateExpression:          updateExpression,
	}

	_, err = d.dynamoClient.UpdateItem(ctx, input)
	if err != nil && strings.Contains(err.Error(), "ConditionalCheckFailedException") {
		return fmt.Errorf("%v: %w", err, models.ErrUpdateSavingNotFound)
	}

	if err != nil {
		return fmt.Errorf("updating saving item: %v", err)
	}

	return nil
}

func getAttributeValues(saving *models.Saving) (map[string]types.AttributeValue, error) {
	m := make(map[string]types.AttributeValue)

	savingGoalID, err := attributevalue.Marshal(saving.SavingGoalID)
	if err != nil {
		return nil, err
	}

	amount, err := attributevalue.Marshal(saving.Amount)
	if err != nil {
		return nil, err
	}

	updatedDate, err := attributevalue.Marshal(time.Now())
	if err != nil {
		return nil, err
	}

	if saving.SavingGoalID != "" {
		m[":saving_goal_id"] = savingGoalID
	}

	if saving.Amount > 0 {
		m[":amount"] = amount
	}

	m[":updated_date"] = updatedDate

	return m, nil
}

func getUpdateExpression(attributeValues map[string]types.AttributeValue) *string {
	attributes := make([]string, 0)

	if _, ok := attributeValues[":saving_goal_id"]; ok {
		attributes = append(attributes, "saving_goal_id = :saving_goal_id")
	}

	if _, ok := attributeValues[":amount"]; ok {
		attributes = append(attributes, "amount = :amount")
	}

	if _, ok := attributeValues[":updated_date"]; ok {
		attributes = append(attributes, "updated_date = :updated_date")
	}

	return aws.String("SET " + strings.Join(attributes, ", "))
}

func generateSavingID() string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, 20)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}

	return "SV" + string(b)
}

func (d *DynamoRepository) DeleteSaving(ctx context.Context, savingID, email string) error {
	emailAtr, err := attributevalue.Marshal(email)
	if err != nil {
		return fmt.Errorf("marshaling email key: %v", err)
	}

	savingIDAtr, err := attributevalue.Marshal(savingID)
	if err != nil {
		return fmt.Errorf("marshaling saving id key: %v", err)
	}

	input := &dynamodb.DeleteItemInput{
		Key: map[string]types.AttributeValue{
			"email":     emailAtr,
			"saving_id": savingIDAtr,
		},
		TableName:           aws.String(tableName),
		ConditionExpression: aws.String("attribute_exists(saving_id)"),
	}

	_, err = d.dynamoClient.DeleteItem(ctx, input)
	if err != nil && strings.Contains(err.Error(), "ConditionalCheckFailedException") {
		return fmt.Errorf("%v: %w", err, models.ErrUpdateSavingNotFound)
	}

	if err != nil {
		return fmt.Errorf("deleting item: %v", err)
	}

	return nil
}
