package firebase

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

type MyApp struct {
	Auth      *Auth
	Firestore *Firestore
	Storage   *Storage
}

func NewApp(credentialsFile string) (*MyApp, error) {
	ctx := context.TODO()

	opt := option.WithCredentialsFile(credentialsFile)
	config := firebase.Config{
		StorageBucket: "wearittomorrow-ab06f.appspot.com",
	}

	app, err := firebase.NewApp(ctx, &config, opt)
	if err != nil {
		return nil, fmt.Errorf("new app: %w", err)
	}

	return &MyApp{
		Auth:      NewAuth(*app),
		Firestore: NewFirestore(*app),
		Storage:   NewStorage(*app),
	}, nil
}
