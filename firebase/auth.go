package firebase

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
)

type Auth struct {
	Client *auth.Client
}

func NewAuth(app firebase.App) *Auth {
	auth, err := app.Auth(context.TODO())
	if err != nil {
		panic(err)
	}

	return &Auth{
		Client: auth,
	}
}

func (auth *Auth) GetUidByIdToken(idToken string) (string, error) {
	token, err := auth.Client.VerifyIDToken(context.TODO(), idToken)
	if err != nil {
		return "", fmt.Errorf("get uid by id token: %w", err)
	}

	return token.UID, nil
}

func (a *Auth) UpdateProfilePhoto(uid string, photoUrl string) error {
	var user auth.UserToUpdate
	user = *(user.PhotoURL(photoUrl))

	_, err := a.Client.UpdateUser(context.TODO(), uid, &user)
	if err != nil {
		return fmt.Errorf("update profile photo: %w", err)
	}

	return nil
}
