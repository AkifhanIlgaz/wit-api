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

func (service *UserService) AddUser(uid string, user User) error {
	collection := service.Client.Collection(usersCollection)

	_, err := collection.Doc(uid).Set(context.TODO(), user)
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
		return fmt.Errorf("follow | update following: %w", err)
	}

	// TODO: If an error occurs on second action undo the first action
	_, err = collection.Doc(followedUid).Update(context.TODO(), []firestore.Update{
		{
			Path:  "followers",
			Value: firestore.ArrayUnion(currentUid),
		},
	})
	if err != nil {
		return fmt.Errorf("follow | update followers: %w", err)
	}

	return nil
}

func (service *UserService) Unfollow(currentUid, unfollowedUid string) error {
	collection := service.Client.Collection(usersCollection)

	_, err := collection.Doc(currentUid).Update(context.TODO(), []firestore.Update{
		{
			Path:  "followings",
			Value: firestore.ArrayRemove(unfollowedUid),
		},
	})
	if err != nil {
		return fmt.Errorf("unfollow | update following: %w", err)
	}

	// TODO: If an error occurs on second action undo the first action
	_, err = collection.Doc(unfollowedUid).Update(context.TODO(), []firestore.Update{
		{
			Path:  "followers",
			Value: firestore.ArrayRemove(currentUid),
		},
	})
	if err != nil {
		return fmt.Errorf("unfollow | update followers: %w", err)
	}
	return nil
}
