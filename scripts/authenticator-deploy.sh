#!/bin/bash
set -o pipefail
env GOOS=linux GOARCH=amd64 go build -o api/deploy/bin/authenticator/main  github.com/JoelD7/money/auth/authenticator
zip -j api/deploy/bin/authenticator/main.zip api/deploy/bin/authenticator/main
aws lambda update-function-code --function-name money-authenticator --zip-file fileb://api/deploy/bin/authenticator/main.zip