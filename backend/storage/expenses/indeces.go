package expenses

/*
We get a validation error from Dynamo("The provided starting key is invalid") if the ExclusiveStartKey used in a
paginated query doesn't include the main table's primary key attributes: ExpenseID and Username. Note that even though
expense_id and username aren't part of neither the username-created_date-index nor the period_user-created_date-index
primary keys, they are included in the respective indeces' structs.
*/

type keys struct {
	ExpenseID string `json:"expense_id" dynamodbav:"expense_id"`
	Username  string `json:"username" dynamodbav:"username"`
}

type keysPeriodUserIndex struct {
	ExpenseID  string `json:"expense_id" dynamodbav:"expense_id"`
	Username   string `json:"username,omitempty" dynamodbav:"username"`
	PeriodUser string `json:"period_user,omitempty" dynamodbav:"period_user"`
}

type keysUsernameCreatedDateIndex struct {
	Username    string `json:"username" dynamodbav:"username"`
	ExpenseID   string `json:"expense_id" dynamodbav:"expense_id"`
	CreatedDate string `json:"created_date" dynamodbav:"created_date"`
}

type keysPeriodUserCreatedDateIndex struct {
	ExpenseID   string `json:"expense_id" dynamodbav:"expense_id"`
	PeriodUser  string `json:"period_user,omitempty" dynamodbav:"period_user"`
	Username    string `json:"username" dynamodbav:"username"`
	CreatedDate string `json:"created_date" dynamodbav:"created_date"`
}
