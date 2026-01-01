package bind

import (
	"chattery/internal/utils/errors"
	"encoding/json"
	"io"
	"net/http"
)

func JSON[T any](request *http.Request) (*T, error) {
	var res *T

	body, err := io.ReadAll(request.Body)
	if err != nil {
		return nil, errors.E(err).Kind(errors.InvalidRequest).Debug("io.ReadAll")
	}

	if err := json.Unmarshal(body, res); err != nil {
		return nil, errors.E(err).Kind(errors.InvalidRequest).Debug("json.Unmarshal")
	}

	return res, nil
}
