package models

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/storage"
	"github.com/google/uuid"
	"google.golang.org/api/iterator"
)

const outfitCollection = "outfits"

type Outfit struct {
	Id        string    `firestore:"id"`
	Uid       string    `firestore:"uid"`
	PhotoURL  string    `firestore:"photoURL"`
	Links     []Link    `firestore:"links"`
	Likes     []string  `firestore:"likes"`
	CreatedAt time.Time `firestore:"createdAt"`
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
	Bucket     *storage.BucketHandle
}

func (service *OutfitService) UploadPhoto(imageBase64 string) string {
	id := uuid.New()

	fmt.Println(service.Bucket)

	object := service.Bucket.Object("outfits/second")
	writer := object.NewWriter(context.TODO())

	writer.ObjectAttrs.Metadata = map[string]string{"firebaseStorageDownloadTokens": id.String()}
	defer writer.Close()
	if _, err := io.Copy(writer, bytes.NewReader([]byte(imageBase64))); err != nil {
		fmt.Println(err)
		return ""
	}

	return "Successfully uploaded"
}

func (service *OutfitService) AddOutfit(outfit *Outfit) error {
	doc := service.Collection.NewDoc()
	outfit.Id = doc.ID
	outfit.CreatedAt = time.Now()
	_, err := doc.Set(context.TODO(), outfit)
	if err != nil {
		return fmt.Errorf("add document: %w", err)
	}

	return nil
}

func (service *OutfitService) DeleteOutfit(uid, outfitId string) error {
	var outfit Outfit
	// ! If current user is not the owner of this outfit return error
	doc := service.Collection.Doc(outfitId)
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
