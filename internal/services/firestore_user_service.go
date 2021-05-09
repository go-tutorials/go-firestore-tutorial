package services

import (
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	"github.com/mitchellh/mapstructure"
	"google.golang.org/api/iterator"

	. "go-service/internal/models"
)

type FirestoreUserService struct {
	Client *firestore.Client
}

func NewUserService(client *firestore.Client) *FirestoreUserService {
	return &FirestoreUserService{Client: client}
}

func (s *FirestoreUserService) GetAll(ctx context.Context) (*[]User, error) {
	iter := s.Client.Collection("users").Documents(ctx)
	var result []User
	for {
		doc, er1 := iter.Next()
		if er1 == iterator.Done {
			break
		}
		if er1 != nil {
			return nil, er1
		}
		// convert map to json
		jsonString, _ := json.Marshal(doc.Data())
		// convert json to struct
		s := User{}
		er2 := json.Unmarshal(jsonString, &s)
		if er2 == nil {
			s.Id = doc.Ref.ID
		}
		result = append(result, s)
	}
	return &result, nil
}

func (s *FirestoreUserService) Load(ctx context.Context, id string) (*User, error) {
	dsnap, er1 := s.Client.Collection("users").Doc(id).Get(ctx)
	var user *User
	if er1 != nil {
		return nil, er1
	}
	er2 := mapstructure.Decode(dsnap.Data(), &user)
	if er2 == nil {
		user.Id = id
	}
	return user, er1
}

func (s *FirestoreUserService) Insert(ctx context.Context, user *User) (int64, error) {
	_, err := s.Client.Collection("users").Doc(user.Id).Set(ctx, user)
	if err != nil {
		return -1, err
	}
	return 1, nil
}

func (s *FirestoreUserService) Update(ctx context.Context, user *User) (int64, error) {
	_, err := s.Client.Collection("users").Doc(user.Id).Set(ctx, user)
	if err != nil {
		return -1, err
	}
	return 1, nil
}

func (s *FirestoreUserService) Patch(ctx context.Context, user map[string]interface{}) (int64, error) {
	//get element by id
	id := user["id"]
	sid := id.(string)
	dsnap, err1 := s.Client.Collection("users").Doc(sid).Get(ctx)
	if err1 != nil {
		return -1, err1
	}

	dest := make(map[string]interface{})
	for k, v := range dsnap.Data() {
		dest[k] = v
	}

	for k, v := range user {
		dest[k] = v
	}
	_, err := s.Client.Collection("users").Doc(sid).Set(ctx, dest)
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
