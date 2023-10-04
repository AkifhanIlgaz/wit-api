package firebase

import (
	"context"
	"fmt"
	"path"
	"time"

	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go/v4"
)

type Storage struct {
	BaseUrl string
	Bucket  *storage.BucketHandle
}

func NewStorage(app firebase.App) *Storage {
	client, err := app.Storage(context.TODO())
	if err != nil {
		panic(err)
	}

	bucket, err := client.DefaultBucket()
	if err != nil {
		panic(err)
	}

	return &Storage{
		BaseUrl: "https://storage.cloud.google.com/wearittomorrow-ab06f.appspot.com",
		Bucket:  bucket}
}

func (service *Storage) GenerateUploadUrl(uid string, timestamp int64, extension string, dir string) (string, string, error) {
	objPath := fmt.Sprintf("%s/%s-%v.%s", dir, uid, timestamp, extension)

	signedUrl, err := service.Bucket.SignedURL(objPath, &storage.SignedURLOptions{
		Method:      "PUT",
		Expires:     time.Now().Add(15 * time.Minute),
		ContentType: fmt.Sprintf("image/%s", extension),
		Scheme:      storage.SigningSchemeV4,
	})
	if err != nil {
		return "", "", fmt.Errorf("generate upload url: %w", err)
	}

	return signedUrl, objPath, nil
}

func (service *Storage) GetDownloadUrl(objPath string) string {
	return path.Join(service.BaseUrl, objPath)
}
