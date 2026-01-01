package session

import (
	"chattery/internal/config"
	"chattery/internal/domain"
	"chattery/internal/utils/errs"
	"net/http"

	"github.com/boj/redistore"
	"github.com/go-chi/render"
	"github.com/gomodule/redigo/redis"
)

const (
	sessionName        = "__Session"
	usernameSessionKey = "username"
)

type Store struct {
	sessions *redistore.RediStore
}

func New(cfg *config.Config, pool *redis.Pool) (*Store, error) {
	store, err := redistore.NewRediStoreWithPool(pool, []byte(cfg.SessionSecretKey))
	if err != nil {
		return nil, errs.E(err, errs.Debug("redistore.NewRediStoreWithPool"))
	}

	return &Store{sessions: store}, nil
}

func (s *Store) Close() error {
	return s.sessions.Close()
}

func (s *Store) SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := s.sessions.Get(r, sessionName)
		if err != nil {
			render.Render(w, r, errs.FromError(err))
			return
		}

		var username string
		if value, ok := session.Values[usernameSessionKey]; ok {
			if fromSession, ok := value.(string); ok {
				username = fromSession
			}
		}

		ctx := r.Context()

		if username != "" {
			ctx = usernameToContext(ctx, domain.Username(username))
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *Store) AuthRequiredMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if UsernameFromContext(r.Context()) != domain.UnknownUsername {
			next.ServeHTTP(w, r)
			return
		}
	})
}
