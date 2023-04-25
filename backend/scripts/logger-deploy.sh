#!/bin/bash
set -o pipefail
echo "Deploying authenticator"
bash authenticator-deploy.sh
echo "Deploying lambda-authorizer"
bash lambda-authorizer-deploy.sh
echo echo "Deploying users"
bash users-deploy.sh