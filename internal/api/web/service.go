package web_api

import (
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"

	"chattery/internal/domain"
	"chattery/internal/utils/errors"
	"chattery/internal/utils/render"
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
	router.HandleFunc("/app", func(w http.ResponseWriter, r *http.Request) {
		w.Write(web.AppPage)
	})

	router.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if domain.UserIDFromContext(r.Context()) != domain.UserIsUnknown {
			http.Redirect(w, r, "app", http.StatusFound)
			return
		}
		w.Write(web.LoginPage)
	})

	router.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		if domain.UserIDFromContext(r.Context()) != domain.UserIsUnknown {
			http.Redirect(w, r, "app", http.StatusFound)
			return
		}
		w.Write(web.SignupPage)
	})

	assetsFS := http.FileServer(http.FS(web.Assets))

	router.HandleFunc("GET /assets/*", func(w http.ResponseWriter, r *http.Request) {
		newPath, err := url.JoinPath("/dist", r.URL.Path)
		if err != nil {
			render.Error(w, r, errors.E(err).Kind(errors.InvalidRequest))
			return
		}

		r.URL.Path = newPath

		assetsFS.ServeHTTP(w, r)
	})
}
