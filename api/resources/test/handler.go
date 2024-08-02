package test

import (
	"net/http"

    e "employee-worklog-service/api/resources/common/err"
	"employee-worklog-service/utils/ctx"
	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
)

type TestApi struct {
	logger    *zerolog.Logger
	validator *validator.Validate
}

func New(logger *zerolog.Logger, validator *validator.Validate) *TestApi {
	return &TestApi{
		logger:    logger,
		validator: validator,
	}
}

func (t *TestApi) TestEndpoint(w http.ResponseWriter, r *http.Request) {
	reqId := ctx.RequestId(r.Context())

	number := chi.URLParam(r, "number")
}
