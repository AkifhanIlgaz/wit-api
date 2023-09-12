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
