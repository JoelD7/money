package savings

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type keys struct {
	SavingID string `json:"saving_id" dynamodbav:"saving_id"`
	Username string `json:"username" dynamodbav:"username"`
}

type keysPeriodIndex struct {
	SavingID   string `json:"saving_id" dynamodbav:"saving_id"`
	Username   string `json:"username" dynamodbav:"username"`
	PeriodUser string `json:"period_user" dynamodbav:"period_user"`
}

type keysSavingGoalIndex struct {
	SavingID     string `json:"saving_id" dynamodbav:"saving_id"`
	Username     string `json:"username" dynamodbav:"username"`
	SavingGoalID string `json:"saving_goal_id" dynamodbav:"saving_goal_id"`
}

// encodeLastKey encodes the last evaluated key returned by Dynamo in a string format to be used in the next query
// as the start key.
// The "keyType" parameter should be a pointer to a struct that maps ot the primary key of the table or index in question.
func encodeLastKey(lastKey map[string]types.AttributeValue, keyType interface{}) (string, error) {
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

// decodeStartKey parses the start key string into a map of attribute values to be used as ExclusiveStartKey in a paginated
// query.
// The "keyType" parameter should be a pointer to a struct that maps ot the primary key of the table or index in question.
func decodeStartKey(startKey string, keyType interface{}) (map[string]types.AttributeValue, error) {
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
