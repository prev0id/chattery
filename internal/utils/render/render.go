package render

import (
	"encoding/json"
	"net/http"
	"time"

	"chattery/internal/utils/errors"
	"chattery/internal/utils/logger"
)

const (
	formatTime     = "15:04"
	formatDate     = "Jan 2"
	formatDateTime = formatDate + ", " + formatTime
)

func Timestamp(t time.Time) string {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	yesterday := today.AddDate(0, 0, -1)
	msgDate := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())

	if msgDate.Equal(today) {
		return "Today, " + t.Format(formatTime)
	}
	if msgDate.Equal(yesterday) {
		return "Yesterday, " + t.Format(formatTime)
	}
	return t.Format(formatDateTime)
}

func Json[T any](w http.ResponseWriter, r *http.Request, value T) {
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
	logger.ErrorCtx(r.Context(), err, "request ended with an error")

	statusCode := errors.E(err).GetKind().StatusCode()

	response, _ := json.Marshal(responseError{Message: err.Error()})

	setContentTypeJSON(w)
	w.WriteHeader(statusCode)
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
