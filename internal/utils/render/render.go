package render

import (
	"chattery/internal/utils/errors"
	"chattery/internal/utils/logger"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

func Json(w http.ResponseWriter, r *http.Request, value any) {
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

	logger.Error(err, "request ended with an error", slog.String("request_id", requestID))

	response, _ := json.Marshal(responseError{Message: err.Error()})

	setContentTypeJSON(w)
	w.Write(response)
}

func setContentTypeJSON(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

func JsonBytes(value any) ([]byte, error) {
	response, err := json.Marshal(value)
	if err != nil {
		return nil, errors.E(err).Debug("json.Marshal")
	}
	return response, nil
}

func JsonString(value any) (string, error) {
	bytes, err := JsonBytes(value)
	return string(bytes), err
}
