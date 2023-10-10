package setup

import (
	ctrls "github.com/AkifhanIlgaz/wit-api/controllers"
)

type controllers struct {
	UidMiddleware      *ctrls.UidMiddleware
	FirebaseController *ctrls.FirebaseController
	OutfitsController  *ctrls.OutfitsController
	UsersController    *ctrls.UsersController
}

func Controllers(services *services) *controllers {
	uidMiddleware := ctrls.UidMiddleware{
		Auth: services.MyApp.Auth,
	}

	firebaseController := ctrls.FirebaseController{
		Storage: services.MyApp.Storage,
	}

	outfitsController := ctrls.OutfitsController{
		Storage:       services.MyApp.Storage,
		OutfitService: services.OutfitService,
		UserService:   services.UserService,
	}

	usersController := ctrls.UsersController{
		UserService:   services.UserService,
		OutfitService: services.OutfitService,
		Auth:          services.MyApp.Auth,
		Storage:       services.MyApp.Storage,
	}

	return &controllers{
		UidMiddleware:      &uidMiddleware,
		FirebaseController: &firebaseController,
		OutfitsController:  &outfitsController,
		UsersController:    &usersController,
	}
}
