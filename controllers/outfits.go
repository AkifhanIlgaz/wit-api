package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/AkifhanIlgaz/wit-api/ctx"
	"github.com/AkifhanIlgaz/wit-api/firebase"
	"github.com/AkifhanIlgaz/wit-api/models"
)

type OutfitsController struct {
	Storage       *firebase.Storage
	OutfitService *models.OutfitService
	UserService   *models.UserService
}

func (controller *OutfitsController) NewOutfit(w http.ResponseWriter, r *http.Request) {
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
	outfit.PhotoUrl = controller.Storage.GetDownloadUrl(outfit.PhotoUrl)

	err = controller.OutfitService.AddOutfit(outfit)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
}

func (controller *OutfitsController) Home(w http.ResponseWriter, r *http.Request) {
	type response struct {
		models.Outfit
		isLiked      bool ``
		likeCount    int
		ProfilePhoto string `json:"profilePhoto"`
		DisplayName  string `json:"displayName"`
	}

	uid := ctx.Uid(r.Context())
	user, err := controller.UserService.GetUser(*uid)
	if err != nil {
		// TODO: Check if user not exist or there was an error with firestore
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	var last time.Time

	last = convertToTime(r.URL.Query().Get("last"))

	outfits, err := controller.OutfitService.GetOutfits(user.Followings, last)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	var respBody []response

	for _, outfit := range outfits {
		var resp response

		isLiked, likeCount := controller.OutfitService.GetLikeStatus(&outfit, *uid)
		resp.Outfit = outfit
		resp.isLiked = isLiked
		resp.likeCount = likeCount
		outfitOwner, err := controller.UserService.GetUser(outfit.Uid)
		if err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
		resp.ProfilePhoto = outfitOwner.PhotoUrl
		resp.DisplayName = outfitOwner.DisplayName
		respBody = append(respBody, resp)
	}

	enc := json.NewEncoder(w)
	if err := enc.Encode(&respBody); err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

}

func convertToTime(timestamp string) time.Time {
	t, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		t = time.Now()
	}

	return t
}
