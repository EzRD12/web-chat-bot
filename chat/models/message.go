package models

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Message struct {
	Id           string    `json:"id" bson:"id,omitempty" required:"true"`
	Value        string    `required:"true" max:"160"`
	ChatRoomId   string    `required:"true"`
	SenderUserId string    `required:"true"`
	Created      time.Time `required:"true"`
}

func AddMessage(m Message, collection *mongo.Collection, ctx context.Context) (Message, error) {
	res, insertErr := collection.InsertOne(ctx, m)
	if insertErr != nil {
		log.Fatal(insertErr)
	}
	fmt.Println(res)

	return m, nil
}

func GetRoomMessages(id string, collection *mongo.Collection, ctx context.Context) ([]*Message, error) {
	objectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		log.Println("Invalid id")
	}

	cur, currErr := collection.Find(ctx, bson.D{{"chatRoomId", objectId}})

	if currErr != nil {
		panic(currErr)
	}
	defer cur.Close(ctx)

	var messagesCollection []*Message
	if err = cur.All(ctx, &messagesCollection); err != nil {
		panic(err)
	}
	fmt.Println(messagesCollection)
	return messagesCollection, err
}
