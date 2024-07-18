#!/bin/bash
set -o pipefail
echo "Deploying recurrent-expense-period-setter"
GOOS=linux GOARCH=arm64 go build -tags lambda.norpc -o lambda/bin/recurrent-expense-period-setter/bootstrap github.com/JoelD7/money/backend/lambda/recurrent-expense-period-setter
zip -j lambda/bin/recurrent-expense-period-setter/bootstrap.zip lambda/bin/recurrent-expense-period-setter/bootstrap
aws lambda update-function-code --function-name money-recurrent-expense-period-setter --zip-file fileb://lambda/bin/recurrent-expense-period-setter/bootstrap.zip | tee