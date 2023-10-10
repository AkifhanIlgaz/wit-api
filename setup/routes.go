package setup

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func Routes(controllers *controllers) *chi.Mux {
	r := chi.NewMux()

	r.Use(controllers.UidMiddleware.SetUid)
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

	r.Get("/generate-upload-url", controllers.FirebaseController.GenerateUploadUrl)

	r.Route("/outfit", func(r chi.Router) {
		r.Post("/new", controllers.OutfitsController.New)

		r.Get("/home", controllers.OutfitsController.Home)

		r.Put("/like", controllers.OutfitsController.Like)
		r.Put("/unlike", controllers.OutfitsController.Unlike)

		r.Get("/all", controllers.OutfitsController.All)

		r.Get("/count", controllers.OutfitsController.Count)

		r.Delete("/links", controllers.OutfitsController.RemoveLink)
	})

	r.Route("/user", func(r chi.Router) {
		r.Post("/new", controllers.UsersController.New)
		r.Get("/", controllers.UsersController.User)
		r.Put("/update", controllers.UsersController.Update)

		r.Put("/follow", controllers.UsersController.Follow)
		r.Put("/unfollow", controllers.UsersController.Unfollow)

		r.Get("/saved", controllers.UsersController.Saved)
		r.Put("/save-outfit", controllers.UsersController.SaveOutfit)
		r.Put("/unsave-outfit", controllers.UsersController.UnsaveOutfit)

		r.Get("/followers", controllers.UsersController.Followers)
		r.Get("/followings", controllers.UsersController.Followings)

		r.Get("/filter", controllers.UsersController.Filter)
	})

	return r
}
