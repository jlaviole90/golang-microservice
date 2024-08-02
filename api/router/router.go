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

        r.Route("/health", func(r chi.Router) {
            healthApi := health.New(l, v)
            r.Method(
                http.MethodGet, 
                "/health",
                requestlog.NewHandler(healthApi.Health, l),
            )
            r.Method(
                http.MethodPost, 
                "/health", 
                requestlog.NewHandler(healthApi.BodyRequest, l),
            )
            r.Method(
                http.MethodGet,
                "/health/{number}",
                requestlog.NewHandler(healthApi.PathVariableRequest, l),
            )
        })
    })

    return r
}
