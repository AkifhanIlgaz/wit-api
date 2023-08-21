package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/AkifhanIlgaz/wit-api/firebase"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	credentialsFile := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")

	myApp, err := firebase.NewApp(credentialsFile)
	if err != nil {
		panic(err)
	}

	// links := []firebase.Link{
	// 	{Title: "first link",
	// 		Href:     "first link href",
	// 		Position: firebase.Position{Left: "50%", Top: "50%"},
	// 	},
	// 	{Title: "second link",
	// 		Href:     "second link href",
	// 		Position: firebase.Position{Left: "25%", Top: "25%"},
	// 	},
	// 	{Title: "third link",
	// 		Href:     "third link href",
	// 		Position: firebase.Position{Left: "75%", Top: "75%"},
	// 	},
	// }

	// outfit := firebase.Outfit{
	// 	Uid:      "user id",
	// 	PhotoURL: "photo url",
	// 	Links:    links,
	// }

	r := chi.NewRouter()

	r.Post("/add-outfit", func(w http.ResponseWriter, r *http.Request) {
		// TODO: Embed uid inside of request's context by middleware
		idToken := r.Header.Get("idToken")
		uid, err := myApp.AuthService.GetUidByIdToken(idToken)
		if err != nil {
			fmt.Fprintln(w, err)
		}
		fmt.Println(uid)

		photoURL := r.FormValue("photoURL")
		fmt.Println(photoURL)
	})

	fmt.Println("Starting app")
	http.ListenAndServe(":3000", r)

}
