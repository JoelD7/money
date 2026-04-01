#!/bin/bash
set -o pipefail

bash scripts/authenticator-deploy.sh &
bash scripts/lambda-authorizer-deploy.sh &
bash scripts/users-deploy.sh &
bash scripts/expenses-deploy.sh &
bash scripts/income-deploy.sh &
bash scripts/recurrent-expense-period-setter-deploy.sh
