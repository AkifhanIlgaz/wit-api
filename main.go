package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/AkifhanIlgaz/wit-api/controllers"
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

	outfitService := &models.OutfitService{
		Client: myApp.Firestore.Client,
	}

	userService := &models.UserService{
		Client: myApp.Firestore.Client,
	}

	uidMiddleware := controllers.UidMiddleware{
		Auth: myApp.Auth,
	}

	firebaseController := controllers.FirebaseController{
		Storage: myApp.Storage,
	}

	outfitsController := controllers.OutfitsController{
		Storage:       myApp.Storage,
		OutfitService: outfitService,
		UserService:   userService,
	}

	usersController := controllers.UsersController{
		UserService:   userService,
		OutfitService: outfitService,
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

	r.Post("/outfit/new", outfitsController.NewOutfit)
	r.Get("/outfit/home", outfitsController.Home)
	r.Put("/outfit/like", outfitsController.Like)
	r.Put("/outfit/unlike", outfitsController.Unlike)

	r.Post("/user/new", usersController.NewUser)
	r.Put("/user/follow", usersController.Follow)
	r.Put("/user/unfollow", usersController.Unfollow)
	r.Put("/user/save-outfit", usersController.SaveOutfit)
	r.Put("/user/unsave-outfit", usersController.UnsaveOutfit)
	r.Get("/user/followers", func(w http.ResponseWriter, r *http.Request) {
		uid := r.URL.Query().Get("uid")
		// user, err := usersController.UserService.GetUser(uid)
		// if err != nil {
		// 	http.Error(w, "Something went wrong", http.StatusInternalServerError)
		// 	return
		// }

		type response struct {
			Uid         string `json:"uid"`
			PhotoUrl    string `json:"photoUrl"`
			DisplayName string `json:"displayName"`
			IsFollowed  bool   `json:"isFollowed"`
			Count       int    `json:"count"`
		}

		followers, err := usersController.UserService.GetFollowers(uid, "Ula Maty")
		if err != nil {
			fmt.Println(err)
			return
		}

		enc := json.NewEncoder(w)
		err = enc.Encode(&followers)
		if err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
	})
	r.Get("/user/followings", usersController.Followings)

	fmt.Println("Starting app")
	http.ListenAndServe(":3000", r)

}
