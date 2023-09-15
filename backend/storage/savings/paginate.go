package savings

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func encodeLastKey(lastKey map[string]types.AttributeValue) (string, error) {
	if len(lastKey) == 0 {
		return "", nil
	}

	data, err := json.Marshal(lastKey)
	if err != nil {
		return "", fmt.Errorf("encoding last key: %v", err)
	}

	encoded := base64.URLEncoding.EncodeToString(data)

	return encoded, nil
}

func decodeStartKey(startKey string) (map[string]types.AttributeValue, error) {
	decoded, err := base64.URLEncoding.DecodeString(startKey)
	if err != nil {
		return nil, fmt.Errorf("decoding last key: %v", err)
	}

	exclusiveStartKey := map[string]types.AttributeValue{}

	err = json.Unmarshal(decoded, &exclusiveStartKey)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling last key: %v", err)
	}

	return exclusiveStartKey, nil
}
