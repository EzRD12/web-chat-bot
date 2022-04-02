package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/ezrod12/chat/helpers"
	"github.com/ezrod12/chat/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetUserById(id string, collection *mongo.Collection, ctx context.Context) (models.User, error) {
	user := models.User{}
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

func GetUserByUsername(username string, collection *mongo.Collection, ctx context.Context) (models.User, error) {
	user := models.User{}

	err := collection.FindOne(ctx, bson.D{{"username", username}}).Decode(&user)

	if err != nil {
		log.Println("User not found with username: " + username)
		return user, errors.New("Invalid user")
	}

	return user, nil
}

func AddUser(u models.User, collection *mongo.Collection, ctx context.Context) (models.User, error) {
	u.Password = helpers.GetHash([]byte(u.Password))

	res, insertErr := collection.InsertOne(ctx, u)
	if insertErr != nil {
		log.Fatal(insertErr)
	}
	fmt.Println(res)

	return u, nil
}

func UpdateLastConnection(u models.User, collection *mongo.Collection, ctx context.Context) (models.User, error) {
	fmt.Println(u)
	u.LastConnection = time.Now()
	objectId, err := primitive.ObjectIDFromHex(u.Id)

	if err != nil {
		log.Println("Invalid id")
	}

	_, insertErr := collection.UpdateOne(ctx, bson.M{"_id": objectId},
		bson.D{
			{"$set", bson.D{{"lastconnection", time.Now()}}},
		})
	if insertErr != nil {
		log.Fatal(insertErr)
	}

	return u, nil
}
