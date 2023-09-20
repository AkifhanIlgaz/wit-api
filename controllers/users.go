package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/AkifhanIlgaz/wit-api/ctx"
	"github.com/AkifhanIlgaz/wit-api/firebase"
	"github.com/AkifhanIlgaz/wit-api/models"
)

type UsersController struct {
	UserService   *models.UserService
	OutfitService *models.OutfitService
}

func (controller *UsersController) NewUser(w http.ResponseWriter, r *http.Request) {
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

	err = controller.UserService.AddUser(*uid, user)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
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

func (controller *UsersController) SaveOutfit(w http.ResponseWriter, r *http.Request) {
	uid := ctx.Uid(r.Context())
	outfitId := r.FormValue("outfitId")

	fmt.Println(outfitId)
	err := controller.UserService.SaveOutfit(outfitId, *uid)
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
	lastFollower := r.URL.Query().Get("lastFollower")

	followers, err := controller.UserService.GetFollowers(uid, lastFollower)
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
}

func (controller *UsersController) Followings(w http.ResponseWriter, r *http.Request) {
	uid := r.URL.Query().Get("uid")
	lastFollowing := r.URL.Query().Get("lastFollowing")

	followings, err := controller.UserService.GetFollowings(uid, lastFollowing)
	if err != nil {
		fmt.Println(err)
		return
	}

	enc := json.NewEncoder(w)
	err = enc.Encode(&followings)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

}

func (controller *UsersController) GetUser(w http.ResponseWriter, r *http.Request) {
	uid := r.URL.Query().Get("uid")
	lastOutfit := convertToTime(r.URL.Query().Get("lastOutfit"))
	lastSaved := convertToTime(r.URL.Query().Get("lastSaved"))
	lastFollower := r.URL.Query().Get("lastFollower")
	lastFollowing := r.URL.Query().Get("lastFollowing")

	u, err := controller.UserService.GetUser(uid)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	type response struct {
		DisplayName string          `json:"displayName"`
		PhotoUrl    string          `json:"photoUrl"`
		Outfits     []models.Outfit `json:"outfits"`
		Saved       []models.Outfit `json:"saved"`
		Followers   []models.User   `json:"followers"`
		Followings  []models.User   `json:"followings"`
	}

	res := response{
		DisplayName: u.DisplayName,
		PhotoUrl:    u.PhotoUrl,
	}

	outfits, err := controller.OutfitService.GetUserOutfits(uid, lastOutfit)
	if err != nil {
		// ?
	}
	res.Outfits = outfits

	saved, err := controller.OutfitService.GetOutfits(u.Saved, lastSaved)
	if err != nil {
		// ?
	}
	res.Saved = saved

	followers, err := controller.UserService.GetFollowers(uid, lastFollower)
	if err != nil {
		// ?
	}
	res.Followers = followers

	followings, err := controller.UserService.GetFollowings(uid, lastFollowing)
	if err != nil {
		// ?
	}
	res.Followings = followings

	encoder := json.NewEncoder(w)
	err = encoder.Encode(&res)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
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
