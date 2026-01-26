package web_api

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"chattery/web"
)

type Server struct {
}

func New() *Server {
	return &Server{}
}

func (s *Server) Pattern() string {
	return "/"
}

func (s *Server) Route(router chi.Router) {
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(web.RootPage)
	})

	router.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		w.Write(web.LoginPage)
	})

	router.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		w.Write(web.SignupPage)
	})

	router.Handle("GET /src/*", http.FileServer(http.FS(web.Src)))
}
