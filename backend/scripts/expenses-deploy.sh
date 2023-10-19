#!/bin/bash
set -o pipefail
GOOS=linux GOARCH=arm64 go build -tags lambda.norpc -o api/deploy/bin/expenses/bootstrap github.com/JoelD7/money/backend/api/functions/expenses
zip -j api/deploy/bin/expenses/bootstrap.zip api/deploy/bin/expenses/bootstrap
aws lambda update-function-code --function-name money-expenses-handler --zip-file fileb://api/deploy/bin/expenses/bootstrap.zip | tee