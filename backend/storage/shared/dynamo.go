package shared

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/JoelD7/money/backend/shared/env"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"time"
)

var (
	batchWriteRetries       = env.GetInt("BATCH_WRITE_RETRIES", 3)
	batchWriteBaseDelay     = env.GetInt("BATCH_WRITE_BASE_DELAY_IN_MS", 300)
	batchWriteBackoffFactor = env.GetInt("BATCH_WRITE_BACKOFF_FACTOR", 2)
	dynamoDBMaxBatchWrite   = env.GetInt("DYNAMODB_MAX_BATCH_WRITE", 25)
)

// BuildPeriodUser builds a combined string of period and username required to identify an item of certain period and user.
func BuildPeriodUser(username, period string) *string {
	p := fmt.Sprintf("%s:%s", period, username)
	return &p
}

// EncodePaginationKey encodes the last evaluated key returned by Dynamo in a string format to be used in the next query
// as the start key.
// The "keyType" parameter should be a pointer to a struct that maps ot the primary key of the table or index in question.
func EncodePaginationKey(lastKey map[string]types.AttributeValue, keyType interface{}) (string, error) {
	if len(lastKey) == 0 {
		return "", nil
	}

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
func DecodePaginationKey(startKey string, keyType interface{}) (map[string]types.AttributeValue, error) {
	decoded, err := base64.URLEncoding.DecodeString(startKey)
	if err != nil {
		return nil, fmt.Errorf("decoding last key: %v", err)
	}

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

//func BatchWrite(input *dynamodb.BatchWriteItemInput) error {
//	start := 0
//	end := dynamoDBMaxBatchWrite
//	requestItemsInBatch := copyRequestItems(input.RequestItems)
//
//	if len(requestItemsInBatch) > dynamoDBMaxBatchWrite {
//		requestItemsInBatch = entities[start:end]
//	}
//
//	for {
//		input := &dynamodb.BatchWriteItemInput{
//			RequestItems: map[string][]types.WriteRequest{
//				tableName: getBatchWriteRequests(requestItemsInBatch, log),
//			},
//		}
//
//		result, err := d.dynamoClient.BatchWriteItem(ctx, input)
//		if err != nil {
//			return fmt.Errorf("batch write recurring expenses failed: %v", err)
//		}
//
//		if result != nil && len(result.UnprocessedItems) > 0 {
//			return shared.HandleBatchWriteRetries(ctx, d.dynamoClient, result.UnprocessedItems)
//		}
//
//		if end >= len(entities) {
//			break
//		}
//
//		start += dynamoDBMaxBatchWrite
//		end += dynamoDBMaxBatchWrite
//
//		if len(entities[start:]) > dynamoDBMaxBatchWrite {
//			requestItemsInBatch = entities[start:end]
//			continue
//		}
//
//		requestItemsInBatch = entities[start:]
//	}
//}

func copyRequestItems(requestItems map[string][]types.WriteRequest) map[string][]types.WriteRequest {
	c := make(map[string][]types.WriteRequest)
	for key, value := range requestItems {
		c[key] = value
	}

	return c
}

func splitRequestItems(requestItems map[string][]types.WriteRequest) []map[string][]types.WriteRequest {
	var result []map[string][]types.WriteRequest
	batch := make(map[string][]types.WriteRequest)

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
				}

				break
			}

			if remSlots < len(itemsWoProcess) {
				batch[table] = itemsWoProcess[:remSlots]
				itemsWoProcess = itemsWoProcess[remSlots:]
				result = append(result, batch)
				remSlots = dynamoDBMaxBatchWrite
				continue
			}
		}
	}

	return result
}

// HandleBatchWriteRetries is a helper function to handle retries for batch write operations in DynamoDB.
func HandleBatchWriteRetries(ctx context.Context, d *dynamodb.Client, unprocessedItems map[string][]types.WriteRequest) error {
	var result *dynamodb.BatchWriteItemOutput
	var err error

	delay := time.Duration(batchWriteBaseDelay) * time.Millisecond

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
