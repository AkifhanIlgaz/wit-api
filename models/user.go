package models

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"firebase.google.com/go/v4/auth"
)

const usersCollection = "users"

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
	Auth   *auth.Client
}

func (service *UserService) AddUser(user User) error {
	collection := service.Client.Collection(usersCollection)

	doc := collection.NewDoc()
	_, err := doc.Set(context.TODO(), user)
	if err != nil {
		return fmt.Errorf("add user: %w", err)
	}

	return nil
}

func (service *UserService) UpdateUser(displayName string, photoUrl string) error {
	// ? Merge
	panic("Implement")
}

func (service *UserService) Follow(currentUid, followedUid string) error {
	collection := service.Client.Collection(usersCollection)

	_, err := collection.Doc(currentUid).Update(context.TODO(), []firestore.Update{
		{
			Path:  "followings",
			Value: firestore.ArrayUnion(followedUid),
		},
	})
	if err != nil {
		return fmt.Errorf("update following: %w", err)
	}

	_, err = collection.Doc(followedUid).Update(context.TODO(), []firestore.Update{
		{
			Path:  "followers",
			Value: firestore.ArrayUnion(currentUid),
		},
	})
	if err != nil {
		return fmt.Errorf("update followers: %w", err)
	}

	return nil
}
