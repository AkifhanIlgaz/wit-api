package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/AkifhanIlgaz/wit-api/controllers"
	"github.com/AkifhanIlgaz/wit-api/ctx"
	"github.com/AkifhanIlgaz/wit-api/firebase"
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

	outfitsController := controllers.OutfitsController{
		OutfitService: *myApp.FirestoreService.OutfitService,
	}

	uidMiddleware := controllers.UidMiddleware{
		AuthService: myApp.AuthService,
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

	r.Get("/generate-upload-url", func(w http.ResponseWriter, r *http.Request) {
		uid := ctx.Uid(r.Context())
		fileExtension := r.Header.Get("fileExtension")

		uploadUrl, err := myApp.StorageService.GenerateUploadUrl(*uid, fileExtension)
		if err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, uploadUrl)
	})

	r.Route("/outfits", func(r chi.Router) {
		r.Get("/all/{uid}", outfitsController.GetAllOutfitsByUid)
		r.Get("/{outfitId}", outfitsController.GetOutfitById)
		r.Delete("/{outfitId}", outfitsController.DeleteOutfit)
		r.Post("/add", outfitsController.Add)
	})

	fmt.Println("Starting app")
	http.ListenAndServe(":3000", r)

}
