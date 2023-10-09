package setup

import (
	"github.com/AkifhanIlgaz/wit-api/controllers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func Routes(uidMiddleware *controllers.UidMiddleware, firebaseController *controllers.FirebaseController, outfitsController *controllers.OutfitsController, usersController *controllers.UsersController) *chi.Mux {
	r := chi.NewMux()

	r.Use(uidMiddleware.SetUid)
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*", "Accept", "Authorization", "Content-Type", "X-CSRF-Token", "idToken", "fileExtension", "type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Get("/generate-upload-url", firebaseController.GenerateUploadUrl)

	r.Route("/outfit", func(r chi.Router) {
		r.Post("/new", outfitsController.New)
		
		r.Get("/home", outfitsController.Home)

		r.Put("/like", outfitsController.Like)
		r.Put("/unlike", outfitsController.Unlike)

		r.Get("/all", outfitsController.All)

		r.Get("/count", outfitsController.Count)

		r.Delete("/links", outfitsController.RemoveLink)
	})

	r.Route("/user", func(r chi.Router) {
		r.Post("/new", usersController.New)
		r.Get("/", usersController.User)
		r.Put("/update", usersController.Update)

		r.Put("/follow", usersController.Follow)
		r.Put("/unfollow", usersController.Unfollow)

		r.Get("/saved", usersController.Saved)
		r.Put("/save-outfit", usersController.SaveOutfit)
		r.Put("/unsave-outfit", usersController.UnsaveOutfit)

		r.Get("/followers", usersController.Followers)
		r.Get("/followings", usersController.Followings)

		r.Get("/filter", usersController.Filter)
	})

	return r
}
