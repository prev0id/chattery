package user

import (
	"chattery/internal/domain"
	"chattery/internal/utils/errors"
	"chattery/internal/utils/logger"
	"chattery/internal/utils/render"
	"context"
	"net/http"
	"time"
)

const sessionCookieName = "__Session"

func (s *Service) SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(sessionCookieName)
		if err != nil {
			clearSessionCookie(w)
			next.ServeHTTP(w, r)
			return
		}

		ctx := r.Context()
		session := domain.Session(cookie.Value)
		// save expiration before GetEx
		expiresAt := time.Now().Add(s.expiration)

		user := s.cache.UsernameFromSession(ctx, session, s.expiration)

		if user != domain.UserUnknown {
			ctx = domain.UsernameToContext(ctx, domain.Username(user))
			writeSessionsCookie(w, session, expiresAt)
		} else {
			clearSessionCookie(w)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *Service) CreateSession(ctx context.Context, w http.ResponseWriter, user domain.Username) error {
	session := domain.NewSession()

	expiresAt := time.Now().Add(s.expiration)

	if err := s.cache.WriteSession(ctx, session, user, s.expiration); err != nil {
		return errors.E(err).Debug("s.cache.WriteSession")
	}

	writeSessionsCookie(w, session, expiresAt)

	return nil
}

func (s *Service) ClearSession(ctx context.Context, w http.ResponseWriter, session domain.Session) {
	if err := s.cache.ExpireSession(ctx, session); err != nil {
		logger.ErrorCtx(ctx, err, "s.cache.ExpireSession")
	}
	clearSessionCookie(w)
}

func (s *Service) AuthRequiredMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if domain.UsernameFromContext(r.Context()) != domain.UserUnknown {
			next.ServeHTTP(w, r)
			return
		}
		render.Error(w, r, errors.E().Kind(errors.Unauthorized).Message("login required"))
	})
}

func clearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
}

func writeSessionsCookie(w http.ResponseWriter, session domain.Session, expiresAt time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    session.String(),
		Path:     "/",
		Expires:  expiresAt,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
}
