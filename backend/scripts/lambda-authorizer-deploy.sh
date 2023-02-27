#!/bin/bash
set -o pipefail
env GOOS=linux GOARCH=amd64 go build -o api/deploy/bin/lambda-authorizer/main  github.com/JoelD7/money/backend/auth/lambda-authorizer
zip -j api/deploy/bin/lambda-authorizer/main.zip api/deploy/bin/lambda-authorizer/main
aws lambda update-function-code --function-name money-lambda-authorizer --zip-file fileb://api/deploy/bin/lambda-authorizer/main.zip