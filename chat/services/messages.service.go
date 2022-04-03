package services

import (
	"context"
	"fmt"
	"log"

	"github.com/ezrod12/chat/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddMessage(m models.Message, collection *mongo.Collection, ctx context.Context) (models.Message, error) {
	res, insertErr := collection.InsertOne(ctx, m)
	if insertErr != nil {
		log.Fatal(insertErr)
	}
	fmt.Println(res)

	return m, nil
}

func GetRoomMessages(id string, collection *mongo.Collection, ctx context.Context) ([]*models.Message, error) {
	objectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		log.Println("Invalid id")
	}

	cur, currErr := collection.Find(ctx, bson.D{{"chatRoomId", objectId}})

	if currErr != nil {
		panic(currErr)
	}
	defer cur.Close(ctx)

	var messagesCollection []*models.Message
	if err = cur.All(ctx, &messagesCollection); err != nil {
		panic(err)
	}
	fmt.Println(messagesCollection)
	return messagesCollection, err
}
