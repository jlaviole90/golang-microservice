package router

import (
	"employee-worklog-service/api/resources/health"
	"employee-worklog-service/api/router/middleware"
	"employee-worklog-service/api/router/middleware/requestlog"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
)

func New(l *zerolog.Logger, v *validator.Validate) *chi.Mux {
    r := chi.NewRouter()

    r.Use(cors.Handler(cors.Options{
        AllowedOrigins: []string{"*"},
        AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    }))

    r.Route("/api/v1", func(r chi.Router) {
        r.Use(middleware.RequestId)
        r.Use(middleware.ContentTypeJSON)

        r.Get("/health", health.Read)

        r.Route("/test", func(r chi.Router) {
            testApi := test.New(l, v)
            r.Method(
                http.MethodPost,
                "/{number}",
                requestlog.NewHandler(testApi.TestEndpoint, l),
            ),
        })
    })

    return r
}
