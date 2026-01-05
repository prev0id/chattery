package bind

import (
	"encoding/json"
	"io"
	"net/http"

	"chattery/internal/utils/errors"
)

func Json[T any](request *http.Request) (*T, error) {
	body, err := io.ReadAll(request.Body)
	if err != nil {
		return nil, errors.E(err).Kind(errors.InvalidRequest).Debug("io.ReadAll")
	}

	return JsonBytes[T](body)
}

func JsonString[T any](raw string) (*T, error) {
	return JsonBytes[T]([]byte(raw))
}

func JsonBytes[T any](raw []byte) (*T, error) {
	result := new(T)
	if err := json.Unmarshal(raw, result); err != nil {
		return nil, errors.E(err).
			Kind(errors.InvalidRequest).
			Message("invalid json provided").
			Debug("json.Unmarshal")
	}
	return result, nil
}
