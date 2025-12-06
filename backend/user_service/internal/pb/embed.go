package pb

import (
	_ "embed"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

//go:embed chat_service/chat_service.swagger.json
var swaggerJSON []byte

func RegisterSwaggerHandlers() {
	http.HandleFunc("GET /swagger/", httpSwagger.Handler(httpSwagger.URL("/swagger/doc.json")))
	http.HandleFunc("GET /swagger/doc.json", func(w http.ResponseWriter, _ *http.Request) { w.Write(swaggerJSON) })

}
