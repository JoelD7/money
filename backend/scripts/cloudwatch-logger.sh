#!/bin/bash
set -o pipefail
env GOOS=linux GOARCH=amd64 go build -o elk/deploy/bin/cloudwatch-logger/main  github.com/JoelD7/money/backend/elk/cloudwatch-logger
zip -j elk/deploy/bin/cloudwatch-logger/main.zip elk/deploy/bin/cloudwatch-logger/main
aws lambda update-function-code --function-name money-cloudwatch-logger --zip-file fileb://elk/deploy/bin/cloudwatch-logger/main.zip