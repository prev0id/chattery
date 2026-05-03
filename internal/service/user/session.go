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
		if !session.Valid() {
			s.clearSessionCookie(w)
			next.ServeHTTP(w, r)
			return
		}

		ctx := r.Context()
		expiresAt := time.Now().Add(s.expiration)

		userID, refreshed, err := s.cache.UserIDFromSession(ctx, session, s.expiration, s.refreshBefore)
		if errutil.Is(errutil.NotFound, err) {
			s.clearSessionCookie(w)
			next.ServeHTTP(w, r)
			return
		}
		if err != nil {
			logger.ErrorCtx(ctx, err, "s.cache.UserIDFromSession")
			ctx = domain.SessionErrorToContext(ctx, err)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		if userID != domain.UserIsUnknown {
			ctx = domain.UserIDToContext(ctx, userID)
			if refreshed {
				s.writeSessionsCookie(w, session, expiresAt)
			}
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
	if session == domain.NoSession {
		s.clearSessionCookie(w)
		return
	}
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
		if err := domain.SessionErrorFromContext(r.Context()); err != nil {
			render.Error(w, r, errutil.E(err).Kind(errutil.Internal).Message("session storage unavailable"))
			return
		}
		render.Error(w, r, errutil.E().Kind(errutil.Unauthorized).Message("login required"))
	})
}

func (s *Service) clearSessionCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     domain.SessionCookieName,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   !s.debug,
		SameSite: http.SameSiteLaxMode,
	}
	if s.cookieDomain != "" {
		cookie.Domain = s.cookieDomain
	}
	http.SetCookie(w, cookie)
}

func (s *Service) writeSessionsCookie(w http.ResponseWriter, session domain.Session, expiresAt time.Time) {
	cookie := &http.Cookie{
		Name:     domain.SessionCookieName,
		Value:    session.String(),
		Path:     "/",
		Expires:  expiresAt,
		HttpOnly: true,
		Secure:   !s.debug,
		SameSite: http.SameSiteLaxMode,
	}
	if s.cookieDomain != "" {
		cookie.Domain = s.cookieDomain
	}
	http.SetCookie(w, cookie)
}
