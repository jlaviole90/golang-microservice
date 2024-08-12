package health

import (
	"encoding/json"
	"net/http"

	e "employee-worklog-service/api"
	exampleservice "employee-worklog-service/services/example_service"
	"employee-worklog-service/utils/ctx"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
)

type HealthApi struct {
	logger    *zerolog.Logger
	validator *validator.Validate
}

func New(logger *zerolog.Logger, validator *validator.Validate) *HealthApi {
	return &HealthApi{
		logger:    logger,
		validator: validator,
	}
}

func (h *HealthApi) Health(w http.ResponseWriter, r *http.Request) {
    reqId := ctx.RequestId(r.Context())

    j := exampleservice.GetExample()

    if err := json.NewEncoder(w).Encode(j); err != nil {
        h.logger.Error().Str("KeyReqId", reqId).Err(err).Msg("Failed to write response")
        e.IOError(w, e.IOErrorResponse)
        return
	}
}

func (h *HealthApi) PathVariableRequest(w http.ResponseWriter, r *http.Request) {
	reqId := ctx.RequestId(r.Context())

	pathVar := chi.URLParam(r, "number")

	if is, err := validateRequest(pathVar, h.validator); !is {
		h.logger.Error().Str("KeyReqId", reqId).Msg(err)
		e.BadRequest(w, e.BadRequestResponse)
		return
	}

    if err := json.NewEncoder(w).Encode(pathVar); err != nil {
		h.logger.Error().Str("KeyReqId", reqId).Err(err).Msg("Failed to write response")
		e.IOError(w, e.IOErrorResponse)
		return
	}
}

func (h *HealthApi) BodyRequest(w http.ResponseWriter, r *http.Request) {
	reqId := ctx.RequestId(r.Context())

	var body string
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		h.logger.Error().Str("KeyReqId", reqId).Err(err).Msg("Failed to decode request body")
		e.JSONDeserializeError(w, e.JSONDeserializeFailureResponse)
		return
	}

	if is, err := validateRequest(body, h.validator); !is {
		h.logger.Error().Str("KeyReqId", reqId).Msg(err)
		e.BadRequest(w, e.BadRequestResponse)
		return
	}

    if err := json.NewEncoder(w).Encode(body); err != nil {
		h.logger.Error().Str("KeyReqId", reqId).Err(err).Msg("Failed to write response")
		e.IOError(w, e.IOErrorResponse)
		return
	}
}

func validateRequest(req string, v *validator.Validate) (bool, string) {
	if err := v.Var(req, "required"); err != nil {
		return false, "Request body is required"
	}
	if err := v.Var(req, "alpha_space"); err != nil {
		return false, "Request body must contain only letters and spaces"
	}

	return true, ""
}
