package firebase

import (
	"context"
	"fmt"

	"firebase.google.com/go/v4/auth"
)

type AuthService struct {
	Client *auth.Client
}

func (service *AuthService) GetUidByIdToken(idToken string) (string, error) {
	token, err := service.Client.VerifyIDToken(context.TODO(), idToken)
	if err != nil {
		return "", fmt.Errorf("get uid by id token: %w", err)
	}

	return token.UID, nil
}
