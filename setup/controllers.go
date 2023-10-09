package setup

import (
	"github.com/AkifhanIlgaz/wit-api/controllers"
	"github.com/AkifhanIlgaz/wit-api/firebase"
	"github.com/AkifhanIlgaz/wit-api/models"
)

func Controllers(myApp *firebase.MyApp, outfitService *models.OutfitService, userService *models.UserService) (*controllers.UidMiddleware, *controllers.FirebaseController, *controllers.OutfitsController, *controllers.UsersController) {
	uidMiddleware := controllers.UidMiddleware{
		Auth: myApp.Auth,
	}

	firebaseController := controllers.FirebaseController{
		Storage: myApp.Storage,
	}

	outfitsController := controllers.OutfitsController{
		Storage:       myApp.Storage,
		OutfitService: outfitService,
		UserService:   userService,
	}

	usersController := controllers.UsersController{
		UserService:   userService,
		OutfitService: outfitService,
		Auth:          myApp.Auth,
		Storage:       myApp.Storage,
	}

	return &uidMiddleware, &firebaseController, &outfitsController, &usersController
}
