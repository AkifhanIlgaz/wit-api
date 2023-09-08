package firebase

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

type MyApp struct {
	AuthService      *AuthService
	FirestoreService *FirestoreService
	StorageService   *StorageService
}

func NewApp(credentialsFile string) (*MyApp, error) {
	ctx := context.TODO()

	opt := option.WithCredentialsFile(credentialsFile)
	config := firebase.Config{
		StorageBucket: "wearittomorrow-ab06f.appspot.com/",
	}

	app, err := firebase.NewApp(ctx, &config, opt)
	if err != nil {
		return nil, fmt.Errorf("new app: %w", err)
	}

	return &MyApp{
		AuthService:      NewAuthService(*app),
		FirestoreService: NewFirestoreService(*app),
		StorageService:   NewStorageService(*app),
	}, nil
}
