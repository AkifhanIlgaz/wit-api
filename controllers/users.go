package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/AkifhanIlgaz/wit-api/ctx"
	"github.com/AkifhanIlgaz/wit-api/firebase"
	"github.com/AkifhanIlgaz/wit-api/models"
	"golang.org/x/exp/slices"
)

type UsersController struct {
	UserService   *models.UserService
	OutfitService *models.OutfitService
	Auth          *firebase.Auth
	Storage       *firebase.Storage
}

func (controller *UsersController) User(w http.ResponseWriter, r *http.Request) {
	uid := ctx.Uid(r.Context())
	user, err := controller.UserService.Get(r.URL.Query().Get("uid"))
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	user.IsFollowed = controller.UserService.IsFollowed(user.Followers, *uid)

	err = writeToResponse(w, user)
	if err != nil {
		return
	}

}

func (controller *UsersController) New(w http.ResponseWriter, r *http.Request) {
	uid := ctx.Uid(r.Context())

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Please provide body", http.StatusBadRequest)
		return
	}

	var user models.User
	err = json.Unmarshal(body, &user)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	err = controller.UserService.Add(*uid, user)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
}

func (controller *UsersController) Update(w http.ResponseWriter, r *http.Request) {
	uid := ctx.Uid(r.Context())

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Please provide body", http.StatusBadRequest)
		return
	}

	var res struct {
		PhotoUrl    string `json:"photoUrl"`
		DisplayName string `json:"displayName"`
	}

	err = json.Unmarshal(body, &res)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	res.PhotoUrl = controller.Storage.GetDownloadUrl(res.PhotoUrl)

	err = controller.Auth.UpdateUser(*uid, res.PhotoUrl, res.DisplayName)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	err = controller.UserService.Update(*uid, res.DisplayName, res.PhotoUrl)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

}

func (controller *UsersController) Filter(w http.ResponseWriter, r *http.Request) {
	filterString := r.URL.Query().Get("filterString")

	users, err := controller.UserService.Filter(filterString)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	err = writeToResponse(w, users)
	if err != nil {
		return
	}

}

func (controller *UsersController) Follow(w http.ResponseWriter, r *http.Request) {
	currentUid := ctx.Uid(r.Context())
	followedUid := r.FormValue("uid")

	err := controller.UserService.Follow(*currentUid, followedUid)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (controller *UsersController) Unfollow(w http.ResponseWriter, r *http.Request) {
	currentUid := ctx.Uid(r.Context())
	unfollowedUid := r.FormValue("uid")

	err := controller.UserService.Unfollow(*currentUid, unfollowedUid)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (controller *UsersController) Saved(w http.ResponseWriter, r *http.Request) {
	uid := ctx.Uid(r.Context())
	user, err := controller.UserService.Get(*uid)

	var outfitIds []string
	last := r.URL.Query().Get("last")
	if last == "" {
		if len(user.Saved) >= 5 {
			outfitIds = user.Saved[:5]
		} else {
			outfitIds = user.Saved[:]
		}
	} else {
		index := slices.Index(user.Saved, last)
		if index+5 < len(user.Saved) {
			outfitIds = user.Saved[index+1 : index+6]
		} else {
			outfitIds = user.Saved[index+1:]
		}
	}
	saved, err := controller.OutfitService.Outfits(outfitIds)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	err = writeToResponse(w, saved)
	if err != nil {
		return
	}

}

func (controller *UsersController) SaveOutfit(w http.ResponseWriter, r *http.Request) {
	uid := ctx.Uid(r.Context())
	outfitId := r.FormValue("outfitId")

	err := controller.UserService.SaveOutfit(outfitId, *uid, time.Now())
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (controller *UsersController) UnsaveOutfit(w http.ResponseWriter, r *http.Request) {
	uid := ctx.Uid(r.Context())
	outfitId := r.FormValue("outfitId")

	err := controller.UserService.UnsaveOutfit(outfitId, *uid)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (controller *UsersController) Followers(w http.ResponseWriter, r *http.Request) {
	uid := r.URL.Query().Get("uid")
	lastFollower := r.URL.Query().Get("last")

	followers, err := controller.UserService.Followers(uid, lastFollower)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = writeToResponse(w, followers)
	if err != nil {
		return
	}
}

func (controller *UsersController) Followings(w http.ResponseWriter, r *http.Request) {
	uid := r.URL.Query().Get("uid")
	lastFollowing := r.URL.Query().Get("lastFollowing")

	followings, err := controller.UserService.Followings(uid, lastFollowing)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = writeToResponse(w, followings)
	if err != nil {
		return
	}

}

type UidMiddleware struct {
	Auth *firebase.Auth
}

func (umw UidMiddleware) SetUid(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idToken := r.Header.Get("Authorization")

		uid, err := umw.Auth.GetUidByIdToken(idToken)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		r = r.WithContext(ctx.WithUid(r.Context(), uid))

		next.ServeHTTP(w, r)
	})
}
