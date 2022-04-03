package services

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/ezrod12/chat/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetChatRoomDetailById(id string, collection *mongo.Collection, ctx context.Context) (models.ChatRoomDetails, error) {
	chat := models.ChatRoomDetails{}
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Invalid id")
	}

	err = collection.FindOne(ctx, bson.D{{"id", objectId}}).Decode(&chat)

	if err != nil {
		log.Println("Chat room not found")
		return chat, errors.New("Chat room not found")
	}

	return chat, nil
}

func GetUsersCountInRoomByRoomId(id string, collection *mongo.Collection, ctx context.Context) (int64, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Invalid id")
	}
	var count int64
	count, err = collection.CountDocuments(ctx, bson.D{{"roomId", objectId}})

	if err != nil {
		log.Println("Chat room not found")
		return 0, errors.New("Chat room not found")
	}

	return count, nil
}

func GetRooms(collection *mongo.Collection, ctx context.Context) []*models.ChatRoom {
	cur, currErr := collection.Find(ctx, bson.D{})

	if currErr != nil {
		panic(currErr)
	}
	defer cur.Close(ctx)

	var roomsCollection []*models.ChatRoom
	if err := cur.All(ctx, &roomsCollection); err != nil {
		panic(err)
	}
	return roomsCollection
}

func GetRoomsByName(name string, collection *mongo.Collection, ctx context.Context) (models.ChatRoom, error) {
	room := models.ChatRoom{}

	err := collection.FindOne(ctx, bson.D{{"name", name}}).Decode(&room)

	if err != nil {
		fmt.Println(err)
		return models.ChatRoom{}, errors.New("room not found")
	}

	return room, nil
}

func DoesExistsUserInRoom(roomId, userId string, collection *mongo.Collection, ctx context.Context) (bool, error) {

	count, err := collection.CountDocuments(ctx, bson.M{"roomId": roomId, "userid": userId})
	if err != nil {
		return false, errors.New("user not assigned to the room")
	}

	return count > 0, nil
}

func AddRoom(u models.ChatRoom, collection *mongo.Collection, ctx context.Context) (models.ChatRoom, error) {
	res, insertErr := collection.InsertOne(ctx, u)
	if insertErr != nil {
		log.Fatal(insertErr)
	}
	fmt.Println(res)

	return u, nil
}

func AddUserChatRoom(u models.ChatRoomUsers, collection *mongo.Collection, ctx context.Context) (models.ChatRoomUsers, error) {
	res, insertErr := collection.InsertOne(ctx, u)
	if insertErr != nil {
		log.Fatal(insertErr)
	}
	fmt.Println(res)

	return u, nil
}
