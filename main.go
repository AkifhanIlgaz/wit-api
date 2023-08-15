package main

import (
	"fmt"
	"os"

	"github.com/AkifhanIlgaz/wit-api/firebase"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	credentialsFile := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")

	MyApp, err := firebase.NewApp(credentialsFile)
	if err != nil {
		panic(err)
	}

	fmt.Println(MyApp)
}
