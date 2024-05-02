package main

import (
	"context"
	"encoding/json"
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
	log            logger.LogAPI
	startingTime   time.Time
	err            error
	secretsManager secrets.SecretManager
}

func jwksHandler(ctx context.Context, log logger.LogAPI, request *apigateway.Request) (*apigateway.Response, error) {
	if jwksRequest == nil {
		jwksRequest = new(requestJwksHandler)
	}

	jwksRequest.initJwksHandler(log)
	defer jwksRequest.finish()

	return jwksRequest.processJWKS(ctx, request)
}

func (req *requestJwksHandler) initJwksHandler(log logger.LogAPI) {
	jwksOnce.Do(func() {
		req.secretsManager = secrets.NewAWSSecretManager()
		req.log = log
		req.log.SetHandler("jwks")
	})
	req.startingTime = time.Now()
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
