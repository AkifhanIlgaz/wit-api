package main

import (
	"context"
	"fmt"
	"os"

	"github.com/AkifhanIlgaz/wit-api/firebase"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	credentialsFile := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")

	myApp, err := firebase.NewApp(credentialsFile)
	if err != nil {
		panic(err)
	}

	u, err := myApp.Auth.Client.GetUser(context.TODO(), "xQFQncknojU5vUnsmIl2bIevBdE2")

	fmt.Println(u.UserInfo.PhotoURL)
}

func x(fileType string) string {
	var dir string

	switch fileType {
	case "outfit":
		return "outfits"
	}
	return dir

}
