package service

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3Client клиент для работы с внешним S3
type S3Client struct {
	client   *s3.Client
	bucket   string
	endpoint string
	region   string
}

// NewS3ClientWithRegion creates a new S3 client with configurable region.
// Returns an error instead of panicking on config load failure.
func NewS3ClientWithRegion(endpoint, accessKey, secretKey, bucket, region string) (*S3Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
			return aws.Credentials{
				AccessKeyID:     accessKey,
				SecretAccessKey: secretKey,
			}, nil
		})),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	return &S3Client{
		client:   s3.NewFromConfig(cfg, func(o *s3.Options) {
			o.BaseEndpoint = aws.String(endpoint)
			o.UsePathStyle = true
		}),
		bucket:   bucket,
		endpoint: endpoint,
		region:   region,
	}, nil
}

// NewS3Client creates a new S3 client with default us-east-1 region.
// Kept for backward compatibility with existing tests.
func NewS3Client(endpoint, accessKey, secretKey, bucket string) *S3Client {
	client, err := NewS3ClientWithRegion(endpoint, accessKey, secretKey, bucket, "us-east-1")
	if err != nil {
		panic(err)
	}
	return client
}

// UploadAvatar загружает аватар пользователя в S3
func (s *S3Client) UploadAvatar(ctx context.Context, sid string, data []byte) (string, error) {
	key := fmt.Sprintf("avatars/%s.jpg", sid)
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String("image/jpeg"),
	})
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s/%s", s.endpoint, s.bucket, key), nil
}

// UploadPGN загружает запись партии в S3
func (s *S3Client) UploadPGN(ctx context.Context, gameType, matchID, pgn string) (string, error) {
	key := fmt.Sprintf("records/%s/%s.pgn", gameType, matchID)
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        strings.NewReader(pgn),
		ContentType: aws.String("text/plain"),
	})
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s/%s", s.endpoint, s.bucket, key), nil
}

// GetAvatarURL получает URL аватара пользователя
func (s *S3Client) GetAvatarURL(ctx context.Context, sid string) (string, error) {
	key := fmt.Sprintf("avatars/%s.jpg", sid)
	_, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return "", err // Аватар не найден
	}
	return fmt.Sprintf("%s/%s/%s", s.endpoint, s.bucket, key), nil
}
