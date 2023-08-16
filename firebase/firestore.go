package firebase

import "cloud.google.com/go/firestore"

type FirestoreService struct {
	Client *firestore.Client
}
