package s3

import (
	"bytes"
	"context"
	"io"

	"chattery/internal/config"
	"chattery/internal/domain"
	"chattery/internal/utils/errutil"
)

const contentTypeJPEG = "image/jpeg"

type client interface {
	ObjectExists(ctx context.Context, bucket string, name string) (bool, error)
	GetObject(ctx context.Context, bucket string, name string) (io.ReadCloser, error)
	PutObject(ctx context.Context, bucket string, name string, reader io.Reader, size int64, contentType string) error
}

type Adapter struct {
	client client
	bucket string
}

func New(client client, cfg *config.Config) *Adapter {
	return &Adapter{
		client: client,
		bucket: cfg.S3.Bucket,
	}
}

func (a *Adapter) ImageExists(ctx context.Context, imageID domain.ImageID) (bool, error) {
	exists, err := a.client.ObjectExists(ctx, a.bucket, imageID.String())
	if err != nil {
		return false, errutil.E(err).Debug("a.client.ObjectExists")
	}

	return exists, nil
}

func (a *Adapter) GetImage(ctx context.Context, imageID domain.ImageID) ([]byte, error) {
	object, err := a.client.GetObject(ctx, a.bucket, imageID.String())
	if err != nil {
		return nil, errutil.E(err).Debug("a.client.GetObject")
	}
	defer object.Close()

	imgBytes, err := io.ReadAll(object)
	if err != nil {
		return nil, errutil.E(err).Debug("io.ReadAll")
	}

	return imgBytes, nil
}

func (a *Adapter) PutJPEGImage(ctx context.Context, imageID domain.ImageID, imgBytes []byte) error {
	if err := a.client.PutObject(
		ctx,
		a.bucket,
		imageID.String(),
		bytes.NewReader(imgBytes),
		int64(len(imgBytes)),
		contentTypeJPEG,
	); err != nil {
		return errutil.E(err).Debug("a.client.PutObject")
	}

	return nil
}
