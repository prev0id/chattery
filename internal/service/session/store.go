package session

import (
	"chattery/internal/config"
	"chattery/internal/domain"
	"chattery/internal/utils/errors"
	"chattery/internal/utils/render"
	"net/http"

	"github.com/boj/redistore"
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
	store, err := redistore.NewRediStoreWithPool(pool, []byte(cfg.Session.SecretKey))
	if err != nil {
		return nil, errors.E(err).Debug("redistore.NewRediStoreWithPool")
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
			render.Error(w, r, errors.E(err).Kind(errors.Unauthorized).Message("invalid session"))
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

		session.Save(r, w)
	})
}

func (s *Store) WriteUsernameToSession(r *http.Request, user domain.Username) error {
	session, err := s.sessions.Get(r, sessionName)
	if err != nil {
		return errors.E(err).Kind(errors.Unauthorized).Message("invalid session")
	}
	session.Values[usernameSessionKey] = user.String()
	return nil
}

func (s *Store) AuthRequiredMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if UsernameFromContext(r.Context()) != domain.UserUnknown {
			next.ServeHTTP(w, r)
			return
		}
		render.Error(w, r, errors.E().Kind(errors.Unauthorized).Message("login required"))
	})
}
