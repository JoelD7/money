#!/bin/bash
set -o pipefail
env GOOS=linux GOARCH=arm64 go build -tags lambda.norpc -o api/deploy/bin/bootstrap  github.com/JoelD7/money/backend/api/functions/users
zip -j api/deploy/bin/bootstrap.zip api/deploy/bin/bootstrap
aws lambda update-function-code --function-name money-users-handler --zip-file fileb://api/deploy/bin/bootstrap.zip