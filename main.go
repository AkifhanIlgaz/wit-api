package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/AkifhanIlgaz/wit-api/controllers"
	"github.com/AkifhanIlgaz/wit-api/ctx"
	"github.com/AkifhanIlgaz/wit-api/firebase"
	"github.com/AkifhanIlgaz/wit-api/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	credentialsFile := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")

	myApp, err := firebase.NewApp(credentialsFile)
	if err != nil {
		panic(err)
	}

	firebaseController := controllers.FirebaseController{
		Storage: myApp.Storage,
	}

	outfitService := &models.OutfitService{
		Client: myApp.Firestore.Client,
	}

	uidMiddleware := controllers.UidMiddleware{
		Auth: myApp.Auth,
	}

	r := chi.NewRouter()

	r.Use(uidMiddleware.SetUid)
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "idToken", "fileExtension"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Get("/generate-upload-url", firebaseController.GenerateUploadUrl)
	r.Post("/add-outfit", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Please provide body", http.StatusBadRequest)
			return
		}

		outfit := models.Outfit{
			Uid:       *ctx.Uid(r.Context()),
			CreatedAt: time.Now(),
		}
		err = json.Unmarshal(body, &outfit)
		if err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
		outfit.PhotoURL = myApp.Storage.GetDownloadUrl(outfit.PhotoURL)

		err = outfitService.AddOutfit(&outfit)
		if err != nil {
			fmt.Println(err)
			return
		}

	})

	fmt.Println("Starting app")
	http.ListenAndServe(":3000", r)

}
