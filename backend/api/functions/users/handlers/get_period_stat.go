package handlers

import (
	"context"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/storage/dynamo"
	"github.com/JoelD7/money/backend/storage/expenses"
	"github.com/JoelD7/money/backend/storage/income"
	"github.com/JoelD7/money/backend/usecases"
	"net/http"
	"sync"
	"time"
)

var gpstRequest *GetPeriodStatRequest
var gpstOnce sync.Once

type GetPeriodStatRequest struct {
	Log          logger.LogAPI
	startingTime time.Time
	err          error
	ExpensesRepo expenses.Repository
	IncomeRepo   income.Repository
}

func (request *GetPeriodStatRequest) init(ctx context.Context, log logger.LogAPI, envConfig *models.EnvironmentConfiguration) error {
	var err error

	gpstOnce.Do(func() {
		request.Log = log
		request.Log.SetHandler("get-period-stat")
		dynamoClient := dynamo.InitClient(ctx)

		request.IncomeRepo, err = income.NewDynamoRepository(dynamoClient, envConfig)
		if err != nil {
			return
		}

		request.ExpensesRepo, err = expenses.NewDynamoRepository(dynamoClient, envConfig)
		if err != nil {
			return
		}
	})

	request.startingTime = time.Now()

	return err
}

func (request *GetPeriodStatRequest) finish() {
	request.Log.LogLambdaTime(request.startingTime, request.err, recover())
}

func GetPeriodStatHandler(ctx context.Context, log logger.LogAPI, envConfig *models.EnvironmentConfiguration, req *apigateway.Request) (*apigateway.Response, error) {
	if gpstRequest == nil {
		gpstRequest = new(GetPeriodStatRequest)
	}

	err := gpstRequest.init(ctx, log, envConfig)
	if err != nil {
		gpstRequest.err = err

		log.Error("get_period_stat_init_failed", err, []models.LoggerObject{req})

		return req.NewErrorResponse(err), nil
	}

	return gpstRequest.Process(ctx, req)
}

func (request *GetPeriodStatRequest) Process(ctx context.Context, req *apigateway.Request) (*apigateway.Response, error) {
	periodID, ok := req.PathParameters["periodID"]
	if !ok || periodID == "" {
		request.Log.Error("missing_period_id", nil, []models.LoggerObject{req})

		return req.NewErrorResponse(models.ErrMissingPeriodID), nil
	}

	username, err := apigateway.GetUsernameFromContext(req)
	if err != nil {
		request.err = err
		request.Log.Error("get_username_from_context_failed", err, []models.LoggerObject{req})

		return req.NewErrorResponse(err), nil
	}

	getPeriodStats := usecases.NewPeriodStatsGetter(request.ExpensesRepo, request.IncomeRepo)

	periodStats, err := getPeriodStats(ctx, username, periodID)
	if err != nil {
		request.Log.Error("get_period_stats_failed", err, []models.LoggerObject{req})

		return req.NewErrorResponse(err), nil
	}

	return req.NewJSONResponse(http.StatusOK, periodStats), nil
}
