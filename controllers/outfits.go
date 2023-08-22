package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AkifhanIlgaz/wit-api/ctx"
	"github.com/AkifhanIlgaz/wit-api/models"
	"github.com/go-chi/chi/v5"
)

type OutfitsController struct {
	OutfitService models.OutfitService
}

func (controller *OutfitsController) Add(w http.ResponseWriter, r *http.Request) {
	var outfit models.Outfit

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&outfit)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	outfit.Uid = *ctx.Uid(r.Context())

	err = controller.OutfitService.AddOutfit(&outfit)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(&outfit)
}

func (controller *OutfitsController) GetByUid(w http.ResponseWriter, r *http.Request) {
	uid := chi.URLParam(r, "uid")
	if uid == "" {
		http.Error(w, "Please give an valid uid", http.StatusBadRequest)
		return
	}
	fmt.Println(uid)

	outfits, err := controller.OutfitService.GetAllOutfitsByUid(uid)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(&outfits)
}
