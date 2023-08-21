package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AkifhanIlgaz/wit-api/ctx"
	"github.com/AkifhanIlgaz/wit-api/models"
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
