package router

import (
	"github.com/aws/aws-lambda-go/events"
)

type Router struct {
	path           string
	parent         *Router
	subRouter      *Router
	methodHandlers map[string]func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
}

func NewRouter(path string) *Router {
	return &Router{
		methodHandlers: map[string]func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error){},
	}
}

func (router *Router) Route(path string, fn func(r Router)) {

}

func (router *Router) Get(path string) {
	// The path is the same
	if path == "/" {

	}
}
