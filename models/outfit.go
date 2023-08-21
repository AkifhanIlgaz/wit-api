package models

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
)

const outfitCollection = "outfits"

type Outfit struct {
	Uid      string `json:"uid"`
	PhotoURL string `json:"photoURL"`
	Links    []Link `json:"links"`
}

type Link struct {
	Title    string   `json:"title"`
	Href     string   `json:"href"`
	Position Position `json:"position"`
}

type Position struct {
	Left string `json:"left"`
	Top  string `json:"top"`
}

type OutfitService struct {
	Collection *firestore.CollectionRef
}

func (service *OutfitService) AddOutfit(outfit *Outfit) error {
	doc := service.Collection.NewDoc()
	_, err := doc.Set(context.TODO(), outfit)
	if err != nil {
		return fmt.Errorf("add document: %w", err)
	}

	return nil
}
