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

	return &MyApp{
		AuthService:      &AuthService{auth},
		FirestoreService: &FirestoreService{firestore},
	}, nil
}
