package dm

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"

	"chattery/internal/domain"
)

type dmService interface {
	UserDMs(ctx context.Context, userID domain.UserID) ([]*domain.DM, error)
	SearchUsersWithoutDM(ctx context.Context, userID domain.UserID, query string) ([]*domain.User, error)
	CreateDM(ctx context.Context, participant1, participant2 domain.UserID) (domain.DMID, error)
	CreateDMMessage(ctx context.Context, message *domain.DMMessage) error
	FirstPageOfDMMessages(ctx context.Context, userID domain.UserID, cursor *domain.DMCursor) ([]*domain.DMMessage, *domain.DMCursor, error)
	NextPagesOfDMMessages(ctx context.Context, userID domain.UserID, cursor *domain.DMCursor) ([]*domain.DMMessage, *domain.DMCursor, error)
}

type userService interface {
	AuthRequiredMiddleware(next http.Handler) http.Handler
}

type userCache interface {
	GetByID(id domain.UserID) (*domain.User, error)
	ListByID() map[domain.UserID]*domain.User
}

type Server struct {
	user  userService
	dm    dmService
	cache userCache
}

func New(user userService, dm dmService, cache userCache) *Server {
	return &Server{
		user:  user,
		dm:    dm,
		cache: cache,
	}
}

func (*Server) Pattern() string {
	return "/v1/dm"
}

func (s *Server) Route(router chi.Router) {
	router.Group(func(withAuthRouter chi.Router) {
		withAuthRouter.Use(s.user.AuthRequiredMiddleware)

		withAuthRouter.Get("/list", s.GetDMs)
		withAuthRouter.Get("/search", s.SearchUsers)
		withAuthRouter.Post("/create", s.PostCreateDM)
		withAuthRouter.Post("/message", s.PostMessage)
		withAuthRouter.Post("/messages", s.GetMessages)
	})
}
