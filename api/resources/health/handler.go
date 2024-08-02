package health

import (
	"employee-worklog-service/utils/ctx"
	"encoding/json"
	"net/http"

	e "employee-worklog-service/api/resources/common/err"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
)

type HealthApi struct {
    logger *zerolog.Logger
    validator *validator.Validate
}

func New(logger *zerolog.Logger, validator *validator.Validate) *HealthApi {
    return &HealthApi{
        logger: logger,
        validator: validator,
    }
}

func (h *HealthApi) Health(w http.ResponseWriter, r *http.Request) {
	_, e := w.Write([]byte("."))
	if e != nil {
		panic(e)
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

    _, er := w.Write([]byte("Path variable: " + pathVar))
    if er != nil {
        h.logger.Error().Err(er).Msg("Failed to write response")
        e.IOError(w, e.IOErrorResponse)
        return
    }
}

func (h *HealthApi) BodyRequest(w http.ResponseWriter, r *http.Request) {
    reqId := ctx.RequestId(r.Context())

    var body string
    if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
        h.logger.Error().Err(err).Msg("Failed to decode request body")
        e.JSONDeserializeError(w, e.JSONDeserializeFailureResponse)
        return
    }

    if is, err := validateRequest(body, h.validator); !is {
        h.logger.Error().Str("KeyReqId", reqId).Msg(err)
        e.BadRequest(w, e.BadRequestResponse)
        return
    }

	_, er := w.Write([]byte("Hello world!"))
	if er != nil {
        h.logger.Error().Err(er).Msg("Failed to write response")
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
