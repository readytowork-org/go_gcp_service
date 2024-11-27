package go_gcp_service

import (
	"context"
	"errors"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

type bucketClientLogger interface {
	Fatalf(template string, args ...interface{})
	Fatal(args ...interface{})
	Error(args ...interface{})
}

type BucketClientConfig struct {
	Logger            bucketClientLogger
	StorageBucketName string
	ClientOption      *option.ClientOption
}

type BucketClient struct {
	*storage.Client
}

// NewGCPBucketClient creates a new gcp bucket api client
func NewGCPBucketClient(config BucketClientConfig) BucketClient {
	bucketName := config.StorageBucketName
	ctx := context.Background()
	if bucketName == "" {
		config.Logger.Error("Please check your env file for STORAGE_BUCKET_NAME")
	}
	client, err := storage.NewClient(ctx, *config.ClientOption)
	if err != nil {
		config.Logger.Fatal(err.Error())
	}

	bucket := client.Bucket(bucketName)
	_, err = bucket.Attrs(ctx)
	if errors.Is(err, storage.ErrBucketNotExist) {
		config.Logger.Fatalf("Provided bucket %v doesn't exists", bucketName)
	}

	if err != nil {
		config.Logger.Fatalf("Cloud bucket error: %v", err.Error())
	}

	bucketAttrsToUpdate := storage.BucketAttrsToUpdate{
		CORS: []storage.CORS{
			{
				MaxAge:          600,
				Methods:         []string{"PUT", "PATCH", "GET", "POST", "OPTIONS", "DELETE"},
				Origins:         []string{"*"},
				ResponseHeaders: []string{"Content-Type"},
			},
		},
	}
	if _, err := bucket.Update(ctx, bucketAttrsToUpdate); err != nil {
		config.Logger.Fatalf("Cloud bucket update error: %v", err.Error())
	}
	return BucketClient{
		client,
	}
}
