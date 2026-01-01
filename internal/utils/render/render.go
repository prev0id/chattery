package render

import (
	"chattery/internal/utils/errors"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

func JSON(w http.ResponseWriter, r *http.Request, value any) {
	if value == nil {
		return
	}

	if err, ok := value.(error); ok {
		Error(w, r, err)
		return
	}

	response, err := json.Marshal(value)
	if err != nil {
		Error(w, r, errors.E(err).Debug("json.Marshal"))
		return
	}

	setContentTypeJSON(w)
	w.Write(response)
}

type responseError struct {
	Message string `json:"message"`
}

func Error(w http.ResponseWriter, r *http.Request, err error) {
	requestID := middleware.GetReqID(r.Context())

	domainErr := errors.E(err).Log("request ended with an error", slog.String("request_id", requestID))

	response, _ := json.Marshal(responseError{Message: domainErr.Error()})

	setContentTypeJSON(w)
	w.Write(response)
}

func setContentTypeJSON(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}
