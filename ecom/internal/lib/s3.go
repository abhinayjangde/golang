package lib

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Client struct {
	svc     *awss3.Client
	presign *awss3.PresignClient
	bucket  string
}

func NewS3Client(region, bucket string) (*S3Client, error) {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))

	if err != nil {
		return nil, fmt.Errorf("load aws config: %w", err)
	}

	svc := awss3.NewFromConfig(cfg)

	return &S3Client{
		svc:     svc,
		presign: awss3.NewPresignClient(svc),
		bucket:  bucket,
	}, nil
}

func (c *S3Client) ObjectKey(listingID, filename string) string {
	return fmt.Sprintf("listings/%s/%s", listingID, filename)
}

func (c *S3Client) PresignedPutURL(ctx context.Context, key string, ttl time.Duration) (string, error) {
	resp, err := c.presign.PresignPutObject(ctx, &awss3.PutObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	}, awss3.WithPresignExpires(ttl))
	if err != nil {
		return "", fmt.Errorf("presign put object: %w", err)
	}
	return resp.URL, nil
}

func (c *S3Client) GetObject(ctx context.Context, key string) ([]byte, error) {
	out, err := c.svc.GetObject(ctx, &awss3.GetObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("get object %q: %w", key, err)
	}
	defer out.Body.Close()

	data, err := io.ReadAll(out.Body)
	if err != nil {
		return nil, fmt.Errorf("read object %q: %w", key, err)
	}
	return data, nil
}

func (c *S3Client) PutObject(ctx context.Context, key string, data []byte, contentType string) error {
	_, err := c.svc.PutObject(ctx, &awss3.PutObjectInput{
		Bucket:      aws.String(c.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return fmt.Errorf("put object %q: %w", key, err)
	}
	return nil
}

func (c *S3Client) DeleteObject(ctx context.Context, key string) error {
	_, err := c.svc.DeleteObject(ctx, &awss3.DeleteObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("delete object %q: %w", key, err)
	}
	return nil
}

// Bucket exposes the bucket name (for building public URLs).
func (c *S3Client) Bucket() string { return c.bucket }
