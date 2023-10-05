package main

import (
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

	myApp.Auth.UpdateProfilePhoto("xQFQncknojU5vUnsmIl2bIevBdE2", "https://lh3.googleusercontent.com/a/ACg8ocI66bKaZPTj_ZzGiuajojbqzkTAeFyyCVg15CRLAUFj=s96-c")

}

func x(fileType string) string {
	var dir string

	switch fileType {
	case "outfit":
		return "outfits"
	}
	return dir

}
