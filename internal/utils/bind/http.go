package bind

import (
	"chattery/internal/utils/errs"
	"encoding/json"
	"io"
	"net/http"
)

func HttpBody[T any](request *http.Request) (*T, error) {
	var res *T

	body, err := io.ReadAll(request.Body)
	if err != nil {
		return nil, errs.E(err, errs.InvalidRequest, errs.Debug("io.ReadAll"))
	}

	if err := json.Unmarshal(body, res); err != nil {
		return nil, errs.E(err, errs.InvalidRequest, errs.Debug("json.Unmarshal"))
	}

	return res, nil
}
