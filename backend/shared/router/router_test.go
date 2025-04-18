package router

import (
	"context"
	"fmt"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

var endpoint = "/users/categories/number"

func TestHandle(t *testing.T) {
	c := require.New(t)

	rootRouter := NewRouter(&models.EnvironmentConfiguration{})
	ctx := context.Background()

	rootRouter.Route("/", func(r *Router) {
		r.Options("/", dummyHandler())
		r.Route("/users", func(r *Router) {
			r.Route("/{userID}", func(r *Router) {
				r.Get("/", dummyHandler())
				r.Post("/", dummyHandler())
				r.Put("/", dummyHandler())
				r.Delete("/", dummyHandler())

				r.Get("/categories", dummyHandler())
			})
		})

		r.Route("/savings", func(r *Router) {
			r.Get("/{savingID}", dummyHandler())
			r.Get("/", dummyHandler())
			r.Post("/", dummyHandler())
			r.Put("/{savingID}", dummyHandler())
			r.Delete("/", dummyHandler())
		})
	})

	request := &apigateway.Request{
		HTTPMethod: http.MethodGet,
		Resource:   "/users/{userID}",
	}

	response, err := rootRouter.Handle(ctx, request)
	c.Nil(err)
	c.Equal("Method: GET, Endpoint: /users/{userID}", response.Body)

	request.Resource = "/users/{userID}/categories"
	response, err = rootRouter.Handle(ctx, request)
	c.Nil(err)
	c.Equal("Method: GET, Endpoint: /users/{userID}/categories", response.Body)

	request.Resource = "/users/{userID}"
	request.HTTPMethod = http.MethodPost

	response, err = rootRouter.Handle(ctx, request)
	c.Nil(err)
	c.Equal("Method: POST, Endpoint: /users/{userID}", response.Body)

	request = &apigateway.Request{
		HTTPMethod: http.MethodOptions,
		Resource:   "/",
	}

	response, err = rootRouter.Handle(ctx, request)
	c.Nil(err)
	c.Equal(http.StatusOK, response.StatusCode)
}

func TestHandleError(t *testing.T) {
	c := require.New(t)

	rootRouter := NewRouter(&models.EnvironmentConfiguration{})
	ctx := context.Background()

	subRouter := &Router{
		path:           "/users/{userID}",
		parent:         rootRouter,
		root:           rootRouter,
		methodHandlers: nil,
	}

	request := &apigateway.Request{
		HTTPMethod: http.MethodGet,
		Resource:   "/users/{userID}",
	}

	response, err := subRouter.Handle(ctx, request)
	c.Equal(http.StatusInternalServerError, response.StatusCode)
	c.Equal(http.StatusText(http.StatusInternalServerError), response.Body)
	c.ErrorIs(err, errRouterIsNotRoot)

	rootRouter = NewRouter(&models.EnvironmentConfiguration{})

	rootRouter.Route("/users", func(r *Router) {
		r.Post("/dummy", dummyHandler())
	})

	response, err = rootRouter.Handle(ctx, request)
	c.Equal(http.StatusInternalServerError, response.StatusCode)
	c.ErrorIs(err, errPathNotDefined)
}

func TestRoute(t *testing.T) {
	c := require.New(t)

	rootRouter := NewRouter(&models.EnvironmentConfiguration{})

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

	rootRouter = NewRouter(&models.EnvironmentConfiguration{})
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

		r.Route("/savings", func(r *Router) {
			r.Get("/", dummyHandler())
		})
	})

	_, ok = rootRouter.methodHandlers[http.MethodGet]["/savings"]
	c.True(ok)

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
		rootRouter = NewRouter(&models.EnvironmentConfiguration{})
		rootRouter.Get(endpoint, dummyHandler())
	})
}

func dummyRouter() *Router {
	rootRouter := NewRouter(&models.EnvironmentConfiguration{})

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
	return func(ctx context.Context, envConfig *models.EnvironmentConfiguration, request *apigateway.Request) (*apigateway.Response, error) {
		return &apigateway.Response{
			Body: fmt.Sprintf("Method: %s, Endpoint: %s", request.HTTPMethod, request.Resource),
		}, nil
	}
}
