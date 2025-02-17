package dynamo

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/env"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const (
	defaultPageSize = 10
)

// GenerateID generates a hex-based random unique ID with the given prefix
func GenerateID(prefix string) string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, 20)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}

	return prefix + string(b)
}

// BuildPeriodUser builds a combined string of period and username required to identify an item of certain period and user.
func BuildPeriodUser(username, period string) *string {
	p := fmt.Sprintf("%s:%s", period, username)
	return &p
}

// BuildAmountKey builds the amount sort key, which is a combined string of the amount and the item ID.
// Example -> Input: 1234.56, abc123 -> Output: 000001234.56:abc123
func BuildAmountKey(amount float64, id string) string {
	amountStr := fmt.Sprintf("%0.2f", amount)
	return fmt.Sprintf("%012s:%s", amountStr, id)
}

// BuildNameKey builds the name sort key, which is a combined string of the name and the item ID.
func BuildNameKey(name, id string) string {
	return fmt.Sprintf("%s:%s", name, id)
}

// EncodePaginationKey encodes the last evaluated key returned by Dynamo in a string format to be used in the next query
// as the start key.
// The "keyType" parameter should be a pointer to a struct that maps ot the primary key of the table or index in question.
func EncodePaginationKey(lastKey map[string]types.AttributeValue) (string, error) {
	if len(lastKey) == 0 {
		return "", nil
	}

	var keyType interface{}

	err := attributevalue.UnmarshalMap(lastKey, &keyType)
	if err != nil {
		return "", fmt.Errorf("unmarshalling lastKey map: %v", err)
	}

	data, err := json.Marshal(keyType)
	if err != nil {
		return "", fmt.Errorf("json marshalling primary key: %v", err)
	}

	encoded := base64.URLEncoding.EncodeToString(data)

	return encoded, nil
}

// DecodePaginationKey parses the start key string into a map of attribute values to be used as ExclusiveStartKey in a paginated
// query.
// The "keyType" parameter should be a pointer to a struct that maps ot the primary key of the table or index in question.
func DecodePaginationKey(startKey string) (map[string]types.AttributeValue, error) {
	decoded, err := base64.URLEncoding.DecodeString(startKey)
	if err != nil {
		return nil, fmt.Errorf("decoding last key: %v", err)
	}

	var keyType interface{}

	err = json.Unmarshal(decoded, &keyType)
	if err != nil {
		return nil, fmt.Errorf("json unmarshalling primary key: %v", err)
	}

	exclusiveStartKey, err := attributevalue.MarshalMap(keyType)
	if err != nil {
		return nil, fmt.Errorf("marshalling to map of attribute value: %v", err)
	}

	return exclusiveStartKey, nil
}

// BatchWrite is a helper function to handle batch write operations in DynamoDB. Prevents more than 25 items per batch by
// splitting the request items in multiple batches and also handles retries in case of unprocessed items.
func BatchWrite(ctx context.Context, dynamoClient *dynamodb.Client, input *dynamodb.BatchWriteItemInput) error {
	requestItemsBatches := splitRequestItems(input.RequestItems)

	for _, requestItems := range requestItemsBatches {
		input.RequestItems = requestItems

		result, err := dynamoClient.BatchWriteItem(ctx, input)
		if err != nil {
			return fmt.Errorf("batch write recurring expenses failed: %v", err)
		}

		if result != nil && len(result.UnprocessedItems) > 0 {
			return handleBatchWriteRetries(ctx, dynamoClient, result.UnprocessedItems)
		}
	}

	return nil
}

// splitRequestItems splits batch write's request items in batches of 25, as this is DynamoDB's current limit of batch items per request.
// The function works regardless of how many tables are in the request, meaning that a batch may be composed of a single table with
// 25 request items or a group of tables that add up to 25 items.
func splitRequestItems(requestItems map[string][]types.WriteRequest) []map[string][]types.WriteRequest {
	var result []map[string][]types.WriteRequest
	batch := make(map[string][]types.WriteRequest)
	dynamoDBMaxBatchWrite := env.GetInt("DYNAMODB_MAX_BATCH_WRITE", 25)

	remSlots := dynamoDBMaxBatchWrite               // available slots in batch
	itemsWoProcess := make([]types.WriteRequest, 0) // items that could not be processed in the current batch

	for table, items := range requestItems {
		itemsWoProcess = items

		for {
			if remSlots == 0 && len(itemsWoProcess) == 0 {
				break
			}

			if remSlots > 0 && remSlots >= len(itemsWoProcess) {
				batch[table] = itemsWoProcess
				remSlots -= len(itemsWoProcess)

				if remSlots == 0 {
					result = append(result, batch)
					batch = make(map[string][]types.WriteRequest)
					remSlots = dynamoDBMaxBatchWrite
				}

				break
			}

			if remSlots < len(itemsWoProcess) {
				batch[table] = itemsWoProcess[:remSlots]
				itemsWoProcess = itemsWoProcess[remSlots:]
				result = append(result, batch)
				batch = make(map[string][]types.WriteRequest)
				remSlots = dynamoDBMaxBatchWrite
				continue
			}
		}

	}

	if len(batch) > 0 {
		result = append(result, batch)
		batch = make(map[string][]types.WriteRequest)
	}

	return result
}

// handleBatchWriteRetries is a helper function to handle retries for batch write operations in DynamoDB.
func handleBatchWriteRetries(ctx context.Context, d *dynamodb.Client, unprocessedItems map[string][]types.WriteRequest) error {
	var result *dynamodb.BatchWriteItemOutput
	var err error
	batchWriteBaseDelay := env.GetInt("BATCH_WRITE_BASE_DELAY_IN_MS", 300)

	delay := time.Duration(batchWriteBaseDelay) * time.Millisecond
	batchWriteRetries := env.GetInt("BATCH_WRITE_RETRIES", 3)
	batchWriteBackoffFactor := env.GetInt("BATCH_WRITE_BACKOFF_FACTOR", 2)

	for i := 0; i < batchWriteRetries; i++ {
		time.Sleep(delay)

		input := &dynamodb.BatchWriteItemInput{
			RequestItems: unprocessedItems,
		}

		result, err = d.BatchWriteItem(ctx, input)
		if err != nil {
			return fmt.Errorf("batch write failed: %v", err)
		}

		if result != nil && len(result.UnprocessedItems) == 0 {
			return nil
		}

		unprocessedItems = result.UnprocessedItems
		delay *= time.Duration(batchWriteBackoffFactor)
	}

	return nil
}

func InitClient(ctx context.Context) *dynamodb.Client {
	awsRegion := env.GetString("AWS_REGION", "")

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(awsRegion),
		config.WithLogger(logger.NewLogstashDynamo()),
		config.WithClientLogMode(aws.LogRequestWithBody|aws.LogResponseWithBody))
	if err != nil {
		panic(err)
	}

	return dynamodb.NewFromConfig(cfg)
}

// SetExclusiveStartKey sets the ExclusiveStartKey in a DynamoDB query input based on the provided startKey string.
func SetExclusiveStartKey(startKey string, input *dynamodb.QueryInput) error {
	if startKey == "" {
		return nil
	}

	decodedStartKey, err := DecodePaginationKey(startKey)
	if err != nil {
		return fmt.Errorf("%v: %w", err, models.ErrInvalidStartKey)
	}

	input.ExclusiveStartKey = decodedStartKey

	return nil
}

func GetPageSize(pageSize int) *int32 {
	if pageSize == 0 {
		return aws.Int32(defaultPageSize)
	}

	return aws.Int32(int32(pageSize))
}
