package router

import (
	"fmt"
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
		methodHandlers: map[string]map[string]Handler{
			http.MethodGet:    make(map[string]Handler),
			http.MethodPost:   make(map[string]Handler),
			http.MethodPut:    make(map[string]Handler),
			http.MethodDelete: make(map[string]Handler),
		},
	}
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

	return endpoint
}

func (router *Router) isRoot() bool {
	return router.root == nil
}
