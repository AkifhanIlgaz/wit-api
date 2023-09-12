package models

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
)

// User stored on Firestore
type User struct {
	DisplayName string   `firestore:"displayName"`
	PhotoUrl    string   `firestore:"photoUrl"`
	Followers   []string `firestore:"followers"` // Store followers by uid
	Followings  []string `firestore:"followings"`
	Saved       []string `firestore:"saved"`
}

type UserService struct {
	Client *firestore.Client
}

func (service *UserService) AddUser(user User) error {
	collection := service.Client.Collection(outfitCollection)

	doc := collection.NewDoc()
	_, err := doc.Set(context.TODO(), user)
	if err != nil {
		return fmt.Errorf("add user: %w", err)
	}

	return nil
}

// TODO: Update user
