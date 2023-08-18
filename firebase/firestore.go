package firebase

import (
	"cloud.google.com/go/firestore"
)

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

func (service *FirestoreService) AddDocument(collection string, data any) error {
	// ref, _, err := service.Client.Collection(collection).Add(context.TODO(), data)
	// if err != nil {
	// 	return fmt.Errorf("add document: %w", err)
	// }

	// snapshot, _ := ref.Get(context.TODO())

	return nil
}
