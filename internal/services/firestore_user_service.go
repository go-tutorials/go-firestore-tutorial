package services

import (
	"cloud.google.com/go/firestore"
	"context"
	"google.golang.org/api/iterator"
	"reflect"
	"strings"

	. "go-service/internal/models"
)

type FirestoreUserService struct {
	Collection *firestore.CollectionRef
}

func NewUserService(client *firestore.Client) *FirestoreUserService {
	collection := client.Collection("users")
	return &FirestoreUserService{Collection: collection}
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
	maps := MakeFirestoreMap(userType)

	uid := json["id"]
	id := uid.(string)
	docRef := s.Collection.Doc(id)
	doc, err1 := docRef.Get(ctx)
	if err1 != nil {
		return -1, err1
	}
	delete(json, "id")

	dest := MapToFirestore(json, doc, maps)
	_, err := docRef.Set(ctx, dest)
	if err != nil {
		return -1, err
	}
	return 1, nil
}

func (s *FirestoreUserService) Delete(ctx context.Context, id string) (int64, error) {
	_, err := s.Collection.Doc(id).Delete(ctx)
	if err != nil {
		return -1, err
	}
	return 1, nil
}

func MakeFirestoreMap(modelType reflect.Type) map[string]string {
	maps := make(map[string]string)
	numField := modelType.NumField()
	for i := 0; i < numField; i++ {
		field := modelType.Field(i)
		key1 := field.Name
		if tag0, ok0 := field.Tag.Lookup("json"); ok0 {
			if strings.Contains(tag0, ",") {
				a := strings.Split(tag0, ",")
				key1 = a[0]
			} else {
				key1 = tag0
			}
		}
		if tag, ok := field.Tag.Lookup("firestore"); ok {
			if tag != "-" {
				if strings.Contains(tag, ",") {
					a := strings.Split(tag, ",")
					if key1 == "-" {
						key1 = a[0]
					}
					maps[key1] = a[0]
				} else {
					if key1 == "-" {
						key1 = tag
					}
					maps[key1] = tag
				}
			}
		} else {
			if key1 == "-" {
				key1 = field.Name
			}
			maps[key1] = key1
		}
	}
	return maps
}
func MapToFirestore(json map[string]interface{}, doc *firestore.DocumentSnapshot, maps map[string]string) map[string]interface{} {
	fs := doc.Data()
	for k, v := range json {
		fk, ok := maps[k]
		if ok {
			fs[fk] = v
		}
	}
	return fs
}
