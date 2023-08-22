package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/AkifhanIlgaz/wit-api/controllers"
	"github.com/AkifhanIlgaz/wit-api/firebase"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	credentialsFile := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")

	myApp, err := firebase.NewApp(credentialsFile)
	if err != nil {
		panic(err)
	}

	outfitsController := controllers.OutfitsController{
		OutfitService: *myApp.FirestoreService.OutfitService,
	}

	uidMiddleware := controllers.UidMiddleware{
		AuthService: myApp.AuthService,
	}

	r := chi.NewRouter()

	r.Use(uidMiddleware.SetUid)

	r.Route("/outfits", func(r chi.Router) {
		r.Get("/all/{uid}", outfitsController.GetAllOutfitsByUid)
		r.Get("/{outfitId}", outfitsController.GetOutfitById)
		r.Delete("/{outfitId}", outfitsController.DeleteOutfit)
		r.Post("/add", outfitsController.Add)
	})

	fmt.Println("Starting app")
	http.ListenAndServe(":3000", r)

}
