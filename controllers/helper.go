package controllers

import (
	"encoding/json"
	"net/http"
)

func writeToResponse(w http.ResponseWriter, data any) error {
	enc := json.NewEncoder(w)
	if err := enc.Encode(&data); err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return err
	}
	return nil
}
