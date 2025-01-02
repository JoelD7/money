package dynamo

import (
	"fmt"
	"github.com/JoelD7/money/backend/shared/env"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSplitRequestItems(t *testing.T) {
	c := require.New(t)

	dynamoDBMaxBatchWrite := env.GetInt("DYNAMODB_MAX_BATCH_WRITE", 25)

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

func TestBuildAmountKey(t *testing.T) {
	cases := []struct {
		name        string
		amount      float64
		id          string
		expectedKey string
	}{
		{
			name:        "Positive amount",
			amount:      1234.56,
			id:          "abc123",
			expectedKey: "01234.560000:abc123",
		},
		{
			name:        "Negative amount",
			amount:      -987.65,
			id:          "xyz789",
			expectedKey: "-987.650000:xyz789",
		},
		{
			name:        "Zero amount",
			amount:      0,
			id:          "zero123",
			expectedKey: "00000.000000:zero123",
		},
		{
			name:        "Large amount",
			amount:      99999999.99,
			id:          "big123",
			expectedKey: "99999999.990000:big123",
		},
		{
			name:        "Small fractional amount",
			amount:      0.000001,
			id:          "tiny123",
			expectedKey: "00000.000001:tiny123",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := BuildAmountKey(tc.amount, tc.id)
			assert.NotNil(t, result)
			assert.Equal(t, tc.expectedKey, result)
		})
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
