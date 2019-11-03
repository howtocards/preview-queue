package api

import (
	"context"
	"facade/src/internal/api/generated/restapi"
	"facade/src/internal/api/generated/restapi/operations"
	"facade/src/internal/queue"
	"fmt"
	"net/http"
	"path"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
	"github.com/powerman/structlog"
	"github.com/sebest/xff"
)

type (
	Queue interface {
		Send(event queue.Event) error
	}

	service struct {
		queue Queue
	}

	Config struct {
		Host     string
		Port     int
		BasePath string
	}
)

// NewServer returns server configured to listen on the TCP network
func NewServer(log *structlog.Logger, queue Queue, cfg Config) (*restapi.Server, error) {
	svc := &service{
		queue: queue,
	}

	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		return nil, fmt.Errorf("loads embedded: %w", err)
	}
	if cfg.BasePath == "" {
		cfg.BasePath = swaggerSpec.BasePath()
	}
	swaggerSpec.Spec().BasePath = cfg.BasePath

	api := operations.NewFacadeAPI(swaggerSpec)
	api.Logger = log.New(structlog.KeyUnit, "swagger").Printf

	api.RenderHandler = operations.RenderHandlerFunc(svc.render)

	server := restapi.NewServer(api)
	server.Host = cfg.Host
	server.Port = cfg.Port

	// The middleware executes before anything.
	globalMiddlewares := func(handler http.Handler) http.Handler {
		xffmw, _ := xff.Default()
		logger := makeLogger(cfg.BasePath)
		accesslog := makeAccessLog(cfg.BasePath)
		redocOpts := middleware.RedocOpts{
			BasePath: cfg.BasePath,
			SpecURL:  path.Join(cfg.BasePath, "/swagger.json"),
		}
		return xffmw.Handler(logger(recovery(accesslog(
			middleware.Spec(cfg.BasePath, restapi.FlatSwaggerJSON,
				middleware.Redoc(redocOpts,
					handler))))))
	}

	server.SetHandler(globalMiddlewares(api.Serve(nil)))

	log.Info("protocol", "version", swaggerSpec.Spec().Info.Version)
	return server, nil
}

func fromRequest(r *http.Request) (context.Context, *structlog.Logger) {
	return r.Context(), structlog.FromContext(r.Context(), nil)
}
