package models

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ChatRoom struct {
	Id      string `json:"id" bson:"id,omitempty"`
	Name    string `max:"65" required:"true"`
	Created time.Time
}

type ChatRoomUsers struct {
	RoomId string `json:"roomId" bson:"roomId"`
	UserId string
	Joined time.Time
}

func GetChatRoomById(id string, collection *mongo.Collection, ctx context.Context) (ChatRoom, error) {
	chat := ChatRoom{}
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

func GetRooms(collection *mongo.Collection, ctx context.Context) []*ChatRoom {
	cur, currErr := collection.Find(ctx, bson.D{})

	if currErr != nil {
		panic(currErr)
	}
	defer cur.Close(ctx)

	var roomsCollection []*ChatRoom
	if err := cur.All(ctx, &roomsCollection); err != nil {
		panic(err)
	}
	fmt.Println(roomsCollection)
	return roomsCollection
}

func AddRoom(u ChatRoom, collection *mongo.Collection, ctx context.Context) (ChatRoom, error) {
	res, insertErr := collection.InsertOne(ctx, u)
	if insertErr != nil {
		log.Fatal(insertErr)
	}
	fmt.Println(res)

	return u, nil
}
