package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/AkifhanIlgaz/wit-api/ctx"
	"github.com/AkifhanIlgaz/wit-api/firebase"
)

type FirebaseController struct {
	Storage *firebase.Storage
}

func (controller *FirebaseController) GenerateUploadUrl(w http.ResponseWriter, r *http.Request) {
	uid := ctx.Uid(r.Context())
	dir, err := storageDirByType(r.Header.Get("type"))
	if err != nil {
		http.Error(w, "Please provide valid type", http.StatusBadRequest)
		return
	}

	fileType := r.Header.Get("fileType")

	timestamp := time.Now().UnixMilli()

	uploadUrl, filePath, err := controller.Storage.GenerateUploadUrl(*uid, timestamp, fileType, dir)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	resp := map[string]string{
		"uploadUrl": uploadUrl,
		"filePath":  filePath,
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(resp)
}

func storageDirByType(t string) (string, error) {
	switch t {
	case "outfit":
		return "outfits", nil
	case "profilePhoto":
		return "profilePhotos", nil
	default:
		return "", fmt.Errorf("storage dir by type: wrong type")
	}
}
