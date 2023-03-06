#!/bin/bash
set -o pipefail
env GOOS=linux GOARCH=amd64 go build -o auth/deploy/bin/authenticator/main  github.com/JoelD7/money/backend/auth/authenticator
zip -j auth/deploy/bin/authenticator/main.zip auth/deploy/bin/authenticator/main
aws lambda update-function-code --function-name money-authenticator --zip-file fileb://auth/deploy/bin/authenticator/main.zip