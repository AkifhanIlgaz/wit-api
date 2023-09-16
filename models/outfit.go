package models

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"golang.org/x/exp/slices"
)

const outfitCollection = "outfits"

// TODO: Change it to 20
const postNumbersPerRequest = 3

type Outfit struct {
	Id        string    `firestore:"-" json:"id"`
	Uid       string    `firestore:"uid" json:"uid"`
	PhotoUrl  string    `firestore:"photoUrl" json:"photoUrl"`
	Links     []Link    `firestore:"links" json:"links"`
	Likes     []string  `firestore:"likes" json:"-"`
	CreatedAt time.Time `firestore:"createdAt" json:"createdAt"`
}

type Link struct {
	Title string `firestore:"title" json:"title"`
	Href  string `firestore:"href" json:"href"`
	Left  string `firestore:"left" json:"left"`
	Top   string `firestore:"top" json:"top"`
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

func (service *OutfitService) GetHomeOutfits(uids []string, last time.Time) ([]Outfit, error) {
	collection := service.Client.Collection(outfitCollection)

	var filter firestore.OrFilter
	for _, uid := range uids {
		filter.Filters = append(filter.Filters, firestore.PropertyFilter{
			Path:     "uid",
			Operator: "==",
			Value:    uid,
		})
	}

	outfitSnapshots, err := collection.WhereEntity(filter).OrderBy("createdAt", firestore.Desc).StartAfter(last).Limit(postNumbersPerRequest).Documents(context.TODO()).GetAll()
	if err != nil {
		return nil, fmt.Errorf("get outfits | query: %w", err)
	}

	var outfits []Outfit
	for _, snapshot := range outfitSnapshots {
		var outfit Outfit

		outfit.Id = snapshot.Ref.ID
		err := snapshot.DataTo(&outfit)
		if err != nil {
			return nil, fmt.Errorf("get outfits | data to: %w", err)
		}

		outfits = append(outfits, outfit)
	}

	return outfits, nil
}

func (service *OutfitService) GetOutfit(outfitId string) (*Outfit, error) {
	var outfit Outfit

	collection := service.Client.Collection(outfitCollection)
	snapshot, err := collection.Doc(outfitId).Get(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("get outfit: %w", err)
	}

	snapshot.DataTo(&outfit)
	outfit.Id = snapshot.Ref.ID

	return &outfit, nil
}

func (service *OutfitService) GetLikeStatus(outfit *Outfit, uid string) (bool, int) {
	return slices.Contains[[]string, string](outfit.Likes, uid), len(outfit.Likes)
}

func (service *OutfitService) Like(outfitId, uid string) error {
	collection := service.Client.Collection(outfitCollection)

	_, err := collection.Doc(outfitId).Update(context.TODO(), []firestore.Update{{
		Path:  "likes",
		Value: firestore.ArrayUnion(uid),
	}})
	if err != nil {
		return fmt.Errorf("like: %w", err)
	}

	return nil
}

func (service *OutfitService) Unlike(outfitId, uid string) error {
	collection := service.Client.Collection(outfitCollection)

	_, err := collection.Doc(outfitId).Update(context.TODO(), []firestore.Update{{
		Path:  "likes",
		Value: firestore.ArrayRemove(uid),
	}})
	if err != nil {
		return fmt.Errorf("unlike: %w", err)
	}

	return nil
}

// *********** OLD ****************

// func (service *OutfitService) DeleteOutfit(uid, outfitId string) error {
// 	collection := service.Client.Collection(outfitCollection)
// 	var outfit Outfit
// 	// ! If current user is not the owner of this outfit return error
// 	doc := collection.Doc(outfitId)
// 	snapshot, err := doc.Get(context.TODO())
// 	if err != nil {
// 		return fmt.Errorf("delete outfit | get: %w", err)
// 	}
// 	err = snapshot.DataTo(&outfit)
// 	if err != nil {
// 		return fmt.Errorf("delete outfit | data to: %w", err)
// 	}

// 	if outfit.Uid != uid {
// 		return fmt.Errorf("current user is not the owner of this outfit")
// 	}

// 	_, err = doc.Delete(context.TODO())
// 	if err != nil {
// 		return fmt.Errorf("delete outfit: %w", err)
// 	}

// 	return nil
// }

// // TODO: Update outfit

// func (service *OutfitService) GetAllOutfitsByUid(uid string) ([]Outfit, error) {
// 	collection := service.Client.Collection(outfitCollection)
// 	outfits := []Outfit{}

// 	iter := collection.Where("Uid", "==", uid).Documents(context.TODO())
// 	for {
// 		var outfit Outfit
// 		doc, err := iter.Next()
// 		if err == iterator.Done {
// 			break
// 		}
// 		if err != nil {
// 			return nil, fmt.Errorf("get outfits by uid: %w", err)
// 		}

// 		err = doc.DataTo(&outfit)
// 		if err != nil {
// 			return nil, fmt.Errorf("get outfits by uid | data to: %w", err)
// 		}

// 		outfits = append(outfits, outfit)
// 	}

// 	return outfits, nil
// }

// func (service *OutfitService) GetOutfitById(outfitId string) (*Outfit, error) {
// 	collection := service.Client.Collection(outfitCollection)
// 	var outfit Outfit

// 	snapshot, err := collection.Doc(outfitId).Get(context.TODO())
// 	if err != nil {
// 		return nil, fmt.Errorf("get outfit by id: %w", err)
// 	}

// 	err = snapshot.DataTo(&outfit)
// 	if err != nil {
// 		return nil, fmt.Errorf("get outfit by id | data to : %w", err)
// 	}

// 	return &outfit, nil
// }
