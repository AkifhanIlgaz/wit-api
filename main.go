package main

import (
	"fmt"
	"net/http"

	"github.com/AkifhanIlgaz/wit-api/setup"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	myApp, outfitService, userService := setup.Services()

	uidMiddleware, firebaseController, outfitsController, usersController := setup.Controllers(myApp, outfitService, userService)

	r := setup.Routes(uidMiddleware, firebaseController, outfitsController, usersController)

	fmt.Println("Starting app")
	http.ListenAndServe(":3000", r)
}
