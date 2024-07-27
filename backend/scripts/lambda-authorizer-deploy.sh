#!/bin/bash
set -o pipefail
echo "Deploying lambda-authorizer"
env GOOS=linux GOARCH=arm64 go build -tags lambda.norpc -o api/deploy/bin/lambda-authorizer/bootstrap  github.com/JoelD7/money/backend/auth/lambda-authorizer
zip -j api/deploy/bin/lambda-authorizer/bootstrap.zip api/deploy/bin/lambda-authorizer/bootstrap
aws lambda update-function-code --function-name money-lambda-authorizer --zip-file fileb://api/deploy/bin/lambda-authorizer/bootstrap.zip | tee