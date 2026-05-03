package s3

import (
	"context"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"chattery/internal/config"
)

type Client struct {
	client *minio.Client
}

func New(cfg config.S3) (*Client, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, err
	}

	return &Client{
		client: client,
	}, nil
}

func (c *Client) ObjectExists(ctx context.Context, bucket string, name string) (bool, error) {
	_, err := c.client.StatObject(ctx, bucket, name, minio.StatObjectOptions{})
	if err != nil {
		response := minio.ToErrorResponse(err)
		if response.Code == "NoSuchKey" || response.Code == "NoSuchBucket" {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (c *Client) GetObject(ctx context.Context, bucket string, name string) (io.ReadCloser, error) {
	return c.client.GetObject(ctx, bucket, name, minio.GetObjectOptions{})
}

func (c *Client) PutObject(ctx context.Context, bucket string, name string, reader io.Reader, size int64, contentType string) error {
	_, err := c.client.PutObject(
		ctx,
		bucket,
		name,
		reader,
		size,
		minio.PutObjectOptions{ContentType: contentType},
	)
	return err
}
