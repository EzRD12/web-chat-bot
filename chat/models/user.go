package models

import (
	"context"
	"errors"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	Id             string `json:"id" bson:"id,omitempty" required:"true"`
	FirstName      string
	LastName       string
	Username       string
	Password       string
	LastConnection string
}

var (
	users  []*User
	nextId = 1
)

func GetUserById(id string, collection *mongo.Collection, ctx context.Context) (User, error) {
	user := User{}
	fmt.Println(id)
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Invalid id")
	}

	err = collection.FindOne(ctx, bson.D{{"id", objectId}}).Decode(&user)

	if err != nil {
		log.Println("User not found")
		return user, errors.New("User not found")
	}

	return user, nil
}

func AddUser(u User, collection *mongo.Collection, ctx context.Context) (User, error) {
	nextId++

	res, insertErr := collection.InsertOne(ctx, u)
	if insertErr != nil {
		log.Fatal(insertErr)
	}
	fmt.Println(res)

	return u, nil
}
