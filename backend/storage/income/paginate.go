package income

type keys struct {
	IncomeID string `json:"income_id" dynamodbav:"income_id"`
	Username string `json:"username" dynamodbav:"username"`
}

type keysPeriodUserIndex struct {
	IncomeID   string `json:"income_id" dynamodbav:"income_id"`
	PeriodUser string `json:"period_user" dynamodbav:"period_user"`
	Username   string `json:"username" dynamodbav:"username"`
}
