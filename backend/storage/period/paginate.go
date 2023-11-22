package period

type keys struct {
	Period   string `json:"period" dynamodbav:"period"`
	Username string `json:"username" dynamodbav:"username"`
}
