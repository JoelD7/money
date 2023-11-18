package savings

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
