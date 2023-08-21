package firebase

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

type FirestoreService struct {
	Client *firestore.Client
}

func (service *FirestoreService) AddOutfit(outfit *Outfit) error {
	doc := service.Client.Collection(outfitCollection).NewDoc()
	_, err := doc.Set(context.TODO(), outfit)
	if err != nil {
		return fmt.Errorf("add document: %w", err)
	}

	fmt.Println(outfit)
	return nil
}
