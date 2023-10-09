package models

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/firestore/apiv1/firestorepb"
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

func (service *OutfitService) Add(outfit Outfit) error {
	collection := service.Client.Collection(outfitCollection)

	doc := collection.NewDoc()
	_, err := doc.Set(context.TODO(), outfit)
	if err != nil {
		return fmt.Errorf("add document: %w", err)
	}

	return nil
}

func (service *OutfitService) Home(uids []string, last time.Time) ([]Outfit, error) {
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

func (service *OutfitService) User(uid string, last time.Time) ([]Outfit, error) {
	collection := service.Client.Collection(outfitCollection)

	snapshots, err := collection.Where("uid", "==", uid).OrderBy("createdAt", firestore.Desc).StartAfter(last).Limit(postNumbersPerRequest).Documents(context.TODO()).GetAll()
	if err != nil {
		return nil, fmt.Errorf("get outfits by uid: %w", err)
	}

	var outfits []Outfit
	for _, snapshot := range snapshots {
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

func (service *OutfitService) Outfit(outfitId string) (*Outfit, error) {
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

func (service *OutfitService) Outfits(outfitIds []string) ([]Outfit, error) {
	collection := service.Client.Collection(outfitCollection)
	var outfits []Outfit

	for _, id := range outfitIds {
		var outfit Outfit

		snapshot, err := collection.Doc(id).Get(context.TODO())
		if err != nil {
			fmt.Println(err)
			continue
		}
		outfit.Id = snapshot.Ref.ID
		err = snapshot.DataTo(&outfit)
		if err != nil {
			return nil, fmt.Errorf("get outfits | data to: %w", err)
		}

		outfits = append(outfits, outfit)

	}

	return outfits, nil
}

func (service *OutfitService) OutfitCountOfUser(uid string) (int, error) {
	collection := service.Client.Collection(outfitCollection)

	query := collection.Where("uid", "==", uid)
	aggregationQuery := query.NewAggregationQuery().WithCount("all")
	results, err := aggregationQuery.Get(context.Background())
	if err != nil {
		return 0, fmt.Errorf("get outfit count of user: %w", err)
	}

	count, ok := results["all"]
	if !ok {
		return 0, fmt.Errorf("get outfit count of user: %w", err)
	}
	countValue := count.(*firestorepb.Value)

	return int(countValue.GetIntegerValue()), nil
}

func (service *OutfitService) LikeStatus(outfit *Outfit, uid string) (bool, int) {
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

func (service *OutfitService) RemoveLink(outfitId string, link Link) error {
	collection := service.Client.Collection(outfitCollection)

	_, err := collection.Doc(outfitId).Update(context.TODO(), []firestore.Update{{
		Path:  "links",
		Value: firestore.ArrayRemove(link),
	}})
	if err != nil {
		return fmt.Errorf("remove link: %w", err)
	}

	return nil
}
