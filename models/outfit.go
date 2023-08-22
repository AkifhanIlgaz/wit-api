package models

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
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

func (service *OutfitService) GetAllOutfitsByUid(uid string) ([]Outfit, error) {
	outfits := []Outfit{}

	iter := service.Collection.Where("Uid", "==", uid).Documents(context.TODO())
	for {
		var outfit Outfit
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("get outfits by uid: %w", err)
		}

		err = doc.DataTo(&outfit)
		if err != nil {
			return nil, fmt.Errorf("get outfits by uid | data to: %w", err)
		}

		outfits = append(outfits, outfit)
	}

	return outfits, nil
}

func (service *OutfitService) GetOutfitById(outfitId string) (Outfit, error) {
	panic("Implement this function")
}
