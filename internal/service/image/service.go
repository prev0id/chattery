package image

import (
	"context"

	"github.com/rrivera/identicon"

	"chattery/internal/config"
	"chattery/internal/domain"
	"chattery/internal/utils/errutil"
)

const (
	identiconNamespace = "chatter"
	identiconBlocks    = 10
	identiconDensity   = 4
	identiconSize      = 300
)

type db interface {
	UserByID(ctx context.Context, userID domain.UserID) (*domain.User, error)
	UpdateUser(ctx context.Context, updated *domain.User) error
}

type objectStorage interface {
	ImageExists(ctx context.Context, imageID domain.ImageID) (bool, error)
	GetImage(ctx context.Context, imageID domain.ImageID) ([]byte, error)
	PutJPEGImage(ctx context.Context, imageID domain.ImageID, imgBytes []byte) error
}

type Service struct {
	db          db
	storage     objectStorage
	identicon   *identicon.Generator
	maxFileSize int64
	maxWidth    int
	maxHeight   int
	jpegQuality int
}

func New(db db, storage objectStorage, cfg *config.Config) (*Service, error) {
	generator, err := identicon.New(identiconNamespace, identiconBlocks, identiconDensity)
	if err != nil {
		return nil, errutil.E(err).Debug("identicon.New")
	}
	return &Service{
		db:          db,
		storage:     storage,
		maxFileSize: cfg.Image.MaxFileSize,
		maxWidth:    cfg.Image.MaxWidth,
		maxHeight:   cfg.Image.MaxHeight,
		jpegQuality: cfg.Image.JPEGQuality,
		identicon:   generator,
	}, nil
}

func objectName(username domain.Username) string {
	return username.String() + ".jpeg"
}
