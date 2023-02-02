package router

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGet(t *testing.T) {
	c := require.New(t)

	rootRouter := NewRouter()
	rootRouter.Route("/", func(r *Router) {
		r.Route("/users", func(r *Router) {
			r.Route("/categories", func(r *Router) {
				r.Get("/number", dummyHandler())
			})
		})
	})
	usersRouter := &Router{path: "/users", parent: rootRouter}
	categoriesRouter := &Router{path: "/users/categories", parent: usersRouter}

	c.Equal("/users/categories", categoriesRouter.path)
}

func dummyHandler() Handler {
	return func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		return events.APIGatewayProxyResponse{}, nil
	}
}
