package user_proxy

import "github.com/go-chi/chi/v5"

type Proxy struct {
}

func NewProxy() *Proxy {
	return nil
}

func (p *Proxy) RegisterHandlers(router *chi.Router) {

}
