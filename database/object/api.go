package object

import (
	"context"
	"io"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"

	"github.com/kiwiworks/rodent/system/errors"
)

type UploadedObject struct {
	Uri            *url.URL
	ChecksumSHA256 string
	UploadedSize   int64
}

func (s *Store) Upload(ctx context.Context, bucketName, path string, data io.Reader, objectSize int64) (
	*UploadedObject,
	error,
) {
	info, err := s.client.PutObject(ctx, bucketName, path, data, objectSize, minio.PutObjectOptions{})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to upload object '%s' to bucket '%s'", path, bucketName)
	}
	return &UploadedObject{
		Uri: &url.URL{
			Scheme: "s3",
			Host:   info.Bucket,
			Path:   info.Key,
			RawQuery: url.Values{
				"etag": []string{info.ETag},
			}.Encode(),
		},
		ChecksumSHA256: info.ChecksumSHA256,
		UploadedSize:   info.Size,
	}, nil
}

func (s *Store) PreSignedDownload(ctx context.Context, uri url.URL, expires time.Duration) (
	*url.URL,
	error,
) {
	bucket, path, err := pathAndKey(uri)
	if err != nil {
		return nil, err
	}
	u, err := s.client.PresignedGetObject(ctx, bucket, path, expires, url.Values{})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create pre-signed download link for object '%s' in bucket '%s'", path, bucket)
	}
	return u, nil
}
