package dynamo

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSplitRequestItems(t *testing.T) {
	c := require.New(t)

	requestItems := map[string][]types.WriteRequest{
		"table1": buildSampleRequestItems(27),
		"table2": buildSampleRequestItems(5),
		"table3": buildSampleRequestItems(37),
		"table4": buildSampleRequestItems(11),
	}

	result := splitRequestItems(requestItems)
	c.Len(result, 4, fmt.Sprint("expected 4 batches, got ", len(result)))

	itemsPerBatch := 0
	for i, batch := range result {
		itemsPerBatch = 0

		for _, items := range batch {
			itemsPerBatch += len(items)
		}

		if i == 3 {
			//The last batch will have 5 items
			c.Equal(5, itemsPerBatch, fmt.Sprintf("Expected %d items per batch, but got %d", 5, itemsPerBatch))
		} else {
			c.Equal(dynamoDBMaxBatchWrite, itemsPerBatch, fmt.Sprintf("Expected %d items per batch, but got %d", dynamoDBMaxBatchWrite, itemsPerBatch))
		}

	}
}

func buildSampleRequestItems(length int) []types.WriteRequest {
	requestItems := make([]types.WriteRequest, 0, length)
	for i := 0; i < length; i++ {
		requestItems = append(requestItems, types.WriteRequest{
			PutRequest: &types.PutRequest{
				Item: map[string]types.AttributeValue{
					"key": &types.AttributeValueMemberS{Value: "value"},
				},
			},
		})
	}

	return requestItems
}
