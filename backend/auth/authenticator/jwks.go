package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/shared/secrets"
	"github.com/JoelD7/money/backend/usecases"
)

type requestJwksHandler struct {
	log            logger.LogAPI
	startingTime   time.Time
	err            error
	secretsManager secrets.SecretManager
}

func jwksHandler(ctx context.Context, log logger.LogAPI, request *apigateway.Request) (*apigateway.Response, error) {
	req := &requestJwksHandler{}

	req.initJwksHandler(log)
	defer req.finish()

	return req.processJWKS(ctx, request)
}

func (req *requestJwksHandler) initJwksHandler(log logger.LogAPI) {
	req.secretsManager = secrets.NewAWSSecretManager()
	req.startingTime = time.Now()
	req.log = logger.NewLoggerWithHandler("jwks")
}

func (req *requestJwksHandler) finish() {
	req.log.LogLambdaTime(req.startingTime, req.err, recover())
}

func (req *requestJwksHandler) processJWKS(ctx context.Context, request *apigateway.Request) (*apigateway.Response, error) {
	jsonWebKeySet, err := usecases.GetJsonWebKeySet(ctx, req.secretsManager, req.log)
	if err != nil {
		return request.NewErrorResponse(err), nil
	}

	jsonResponse, err := json.Marshal(jsonWebKeySet)
	if err != nil {
		req.err = err
		req.log.Error("jwks_response_marshall_failed", err, nil)

		return request.NewErrorResponse(err), nil
	}

	return &apigateway.Response{
		StatusCode: http.StatusOK,
		Body:       string(jsonResponse),
	}, nil
}
