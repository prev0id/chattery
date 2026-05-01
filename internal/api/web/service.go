package web

import (
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"

	"chattery/internal/domain"
	"chattery/internal/utils/errutil"
	"chattery/internal/utils/render"
	"chattery/web"
)

type Server struct {
}

func New() *Server {
	return &Server{}
}

func (*Server) Pattern() string {
	return "/"
}

func (*Server) Route(router chi.Router) {
	router.HandleFunc("/app", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/app/dm", http.StatusFound)
	})
	router.HandleFunc("/app/*", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write(web.AppPage)
	})

	router.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if domain.UserIDFromContext(r.Context()) != domain.UserIsUnknown {
			http.Redirect(w, r, "/app/dm", http.StatusFound)
			return
		}
		_, _ = w.Write(web.LoginPage)
	})

	router.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		if domain.UserIDFromContext(r.Context()) != domain.UserIsUnknown {
			http.Redirect(w, r, "/app/dm", http.StatusFound)
			return
		}
		_, _ = w.Write(web.SignupPage)
	})

	assetsFS := http.FileServer(http.FS(web.Assets))

	router.HandleFunc("GET /assets/*", func(w http.ResponseWriter, r *http.Request) {
		newPath, err := url.JoinPath("/dist", r.URL.Path)
		if err != nil {
			render.Error(w, r, errutil.E(err).Kind(errutil.InvalidRequest))
			return
		}

		r.URL.Path = newPath

		assetsFS.ServeHTTP(w, r)
	})
}
