package firebase

import (
	"context"

	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go/v4"
)

type StorageService struct {
	Bucket *storage.BucketHandle
}

func NewStorageService(app firebase.App) *StorageService {
	client, err := app.Storage(context.TODO())
	if err != nil {
		panic(err)
	}

	bucket, err := client.DefaultBucket()
	if err != nil {
		panic(err)
	}

	return &StorageService{Bucket: bucket}
}
