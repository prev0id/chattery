package user

import (
	"context"
	"net/http"
	"time"

	"chattery/internal/domain"
	"chattery/internal/utils/errutil"
	"chattery/internal/utils/logger"
	"chattery/internal/utils/render"
)

func (s *Service) SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session := domain.GetSessionFromRequest(r)

		if session == domain.NoSession {
			next.ServeHTTP(w, r)
			return
		}

		ctx := r.Context()
		// save expiration before GetEx
		expiresAt := time.Now().Add(s.expiration)

		userID := s.cache.UserIDFromSession(ctx, session, s.expiration)

		if userID != domain.UserIsUnknown {
			ctx = domain.UserIDToContext(ctx, userID)
			s.writeSessionsCookie(w, session, expiresAt)
		} else {
			s.clearSessionCookie(w)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *Service) CreateSession(ctx context.Context, w http.ResponseWriter, userID domain.UserID) error {
	session := domain.NewSession()

	expiresAt := time.Now().Add(s.expiration)

	if err := s.cache.WriteSession(ctx, session, userID, s.expiration); err != nil {
		return errutil.E(err).Debug("s.cache.WriteSession")
	}

	s.writeSessionsCookie(w, session, expiresAt)

	return nil
}

func (s *Service) ClearSession(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	session := domain.GetSessionFromRequest(r)
	if err := s.cache.ClearSession(ctx, session); err != nil {
		logger.ErrorCtx(ctx, err, "s.cache.ExpireSession")
	}
	s.clearSessionCookie(w)
}

func (*Service) AuthRequiredMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if domain.UserIDFromContext(r.Context()) != domain.UserIsUnknown {
			next.ServeHTTP(w, r)
			return
		}
		render.Error(w, r, errutil.E().Kind(errutil.Unauthorized).Message("login required"))
	})
}

func (s *Service) clearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     domain.SessionCookieName,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   !s.debug,
		SameSite: http.SameSiteLaxMode,
	})
}

func (s *Service) writeSessionsCookie(w http.ResponseWriter, session domain.Session, expiresAt time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     domain.SessionCookieName,
		Value:    session.String(),
		Path:     "/",
		Expires:  expiresAt,
		HttpOnly: true,
		Secure:   !s.debug,
		SameSite: http.SameSiteLaxMode,
	})
}
