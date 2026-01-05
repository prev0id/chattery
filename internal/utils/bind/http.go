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
	var res *T
	if err := json.Unmarshal([]byte(raw), res); err != nil {
		return nil, errors.E(err).Debug("json.Unmarshal")
	}
	return res, nil
}

func JsonBytes[T any](raw []byte) (*T, error) {
	var res *T
	if err := json.Unmarshal(raw, res); err != nil {
		return nil, errors.E(err).Debug("json.Unmarshal")
	}
	return res, nil
}
