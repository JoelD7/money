#!/bin/bash
set -o pipefail

bash authenticator-deploy.sh &
bash lambda-authorizer-deploy.sh &
bash users-deploy.sh &
bash expenses-deploy.sh &
bash income-deploy.sh &
bash recurrent-expense-period-setter-deploy.sh
