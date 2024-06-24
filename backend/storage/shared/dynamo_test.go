package shared

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"testing"
)

func TestSplitRequestItems(t *testing.T) {
	//c := require.New(t)

	requestItems := map[string][]types.WriteRequest{
		"table1": buildSampleRequestItems(27),
		"table2": buildSampleRequestItems(5),
		"table3": buildSampleRequestItems(37),
		"table4": buildSampleRequestItems(11),
	}

	result := splitRequestItems(requestItems)
	for _, batch := range result {
		for table, items := range batch {
			fmt.Println(table, len(items))
		}
	}
	//c.Len(result, 3, fmt.Sprint("expected 4 batches, got ", len(result)))
	//c.Len(result[0], 1, fmt.Sprintf("First batch should have 1 table, got %d", len(result[0])))
	//c.Len(result[1], 3, fmt.Sprintf("Second batch should have 3 tables, got %d", len(result[1])))
	//c.Len(result[2], 2, fmt.Sprintf("Third batch should have 2 tables, got %d", len(result[2])))

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
