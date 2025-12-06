package pb

import (
	_ "embed"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

//go:embed user_service/user_service.swagger.json
var swaggerJSON []byte

const (
	docJSONPath = "/swagger/doc.json"
	swaggerPath = "/swagger/*"
)

func RegisterSwaggerHandlers(mux *runtime.ServeMux) {
	swaggerHandler := httpSwagger.Handler(httpSwagger.URL(docJSONPath))

	mux.HandlePath(http.MethodGet, swaggerPath, func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
		swaggerHandler(w, r)
	})

	mux.HandlePath(http.MethodGet, docJSONPath, func(w http.ResponseWriter, _ *http.Request, _ map[string]string) {
		w.Write(swaggerJSON)
	})
}
