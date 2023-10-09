package setup

import (
	"os"

	"github.com/AkifhanIlgaz/wit-api/firebase"
	"github.com/AkifhanIlgaz/wit-api/models"
)

func Services() (*firebase.MyApp, *models.OutfitService, *models.UserService) {
	credentialsFile := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")

	myApp, err := firebase.NewApp(credentialsFile)
	if err != nil {
		panic(err)
	}

	outfitService := &models.OutfitService{
		Client: myApp.Firestore.Client,
	}

	userService := &models.UserService{
		Client: myApp.Firestore.Client,
	}

	return myApp, outfitService, userService
}
