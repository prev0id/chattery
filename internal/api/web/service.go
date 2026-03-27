package web_api

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	web "chattery/web/dist"
)

type Server struct {
}

func New() *Server {
	file, err := web.Assets.ReadFile("dist/assets/Toast-CxIVObqD.js")
	fmt.Println(err)
	fmt.Println(string(file))
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
		w.Write(web.LoginPage)
	})

	// router.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write(web.SignupPage)
	// })

	router.Handle("GET /assets/*", http.FileServer(http.FS(web.Assets)))
}
