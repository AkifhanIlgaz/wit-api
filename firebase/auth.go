package firebase

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
)

type AuthService struct {
	Client *auth.Client
}

func NewAuthService(app firebase.App) *AuthService {
	auth, err := app.Auth(context.TODO())
	if err != nil {
		panic(err)
	}
	return &AuthService{
		Client: auth,
	}
}

func (service *AuthService) GetUidByIdToken(idToken string) (string, error) {
	token, err := service.Client.VerifyIDToken(context.TODO(), idToken)
	if err != nil {
		return "", fmt.Errorf("get uid by id token: %w", err)
	}

	return token.UID, nil
}
