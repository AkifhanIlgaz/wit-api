package firebase

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
)

const outfitCollection = "outfits"

type Outfit struct {
	Id       string
	Uid      string
	PhotoUrl string
	Links    []Link
}

type Link struct {
	Name string
	Href string
}

type FirestoreService struct {
	Client *firestore.Client
}

func (service *FirestoreService) AddOutfit(outfit Outfit) error {
	ref, _, err := service.Client.Collection(outfitCollection).Add(context.TODO(), outfit)
	if err != nil {
		return fmt.Errorf("add document: %w", err)
	}

	snapshot, _ := ref.Get(context.TODO())
	fmt.Println(snapshot.Ref.ID)
	return nil
}
