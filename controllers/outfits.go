package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
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
	outfit.PhotoURL = controller.Storage.GetDownloadUrl(outfit.PhotoURL)

	err = controller.OutfitService.AddOutfit(outfit)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
}

func (controller *OutfitsController) Home(w http.ResponseWriter, r *http.Request) {
	uid := ctx.Uid(r.Context())
	user, err := controller.UserService.GetUser(*uid)
	if err != nil {
		// TODO: Check if user not exist or there was an error with firestore
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	last, err := convertToTime(r.FormValue("last"))
	if err != nil {
		http.Error(w, "Please provide valid timestamp", http.StatusBadRequest)
		return
	}

	outfits, err := controller.OutfitService.GetOutfits(user.Followings, last)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	enc := json.NewEncoder(w)
	if err := enc.Encode(&outfits); err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func convertToTime(timestamp string) (time.Time, error) {
	i, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return time.Now(), fmt.Errorf("convert to time | parse int : %w", err)
	}

	return time.Unix(i, 0), nil
}
