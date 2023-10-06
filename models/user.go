package models

import (
	"context"
	"fmt"
	"strings"
	"time"

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
	Followers   []string `firestore:"followers" json:"-"` // Store followers by uid
	Followings  []string `firestore:"followings" json:"-"`
	Saved       []string `firestore:"saved" json:"-"`
	IsFollowed  bool     `firestore:"-" json:"isFollowed"`
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

func (service *UserService) UpdateUser(uid, displayName, photoUrl string) error {
	collection := service.Client.Collection(usersCollection)

	var updates []firestore.Update

	if displayName != "" {
		updates = append(updates, firestore.Update{
			Path:  "displayName",
			Value: displayName})
	}
	if photoUrl != "" {
		updates = append(updates, firestore.Update{
			Path:  "photoUrl",
			Value: photoUrl})
	}

	_, err := collection.Doc(uid).Update(context.TODO(), updates)
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}

	return nil
}

func (service *UserService) Filter(filterString string) ([]User, error) {
	collection := service.Client.Collection(usersCollection)

	refs := collection.DocumentRefs(context.TODO())

	var users []User

	for ref, err := refs.Next(); err == nil; ref, err = refs.Next() {
		snapshot, err := ref.Get(context.TODO())
		if err != nil {
			fmt.Println(err)
			continue
		}

		displayName, err := snapshot.DataAt("displayName")

		// TODO: Sort by point
	}

	return nil, nil
}

func (service *UserService) isMatchFilter(displayName string, filters []string) bool {
	// TODO: Increment the point for every match.

	for _, filter := range filters {
		if strings.Contains(displayName, filter) {
			return true
		}
	}

	return false
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

func (service *UserService) SaveOutfit(outfitId, uid string, savedAt time.Time) error {
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
func (service *UserService) GetFollowers(uid, last string) ([]User, error) {
	collection := service.Client.Collection(usersCollection)

	snapshot, _ := collection.Doc(uid).Get(context.TODO())
	var user User
	err := snapshot.DataTo(&user)
	if err != nil {
		return nil, fmt.Errorf("get followings | data to: %w", err)
	}

	var filter firestore.OrFilter
	for _, uid := range user.Followers {
		doc, err := collection.Doc(uid).Get(context.TODO())
		if err != nil {
			fmt.Errorf("get followers | get by uid: %w", err)
			break
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
		var u User
		snapshot.DataTo(&u)
		u.Uid = snapshot.Ref.ID
		u.IsFollowed = service.IsFollowed(user.Followings, u.Uid)
		followers = append(followers, u)
	}

	return followers, nil
}

// TODO: Sort and limit
func (service *UserService) GetFollowings(uid, last string) ([]User, error) {
	collection := service.Client.Collection(usersCollection)

	snapshot, _ := collection.Doc(uid).Get(context.TODO())
	var user User
	err := snapshot.DataTo(&user)
	if err != nil {
		return nil, fmt.Errorf("get followings | data to: %w", err)
	}

	var filter firestore.OrFilter
	for _, uid := range user.Followings {
		doc, err := collection.Doc(uid).Get(context.TODO())
		if err != nil {
			break
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
		return nil, fmt.Errorf("get followings | query: %w", err)
	}

	var followings []User
	for _, snapshot := range userSnapshots {
		var user User
		snapshot.DataTo(&user)
		user.Uid = snapshot.Ref.ID
		user.IsFollowed = true
		followings = append(followings, user)
	}

	return followings, nil
}

func (service *UserService) IsFollowed(users []string, uid string) bool {
	return slices.Contains[[]string, string](users, uid)
}

func (service *UserService) IsOutfitSaved(saved []string, outfitId string) bool {
	return slices.Contains[[]string, string](saved, outfitId)
}
