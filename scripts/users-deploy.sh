#!/bin/bash
set -o pipefail
env GOOS=linux GOARCH=amd64 go build -o api/deploy/bin/main  github.com/JoelD7/money/api/users
zip -j api/deploy/bin/main.zip api/deploy/bin/main
aws lambda update-function-code --function-name money-users-handler --zip-file fileb://api/deploy/bin/main.zip