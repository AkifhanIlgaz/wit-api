package models

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"firebase.google.com/go/v4/auth"
	"golang.org/x/exp/slices"
)

const usersCollection = "users"

// User stored on Firestore
type User struct {
	Uid         string   `firestore:"-" json:"uid"`
	DisplayName string   `firestore:"displayName" json:"displayName"`
	PhotoUrl    string   `firestore:"photoUrl" json:"photoUrl"`
	Followers   []string `firestore:"followers" json:"followers"` // Store followers by uid
	Followings  []string `firestore:"followings" json:"followings"`
	Saved       []string `firestore:"saved" json:"saved"`
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

func (service *UserService) GetUser(uid string) (*User, error) {
	var user User
	collection := service.Client.Collection(usersCollection)

	snapshot, err := collection.Doc(uid).Get(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	snapshot.DataTo(&user)

	return &user, nil
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

func (service *UserService) SaveOutfit(outfitId, uid string) error {
	collection := service.Client.Collection(usersCollection)

	_, err := collection.Doc(uid).Update(context.TODO(), []firestore.Update{
		{
			Path:  "saved",
			Value: firestore.ArrayUnion(outfitId),
		},
	})
	if err != nil {
		return fmt.Errorf("save outfit: %w", err)
	}

	return nil
}

func (service *UserService) UnsaveOutfit(outfitId, uid string) error {
	collection := service.Client.Collection(usersCollection)

	_, err := collection.Doc(uid).Update(context.TODO(), []firestore.Update{
		{
			Path:  "saved",
			Value: firestore.ArrayRemove(outfitId),
		},
	})
	if err != nil {
		return fmt.Errorf("unsave outfit: %w", err)
	}

	return nil
}

// TODO: Sort and limit
func (service *UserService) GetFollowers(uid string, last string) ([]User, error) {
	collection := service.Client.Collection(usersCollection)
	ref, _ := collection.Doc(uid).Get(context.TODO())

	uids, err := ref.DataAt("followers")
	if err != nil {
		return nil, fmt.Errorf("get followers | data at: %w", err)
	}

	var filter firestore.OrFilter
	for _, uid := range uids.([]interface{}) {
		doc, err := collection.Doc(uid.(string)).Get(context.TODO())
		if err != nil {
			fmt.Errorf("get followers | get by uid: %w", err)
		}
		// ! check error
		displayName, _ := doc.DataAt("displayName")

		filter.Filters = append(filter.Filters, firestore.PropertyFilter{
			Path:     "displayName",
			Operator: "==",
			Value:    displayName,
		})
	}

	userSnapshots, err := collection.WhereEntity(filter).OrderBy("displayName", firestore.Asc).StartAfter(last).Limit(5).Documents(context.TODO()).GetAll()
	if err != nil {
		return nil, fmt.Errorf("get outfits | query: %w", err)
	}

	var followers []User
	for _, snapshot := range userSnapshots {
		var user User
		snapshot.DataTo(&user)
		user.Uid = snapshot.Ref.ID
		followers = append(followers, user)
	}

	return followers, nil
}

// TODO: Sort and limit
func (service *UserService) GetFollowings(uid string) ([]User, error) {
	user, err := service.GetUser(uid)
	if err != nil {
		return nil, fmt.Errorf("get followings: %w", err)
	}

	var followings []User
	for _, uid := range user.Followings {
		following, _ := service.GetUser(uid)
		followings = append(followings, *following)
	}

	return followings, nil
}

func (service *UserService) IsFollowed(followers []string, uid string) bool {
	return slices.Contains[[]string, string](followers, uid)
}

func (service *UserService) IsOutfitSaved(saved []string, outfitId string) bool {
	return slices.Contains[[]string, string](saved, outfitId)
}
