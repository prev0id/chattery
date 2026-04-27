package bind

import (
	"io"
	"net/http"

	"chattery/internal/utils/errutil"
)

func JSON[T any](request *http.Request) (*T, error) {
	body, err := io.ReadAll(request.Body)
	if err != nil {
		return nil, errutil.E(err).Kind(errutil.InvalidRequest).Debug("io.ReadAll")
	}

	return JSONBytes[T](body)
}
