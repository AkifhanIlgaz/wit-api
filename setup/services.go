package setup

import (
	"fmt"
	"os"

	"github.com/AkifhanIlgaz/wit-api/firebase"
	"github.com/AkifhanIlgaz/wit-api/models"
)

type services struct {
	MyApp         *firebase.MyApp
	OutfitService *models.OutfitService
	UserService   *models.UserService
}

func Services() (*services, error) {
	credentialsFile := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")

	myApp, err := firebase.NewApp(credentialsFile)
	if err != nil {
		return nil, fmt.Errorf("setup services | firebase: %w", err)
	}

	outfitService := &models.OutfitService{
		Client: myApp.Firestore.Client,
	}

	userService := &models.UserService{
		Client: myApp.Firestore.Client,
	}

	return &services{
		MyApp:         myApp,
		OutfitService: outfitService,
		UserService:   userService,
	}, nil
}
