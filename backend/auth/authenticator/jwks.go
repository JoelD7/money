package main

import (
	"context"
	"encoding/json"
	"github.com/JoelD7/money/backend/models"
	"net/http"
	"sync"
	"time"

	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/shared/secrets"
	"github.com/JoelD7/money/backend/usecases"
)

var jwksRequest *requestJwksHandler
var jwksOnce sync.Once

type requestJwksHandler struct {
	startingTime   time.Time
	err            error
	secretsManager secrets.SecretManager
}

func jwksHandler(ctx context.Context, _ *models.EnvironmentConfiguration, request *apigateway.Request) (*apigateway.Response, error) {
	if jwksRequest == nil {
		jwksRequest = new(requestJwksHandler)
	}

	jwksRequest.initJwksHandler()
	defer jwksRequest.finish()

	return jwksRequest.processJWKS(ctx, request)
}

func (req *requestJwksHandler) initJwksHandler() {
	jwksOnce.Do(func() {
		req.secretsManager = secrets.NewAWSSecretManager()
		logger.SetHandler("jwks")
	})
	req.startingTime = time.Now()
}

func (req *requestJwksHandler) finish() {
	logger.LogLambdaTime(req.startingTime, req.err, recover())
}

func (req *requestJwksHandler) processJWKS(ctx context.Context, request *apigateway.Request) (*apigateway.Response, error) {
	jsonWebKeySet, err := usecases.GetJsonWebKeySet(ctx, req.secretsManager)
	if err != nil {
		return request.NewErrorResponse(err), nil
	}

	jsonResponse, err := json.Marshal(jsonWebKeySet)
	if err != nil {
		req.err = err
		logger.Error("jwks_response_marshall_failed", err, nil)

		return request.NewErrorResponse(err), nil
	}

	return &apigateway.Response{
		StatusCode: http.StatusOK,
		Body:       string(jsonResponse),
	}, nil
}
