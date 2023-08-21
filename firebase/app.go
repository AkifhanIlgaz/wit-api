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

	return &MyApp{
		AuthService:      NewAuthService(*app),
		FirestoreService: NewFirestoreService(*app),
	}, nil
}
