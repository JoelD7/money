#!/bin/bash
set -o pipefail
echo "Deploying income"
GOOS=linux GOARCH=arm64 go build -tags lambda.norpc -o api/deploy/bin/income/bootstrap github.com/JoelD7/money/backend/api/functions/income
zip -j api/deploy/bin/income/bootstrap.zip api/deploy/bin/income/bootstrap
aws lambda update-function-code --function-name money-income-handler --zip-file fileb://api/deploy/bin/income/bootstrap.zip | tee