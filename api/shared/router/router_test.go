package router

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

var endpoint = "/users/categories/number"

func TestHandle(t *testing.T) {
	c := require.New(t)

	rootRouter := NewRouter()

	rootRouter.Route("/", func(r *Router) {
		r.Route("/users", func(r *Router) {
			r.Route("/{userID}", func(r *Router) {
				r.Get("/", dummyHandler())
				r.Post("/", dummyHandler())
				r.Put("/", dummyHandler())
				r.Delete("/", dummyHandler())

				r.Get("/categories", dummyHandler())
			})
		})
	})

	request := &events.APIGatewayProxyRequest{
		HTTPMethod: http.MethodGet,
		Resource:   "/users/{userID}",
	}

	response, err := rootRouter.Handle(request)
	c.Nil(err)
	c.Equal("Method: GET, Endpoint: /users/{userID}", response.Body)

	request.Resource = "/users/{userID}/categories"
	response, err = rootRouter.Handle(request)
	c.Nil(err)
	c.Equal("Method: GET, Endpoint: /users/{userID}/categories", response.Body)

	request.Resource = "/users/{userID}"
	request.HTTPMethod = http.MethodPost

	response, err = rootRouter.Handle(request)
	c.Nil(err)
	c.Equal("Method: POST, Endpoint: /users/{userID}", response.Body)
}

func TestHandleError(t *testing.T) {
	c := require.New(t)

	rootRouter := NewRouter()
	subRouter := &Router{
		path:           "/users/{userID}",
		parent:         rootRouter,
		root:           rootRouter,
		methodHandlers: nil,
	}

	request := &events.APIGatewayProxyRequest{
		HTTPMethod: http.MethodGet,
		Resource:   "/users/{userID}",
	}

	response, err := subRouter.Handle(request)
	c.Equal(http.StatusInternalServerError, response.StatusCode)
	c.Equal(http.StatusText(http.StatusInternalServerError), response.Body)
	c.ErrorIs(err, errRouterIsNotRoot)

	rootRouter = NewRouter()

	rootRouter.Route("/users", func(r *Router) {
		r.Post("/dummy", dummyHandler())
	})

	response, err = rootRouter.Handle(request)
	c.Equal(http.StatusInternalServerError, response.StatusCode)
	c.ErrorIs(err, errPathNotDefined)
}

func TestRoute(t *testing.T) {
	c := require.New(t)

	rootRouter := NewRouter()

	rootRouter.Route("/", func(r *Router) {
		r.Route("/users", func(r *Router) {
			r.Route("/categories", func(r *Router) {
				r.Get("/number", dummyHandler())
				r.Post("/number", dummyHandler())
				r.Put("/number", dummyHandler())
				r.Delete("/number", dummyHandler())
			})
		})
	})

	_, ok := rootRouter.methodHandlers[http.MethodGet][endpoint]
	c.True(ok)

	_, ok = rootRouter.methodHandlers[http.MethodPost][endpoint]
	c.True(ok)

	_, ok = rootRouter.methodHandlers[http.MethodPut][endpoint]
	c.True(ok)

	_, ok = rootRouter.methodHandlers[http.MethodDelete][endpoint]
	c.True(ok)

	rootRouter = NewRouter()
	rootRouter.Route("/", func(r *Router) {
		r.Route("/users", func(r *Router) {
			r.Route("/{userID}", func(r *Router) {
				r.Get("/", dummyHandler())
				r.Post("/", dummyHandler())
				r.Put("/", dummyHandler())
				r.Delete("/", dummyHandler())

				r.Get("/categories", dummyHandler())
			})
		})
	})

	_, ok = rootRouter.methodHandlers[http.MethodGet]["/users/{userID}"]
	c.True(ok)

	_, ok = rootRouter.methodHandlers[http.MethodPost]["/users/{userID}"]
	c.True(ok)

	_, ok = rootRouter.methodHandlers[http.MethodPut]["/users/{userID}"]
	c.True(ok)

	_, ok = rootRouter.methodHandlers[http.MethodDelete]["/users/{userID}"]
	c.True(ok)

	_, ok = rootRouter.methodHandlers[http.MethodGet]["/users/{userID}/categories"]
	c.True(ok)
}

func TestPanics(t *testing.T) {
	c := require.New(t)

	rootRouter := dummyRouter()

	c.PanicsWithValue(fmt.Sprintf("The path '%s' already has a handler for method '%s'", endpoint, http.MethodGet), func() {
		rootRouter.Route("/", func(r *Router) {
			r.Route("/users", func(r *Router) {
				r.Route("/categories", func(r *Router) {
					r.Get("/number", dummyHandler())
				})
			})
		})
	})

	c.PanicsWithValue(fmt.Sprintf("This router is a root router. The pattern of a root routers should be '/', but is '%s'", endpoint), func() {
		rootRouter = NewRouter()
		rootRouter.Get(endpoint, dummyHandler())
	})
}

func dummyRouter() *Router {
	rootRouter := NewRouter()

	rootRouter.Route("/", func(r *Router) {
		r.Route("/users", func(r *Router) {
			r.Route("/categories", func(r *Router) {
				r.Get("/number", dummyHandler())
			})
		})
	})

	return rootRouter
}

func dummyHandler() Handler {
	return func(request *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
		return &events.APIGatewayProxyResponse{
			Body: fmt.Sprintf("Method: %s, Endpoint: %s", request.HTTPMethod, request.Resource),
		}, nil
	}
}
