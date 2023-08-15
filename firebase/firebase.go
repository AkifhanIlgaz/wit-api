package firebase

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"firebase.google.com/go/v4/storage"
	"google.golang.org/api/option"
)

type MyApp struct {
	Auth      *auth.Client
	Firestore *firestore.Client
	Storage   *storage.Client
}

func NewApp(credentialsFile string) (*MyApp, error) {
	ctx := context.TODO()
	opt := option.WithCredentialsFile(credentialsFile)

	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, fmt.Errorf("new app: %w", err)
	}

	auth, err := app.Auth(ctx)
	if err != nil {
		return nil, fmt.Errorf("new app | auth: %w", err)
	}

	firestore, err := app.Firestore(ctx)
	if err != nil {
		return nil, fmt.Errorf("new app | firestore: %w", err)
	}

	storage, err := app.Storage(ctx)
	if err != nil {
		return nil, fmt.Errorf("new app | storage: %w", err)
	}

	return &MyApp{
		Auth:      auth,
		Storage:   storage,
		Firestore: firestore,
	}, nil
}
