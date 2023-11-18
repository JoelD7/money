package shared

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
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
