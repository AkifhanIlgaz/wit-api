package firebase

import (
	"context"
	"fmt"
	"time"

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

func (service *StorageService) GenerateUploadUrl(uid, id, extension string) (string, error) {
	objPath := fmt.Sprintf("%s/%s.%s", uid, id, extension)

	signedUrl, err := service.Bucket.SignedURL(objPath, &storage.SignedURLOptions{
		Method:  "PUT",
		Expires: time.Now().Add(15 * time.Minute),
	})
	if err != nil {
		return "", fmt.Errorf("generate upload url: %w", err)
	}

	return signedUrl, nil
}
