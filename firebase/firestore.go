package firebase

import (
	"context"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
)

type Firestore struct {
	Client *firestore.Client
}

func NewFirestore(app firebase.App) *Firestore {
	firestore, err := app.Firestore(context.TODO())
	if err != nil {
		panic(err)
	}

	return &Firestore{
		Client: firestore,
	}
}
