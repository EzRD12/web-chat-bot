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

	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		m.Id = oid.Hex()
	}

	return m, nil
}

func GetRoomMessages(id string, collection *mongo.Collection, ctx context.Context) ([]*models.Message, error) {
	cur, currErr := collection.Find(ctx, bson.D{{"chatroomid", id}})

	if currErr != nil {
		log.Fatal(currErr)
	}
	defer cur.Close(ctx)

	var messagesCollection []*models.Message
	if currErr := cur.All(ctx, &messagesCollection); currErr != nil {
		log.Fatal(currErr)
	}
	fmt.Println(messagesCollection)
	return messagesCollection, currErr
}
