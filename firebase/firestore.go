package firebase

import (
	"context"

	firebase "firebase.google.com/go/v4"
	"github.com/AkifhanIlgaz/wit-api/models"
)

type FirestoreService struct {
	OutfitService *models.OutfitService
}

func NewFirestoreService(app firebase.App) *FirestoreService {
	firestore, err := app.Firestore(context.TODO())
	if err != nil {
		panic(err)
	}

	return &FirestoreService{
		OutfitService: &models.OutfitService{
			Client: firestore,
		},
	}
}
