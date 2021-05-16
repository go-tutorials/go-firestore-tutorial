package services

import (
	"cloud.google.com/go/firestore"
	"context"
	fs "github.com/core-go/firestore"
	"google.golang.org/api/iterator"
	"reflect"

	. "go-service/internal/models"
)

type FirestoreUserService struct {
	Client *firestore.Client
	Collection *firestore.CollectionRef
}

func NewUserService(client *firestore.Client) *FirestoreUserService {
	collection := client.Collection("users")
	return &FirestoreUserService{Collection: collection, Client: client}
}

func (s *FirestoreUserService) GetAll(ctx context.Context) (*[]User, error) {
	iter := s.Collection.Documents(ctx)
	var users []User
	for {
		doc, er1 := iter.Next()
		if er1 == iterator.Done {
			break
		}
		if er1 != nil {
			return nil, er1
		}
		var user User
		er2 := doc.DataTo(&user)
		if er2 != nil {
			return &users, er2
		}

		user.Id = doc.Ref.ID
		user.CreateTime = &doc.CreateTime
		user.UpdateTime = &doc.UpdateTime
		users = append(users, user)
	}
	return &users, nil
}

func (s *FirestoreUserService) Load(ctx context.Context, id string) (*User, error) {
	doc, er1 := s.Collection.Doc(id).Get(ctx)
	var user User
	if er1 != nil {
		return nil, er1
	}
	er2 := doc.DataTo(&user)
	if er2 == nil {
		user.Id = id
		user.CreateTime = &doc.CreateTime
		user.UpdateTime = &doc.UpdateTime
	}
	return &user, er2
}

func (s *FirestoreUserService) Insert(ctx context.Context, user *User) (int64, error) {
	_, err := s.Collection.Doc(user.Id).Set(ctx, user)
	if err != nil {
		return -1, err
	}
	return 1, nil
}

func (s *FirestoreUserService) Update(ctx context.Context, user *User) (int64, error) {
	_, err := s.Collection.Doc(user.Id).Set(ctx, user)
	if err != nil {
		return -1, err
	}
	return 1, nil
}

func (s *FirestoreUserService) Patch(ctx context.Context, json map[string]interface{}) (int64, error) {
	userType := reflect.TypeOf(User{})
	maps := fs.MakeFirestoreMap(userType)

	uid := json["id"]
	id := uid.(string)
	docRef := s.Collection.Doc(id)
	doc, err1 := docRef.Get(ctx)
	if err1 != nil {
		return -1, err1
	}
	delete(json, "id")

	dest := fs.MapToFirestore(json, doc, maps)
	_, err := docRef.Set(ctx, dest)
	if err != nil {
		return -1, err
	}
	return 1, nil
}

func (s *FirestoreUserService) Delete(ctx context.Context, id string) (int64, error) {
	_, err := s.Client.Collection("users").Doc(id).Delete(ctx)
	if err != nil {
		return -1, err
	}
	return 1, nil
}
