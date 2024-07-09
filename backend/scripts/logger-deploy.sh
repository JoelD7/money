#!/bin/bash
set -o pipefail
echo "Deploying authenticator"
bash authenticator-deploy.sh
echo "Deploying lambda-authorizer"
bash lambda-authorizer-deploy.sh
echo "Deploying users"
bash users-deploy.sh
echo "Deploying expenses"
bash expenses-deploy.sh
echo "Deploying income"
bash income-deploy.sh
echo "Deploying recurrent-expense-period-setter"
bash recurrent-expense-period-setter-deploy.sh