#!/bin/bash
set -o pipefail
env GOOS=linux GOARCH=arm64 go build -tags lambda.norpc -o auth/deploy/bin/authenticator/bootstrap  github.com/JoelD7/money/backend/auth/authenticator
zip -j auth/deploy/bin/authenticator/bootstrap.zip auth/deploy/bin/authenticator/bootstrap
aws lambda update-function-code --function-name money-authenticator --zip-file fileb://auth/deploy/bin/authenticator/bootstrap.zip | tee