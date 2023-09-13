package main

import (
	"fmt"
	"os"
	"time"

	"github.com/AkifhanIlgaz/wit-api/firebase"
	"github.com/AkifhanIlgaz/wit-api/models"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	credentialsFile := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")

	myApp, err := firebase.NewApp(credentialsFile)
	if err != nil {
		panic(err)
	}

	outfitService := &models.OutfitService{
		Client: myApp.Firestore.Client,
	}

	outfits, err := outfitService.GetOutfits([]string{"a", "b", "c", "d"}, time.Now())
	if err != nil {
		panic(err)
	}

	for _, outfit := range outfits {
		fmt.Printf("%+v\n", outfit)
	}
}
