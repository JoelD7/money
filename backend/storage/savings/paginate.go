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

func encodeLastKey(lastKey map[string]types.AttributeValue) (string, error) {
	if len(lastKey) == 0 {
		return "", nil
	}

	primaryKey := new(keys)

	err := attributevalue.UnmarshalMap(lastKey, primaryKey)
	if err != nil {
		return "", fmt.Errorf("unmarshalling lastKey map: %v", err)
	}

	data, err := json.Marshal(primaryKey)
	if err != nil {
		return "", fmt.Errorf("marshalling primary key: %v", err)
	}

	encoded := base64.URLEncoding.EncodeToString(data)

	return encoded, nil
}

func decodeStartKey(startKey string) (map[string]types.AttributeValue, error) {
	decoded, err := base64.URLEncoding.DecodeString(startKey)
	if err != nil {
		return nil, fmt.Errorf("decoding last key: %v", err)
	}

	primaryKeyDecoded := new(keys)
	err = json.Unmarshal(decoded, primaryKeyDecoded)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling primary key: %v", err)
	}

	exclusiveStartKey, err := attributevalue.MarshalMap(primaryKeyDecoded)
	if err != nil {
		return nil, fmt.Errorf("marshalling to map of attribute value: %v", err)
	}

	return exclusiveStartKey, nil
}
