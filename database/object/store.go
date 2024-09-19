package object

import (
	"context"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"rodent/system/errors"
)

type (
	Store struct {
		endpoint          string
		client            *minio.Client
		cancelHealthCheck context.CancelFunc
	}
	StoreConfig struct {
		Endpoint  string
		AccessKey string
		SecretKey string
		UseSSL    bool
	}
)

func NewStore(cfg *StoreConfig) (*Store, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:        credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure:       cfg.UseSSL,
		BucketLookup: minio.BucketLookupAuto,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "could not create minio client for endpoint '%s'", cfg.Endpoint)
	}
	return &Store{
		endpoint: cfg.Endpoint,
		client:   client,
	}, nil
}

func (s *Store) OnStart(context.Context) error {
	cancel, err := s.client.HealthCheck(time.Second * 30)
	if err != nil {
		return errors.Wrapf(err, "could not start healthcheck for endpoint '%s'", s.endpoint)
	}
	s.cancelHealthCheck = cancel
	return nil
}

func (s *Store) OnStop(context.Context) error {
	if s.cancelHealthCheck != nil {
		s.cancelHealthCheck()
	}
	return nil
}
