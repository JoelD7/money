package router

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
)

type Handler func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

type Router struct {
	path           string
	parent         *Router
	subRouter      *Router
	methodHandlers map[string]Handler
}

func NewRouter() *Router {
	return &Router{
		methodHandlers: make(map[string]Handler),
	}
}

func (router *Router) Route(path string, fn func(r *Router)) {
	//router.subRouter = subRouter
	router.path = path
	subRouter := &Router{parent: router}

	fn(subRouter)
}

func (router *Router) Get(path string, handler Handler) {
	endpointParts := make([]string, 0)

	if path != "/" {
		endpointParts = append(endpointParts, path)
	}

	// Search root parent in order to build full path
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

	var fullEndpoint string
	for i := len(endpointParts) - 1; i >= 0; i-- {
		fullEndpoint += endpointParts[i]
	}

	fmt.Println("Full endpoint: ", fullEndpoint)
}
