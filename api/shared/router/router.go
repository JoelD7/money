package router

import (
	"github.com/aws/aws-lambda-go/events"
	"net/http"
)

type Handler func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

type Router struct {
	path           string
	parent         *Router
	root           *Router
	methodHandlers map[string]map[string]Handler
}

func NewRouter() *Router {
	return &Router{
		methodHandlers: make(map[string]map[string]Handler),
	}
}

func (router *Router) Route(path string, fn func(r *Router)) {
	router.path = path
	subRouter := &Router{parent: router}

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

func (router *Router) Post(pattern string, handler Handler) {
	router.handlerAssigner(pattern, http.MethodPost, handler)
}

func (router *Router) Put(pattern string, handler Handler) {
	router.handlerAssigner(pattern, http.MethodPut, handler)
}

func (router *Router) Delete(pattern string, handler Handler) {
	router.handlerAssigner(pattern, http.MethodDelete, handler)
}

func (router *Router) handlerAssigner(pattern string, method string, handler Handler) {
	endpoint := router.getEndpoint(pattern)

	if router.root.methodHandlers[method][endpoint] == nil {
		router.root.methodHandlers[method] = make(map[string]Handler)
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

	return endpoint
}
