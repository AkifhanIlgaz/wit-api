package models

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

const outfitCollection = "outfits"

type Outfit struct {
	Id       string `firestore:"id"`
	Uid      string `firestore:"uid"`
	PhotoURL string `firestore:"photoURL"`
	Links    []Link `firestore:"links"`
}

type Link struct {
	Title    string   `firestore:"title"`
	Href     string   `firestore:"href"`
	Position Position `firestore:"position"`
}

type Position struct {
	Left string `firestore:"left"`
	Top  string `firestore:"top"`
}

type OutfitService struct {
	Collection *firestore.CollectionRef
}

func (service *OutfitService) AddOutfit(outfit *Outfit) error {
	doc := service.Collection.NewDoc()
	outfit.Id = doc.ID
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

func (service *OutfitService) GetOutfitById(outfitId string) (*Outfit, error) {
	var outfit Outfit

	snapshot, err := service.Collection.Doc(outfitId).Get(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("get outfit by id: %w", err)
	}

	err = snapshot.DataTo(&outfit)
	if err != nil {
		return nil, fmt.Errorf("get outfit by id | data to : %w", err)
	}

	return &outfit, nil
}
