package models

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

const outfitCollection = "outfits"

// TODO: Change it to 20
const postNumbersPerRequest = 2

type Outfit struct {
	Uid       string    `firestore:"uid"`
	PhotoURL  string    `firestore:"photoURL"`
	Links     []Link    `firestore:"links"`
	Likes     []string  `firestore:"likes"`
	CreatedAt time.Time `firestore:"createdAt"`
}

type Link struct {
	Title string `firestore:"title"`
	Href  string `firestore:"href"`
	Left  string `firestore:"left"`
	Top   string `firestore:"top"`
}

type OutfitService struct {
	Client *firestore.Client
}

func (service *OutfitService) AddOutfit(outfit Outfit) error {
	collection := service.Client.Collection(outfitCollection)

	doc := collection.NewDoc()
	_, err := doc.Set(context.TODO(), outfit)
	if err != nil {
		return fmt.Errorf("add document: %w", err)
	}

	return nil
}

// If lastOutfitTimestamp is not given use time.Now()
func (service *OutfitService) GetOutfits(uids []string, lastOutfitTimestamp time.Time) ([]Outfit, error) {
	collection := service.Client.Collection(outfitCollection)

	var filter firestore.OrFilter
	for _, uid := range uids {
		filter.Filters = append(filter.Filters, firestore.PropertyFilter{
			Path:     "uid",
			Operator: "==",
			Value:    uid,
		})
	}

	outfitSnapshots, err := collection.WhereEntity(filter).OrderBy("createdAt", firestore.Desc).StartAfter(lastOutfitTimestamp).Limit(postNumbersPerRequest).Documents(context.TODO()).GetAll()
	if err != nil {
		return nil, fmt.Errorf("get outfits | query: %w", err)
	}

	var outfits []Outfit
	for _, snapshot := range outfitSnapshots {
		var outfit Outfit

		err := snapshot.DataTo(&outfit)
		if err != nil {
			return nil, fmt.Errorf("get outfits | data to: %w", err)
		}

		outfits = append(outfits, outfit)
	}

	return outfits, nil
}

// *********** OLD ****************

func (service *OutfitService) DeleteOutfit(uid, outfitId string) error {
	collection := service.Client.Collection(outfitCollection)
	var outfit Outfit
	// ! If current user is not the owner of this outfit return error
	doc := collection.Doc(outfitId)
	snapshot, err := doc.Get(context.TODO())
	if err != nil {
		return fmt.Errorf("delete outfit | get: %w", err)
	}
	err = snapshot.DataTo(&outfit)
	if err != nil {
		return fmt.Errorf("delete outfit | data to: %w", err)
	}

	if outfit.Uid != uid {
		return fmt.Errorf("current user is not the owner of this outfit")
	}

	_, err = doc.Delete(context.TODO())
	if err != nil {
		return fmt.Errorf("delete outfit: %w", err)
	}

	return nil
}

// TODO: Update outfit

func (service *OutfitService) GetAllOutfitsByUid(uid string) ([]Outfit, error) {
	collection := service.Client.Collection(outfitCollection)
	outfits := []Outfit{}

	iter := collection.Where("Uid", "==", uid).Documents(context.TODO())
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
	collection := service.Client.Collection(outfitCollection)
	var outfit Outfit

	snapshot, err := collection.Doc(outfitId).Get(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("get outfit by id: %w", err)
	}

	err = snapshot.DataTo(&outfit)
	if err != nil {
		return nil, fmt.Errorf("get outfit by id | data to : %w", err)
	}

	return &outfit, nil
}
