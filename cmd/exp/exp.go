package main

import (
	"fmt"
	"os"

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

	fmt.Println(outfitService.GetOutfitCountOfUser("xQFQncknojU5vUnsmIl2bIevBdE2"))

	fmt.Println(x("outfit"))

}

func x(fileType string) string {
	var dir string

	switch fileType {
	case "outfit":
		return "outfits"
	}
	return dir

}
