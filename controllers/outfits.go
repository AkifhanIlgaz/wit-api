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
)

type OutfitsController struct {
	Storage       *firebase.Storage
	OutfitService *models.OutfitService
	UserService   *models.UserService
}

type outfitResponse struct {
	models.Outfit
	IsLiked      bool   `json:"isLiked"`
	LikeCount    int    `json:"likeCount"`
	IsSaved      bool   `json:"isSaved"`
	ProfilePhoto string `json:"profilePhoto"`
	DisplayName  string `json:"displayName"`
}

func (controller *OutfitsController) New(w http.ResponseWriter, r *http.Request) {
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

	err = controller.OutfitService.Add(outfit)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
}

func (controller *OutfitsController) Home(w http.ResponseWriter, r *http.Request) {
	uid := ctx.Uid(r.Context())
	user, err := controller.UserService.Get(*uid)
	if err != nil {
		// TODO: Check if user not exist or there was an error with firestore
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	last := convertToTime(r.URL.Query().Get("last"))

	outfits, err := controller.OutfitService.Home(user.Followings, last)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	var respBody []outfitResponse

	for _, outfit := range outfits {
		var resp outfitResponse

		isLiked, likeCount := controller.OutfitService.LikeStatus(&outfit, *uid)
		resp.Outfit = outfit
		resp.IsLiked = isLiked
		resp.LikeCount = likeCount
		outfitOwner, err := controller.UserService.Get(outfit.Uid)
		if err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
		resp.IsSaved = controller.UserService.IsOutfitSaved(user.Saved, outfit.Id)
		resp.ProfilePhoto = outfitOwner.PhotoUrl
		resp.DisplayName = outfitOwner.DisplayName
		respBody = append(respBody, resp)
	}

	err = writeToResponse(w, respBody)
	if err != nil {
		return
	}

}

func (controller *OutfitsController) Count(w http.ResponseWriter, r *http.Request) {
	uid := r.URL.Query().Get("uid")
	if uid == "" {
		http.Error(w, "User doesn't exist", http.StatusBadRequest)
		return
	}

	count, err := controller.OutfitService.OutfitCountOfUser(uid)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	respBody := struct {
		OutfitCount int `json:"outfitCount"`
	}{
		OutfitCount: count,
	}

	err = writeToResponse(w, respBody)
	if err != nil {
		return
	}
}

func (controller *OutfitsController) All(w http.ResponseWriter, r *http.Request) {
	uid := r.URL.Query().Get("uid")
	if uid == "" {
		http.Error(w, "User doesn't exist", http.StatusBadRequest)
		return
	}
	user, err := controller.UserService.Get(uid)
	if err != nil {
		// TODO: Check if user not exist or there was an error with firestore
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	last := convertToTime(r.URL.Query().Get("last"))
	outfits, err := controller.OutfitService.User(uid, last)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	var respBody []outfitResponse

	for _, outfit := range outfits {
		var resp outfitResponse

		isLiked, likeCount := controller.OutfitService.LikeStatus(&outfit, uid)
		resp.Outfit = outfit
		resp.IsLiked = isLiked
		resp.LikeCount = likeCount
		outfitOwner, err := controller.UserService.Get(outfit.Uid)
		if err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
		resp.IsSaved = controller.UserService.IsOutfitSaved(user.Saved, outfit.Id)
		resp.ProfilePhoto = outfitOwner.PhotoUrl
		resp.DisplayName = outfitOwner.DisplayName
		respBody = append(respBody, resp)
	}

	err = writeToResponse(w, respBody)
	if err != nil {
		return
	}
}

func (controller *OutfitsController) Like(w http.ResponseWriter, r *http.Request) {
	uid := ctx.Uid(r.Context())
	outfitId := r.FormValue("outfitId")

	err := controller.OutfitService.Like(outfitId, *uid)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (controller *OutfitsController) Unlike(w http.ResponseWriter, r *http.Request) {
	uid := ctx.Uid(r.Context())
	outfitId := r.FormValue("outfitId")

	err := controller.OutfitService.Unlike(outfitId, *uid)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (controller *OutfitsController) RemoveLink(w http.ResponseWriter, r *http.Request) {

	outfitId := r.URL.Query().Get("outfitId")

	fmt.Println(outfitId)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Please provide body", http.StatusBadRequest)
		return
	}

	var link models.Link

	err = json.Unmarshal(body, &link)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	err = controller.OutfitService.RemoveLink(outfitId, link)
	fmt.Println(err)

}

func convertToTime(timestamp string) time.Time {
	t, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		t = time.Now()
	}

	return t
}
