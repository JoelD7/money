package main

import (
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
	secretsManager *secrets.SecretManager
}

func jwksHandler(_ *apigateway.Request) (*apigateway.Response, error) {
	req := &requestJwksHandler{
		log: logger.NewLogger(),
	}

	req.initJwksHandler()
	defer req.finish()

	return req.processJWKS()
}

func (req *requestJwksHandler) initJwksHandler() {
	req.secretsManager = secrets.NewSecretManager()
	req.startingTime = time.Now()
	req.log = logger.NewLogger()
}

func (req *requestJwksHandler) finish() {
	req.log.LogLambdaTime(req.startingTime, req.err, recover())
}

func (req *requestJwksHandler) processJWKS() (*apigateway.Response, error) {
	jsonWebKeySet, err := usecases.GetJsonWebKeySet(req.secretsManager, req.log)
	if err != nil {
		return getErrorResponse(err)
	}

	jsonResponse, err := json.Marshal(jsonWebKeySet)
	if err != nil {
		req.err = err
		req.log.Error("jwks_response_marshall_failed", err, nil)

		return getErrorResponse(err)
	}

	return &apigateway.Response{
		StatusCode: http.StatusOK,
		Body:       string(jsonResponse),
	}, nil
}
