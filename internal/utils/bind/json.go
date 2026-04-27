package bind

import (
	"encoding/json"

	"chattery/internal/utils/errutil"
)

func JSONString[T any](raw string) (*T, error) {
	return JSONBytes[T]([]byte(raw))
}

func JSONBytes[T any](raw []byte) (*T, error) {
	result := new(T)
	if err := json.Unmarshal(raw, result); err != nil {
		return nil, errutil.E(err).
			Kind(errutil.InvalidRequest).
			Message("invalid json provided").
			Debug("json.Unmarshal")
	}
	return result, nil
}
