package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/AkifhanIlgaz/wit-api/ctx"
	"github.com/AkifhanIlgaz/wit-api/firebase"
)

type FirebaseController struct {
	Storage *firebase.Storage
}

func (fc *FirebaseController) GenerateUploadUrl(w http.ResponseWriter, r *http.Request) {
	uid := ctx.Uid(r.Context())
	fileExtension := r.Header.Get("fileExtension")
	timestamp := time.Now().UnixMilli()

	uploadUrl, filePath, err := fc.Storage.GenerateUploadUrl(*uid, timestamp, fileExtension)
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
