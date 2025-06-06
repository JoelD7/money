package router

import (
	"context"
	"errors"
	"fmt"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/storage/shared"
	"net/http"
)

var (
	errRouterIsNotRoot = errors.New("router is not root")
	errPathNotDefined  = errors.New("this path does not have a handler")
)

//TODO: make envConfig a value, not a pointer

// Handler type, defines the function signature for an APIGateway lambda handler.
// Takes the following arguments:
//  1. context provided by the AWS Lambda runtime,
//  2. an environment configuration object, so that the lambda function can access the environment variables,
//  4. and the request object provided by APIGateway.
type Handler func(ctx context.Context, envConfig *models.EnvironmentConfiguration, request *apigateway.Request) (*apigateway.Response, error)

type Router struct {
	path           string
	parent         *Router
	root           *Router
	methodHandlers map[string]map[string]Handler
	envConfig      *models.EnvironmentConfiguration
}

func NewRouter(envConfig *models.EnvironmentConfiguration) *Router {
	return &Router{
		envConfig: envConfig,
		methodHandlers: map[string]map[string]Handler{
			http.MethodGet:     make(map[string]Handler),
			http.MethodHead:    make(map[string]Handler),
			http.MethodPost:    make(map[string]Handler),
			http.MethodPatch:   make(map[string]Handler),
			http.MethodPut:     make(map[string]Handler),
			http.MethodDelete:  make(map[string]Handler),
			http.MethodOptions: make(map[string]Handler),
		},
	}
}

func (router *Router) Handle(ctx context.Context, request *apigateway.Request) (res *apigateway.Response, err error) {
	stackTrace, ctxErr := shared.ExecuteLambda(ctx, func(ctx context.Context) {
		res, err = router.executeHandle(ctx, router.envConfig, request)
	})

	if ctxErr != nil {
		logger.Error("request_timeout", ctxErr, models.Any("stack", map[string]interface{}{
			"s_trace": stackTrace,
		}))

		res = &apigateway.Response{
			StatusCode: http.StatusInternalServerError,
		}

		err = nil
	}

	return
}

func (router *Router) executeHandle(ctx context.Context, envConfig *models.EnvironmentConfiguration, request *apigateway.Request) (*apigateway.Response, error) {
	if !router.isRoot() {
		logger.Error("router_handle_failed", errRouterIsNotRoot, nil)

		return &apigateway.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       http.StatusText(http.StatusInternalServerError),
		}, errRouterIsNotRoot
	}

	if _, ok := router.methodHandlers[request.HTTPMethod][request.Resource]; !ok {
		logger.Error("router_handle_failed", errPathNotDefined, models.Any("router_data", map[string]interface{}{
			"s_path": router.path,
		}), models.Any("request", map[string]interface{}{
			"s_method":   request.HTTPMethod,
			"s_resource": request.Resource,
		}), models.Any("router_data", map[string]interface{}{
			"s_path": router.path,
		}), models.Any("request", map[string]interface{}{
			"s_method":   request.HTTPMethod,
			"s_resource": request.Resource,
		}))

		return &apigateway.Response{
			StatusCode: http.StatusInternalServerError,
		}, nil
	}

	return router.methodHandlers[request.HTTPMethod][request.Resource](ctx, envConfig, request)
}

func (router *Router) Route(path string, fn func(r *Router)) {
	router.path = path
	subRouter := &Router{parent: router}

	// This means that this is the root router
	if router.root == nil {
		subRouter.root = router
	}

	if router.root != nil {
		subRouter.root = router.root
	}

	fn(subRouter)
}

func (router *Router) Get(pattern string, handler Handler) {
	router.handlerAssigner(pattern, http.MethodGet, handler)
}

func (router *Router) Head(pattern string, handler Handler) {
	router.handlerAssigner(pattern, http.MethodHead, handler)
}

func (router *Router) Post(pattern string, handler Handler) {
	router.handlerAssigner(pattern, http.MethodPost, handler)
}

func (router *Router) Patch(pattern string, handler Handler) {
	router.handlerAssigner(pattern, http.MethodPatch, handler)
}

func (router *Router) Put(pattern string, handler Handler) {
	router.handlerAssigner(pattern, http.MethodPut, handler)
}

func (router *Router) Delete(pattern string, handler Handler) {
	router.handlerAssigner(pattern, http.MethodDelete, handler)
}

func (router *Router) Options(pattern string, handler Handler) {
	router.handlerAssigner(pattern, http.MethodOptions, handler)
}

func (router *Router) handlerAssigner(pattern string, method string, handler Handler) {
	if router.isRoot() && pattern != "/" {
		panic(fmt.Sprintf("This router is a root router. The pattern of a root routers should be '/', but is '%s'", pattern))
	}

	endpoint := router.getEndpoint(pattern)

	if router.root.methodHandlers[method][endpoint] != nil {
		panic(fmt.Sprintf("The path '%s' already has a handler for method '%s'", endpoint, method))
	}

	router.root.methodHandlers[method][endpoint] = handler
}

func (router *Router) getEndpoint(pattern string) string {
	endpointParts := make([]string, 0)

	if pattern != "/" {
		endpointParts = append(endpointParts, pattern)
	}

	// Gather relative patterns from this router up to form the endpoint
	cur := router
	for {
		if cur.path != "/" {
			endpointParts = append(endpointParts, cur.path)
		}

		if cur.parent == nil {
			break
		}

		cur = cur.parent
	}

	var endpoint string
	for i := len(endpointParts) - 1; i >= 0; i-- {
		endpoint += endpointParts[i]
	}

	//If the endpoint parts are empty, then the endpoint is the root
	if endpoint == "" {
		return "/"
	}

	return endpoint
}

func (router *Router) isRoot() bool {
	return router.root == nil
}
