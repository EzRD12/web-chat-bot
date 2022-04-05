package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/ezrod12/chat/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func AddMessage(m models.Message, collection *mongo.Collection, ctx context.Context) (models.Message, error) {
	mapMutex := sync.RWMutex{}

	mapMutex.Lock()
	m.Created = time.Now()
	message := models.Message{}

	currErr := collection.FindOne(ctx, bson.M{"chatRoomId": m.ChatRoomId, "username": m.Username, "value": m.Value, "created": m.Created}).Decode(&message)
	if currErr == nil {
		return m, errors.New("message already exists")
	}

	res, insertErr := collection.InsertOne(ctx, m)
	if insertErr != nil {
		log.Fatal(insertErr)
	}

	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		m.Id = oid.Hex()
	}
	mapMutex.Unlock()

	return m, nil
}

func GetRoomMessages(id string, collection *mongo.Collection, ctx context.Context) ([]*models.Message, error) {
	opts := options.Find().SetLimit(50).SetSort(bson.D{{"created", -1}})
	cur, currErr := collection.Find(ctx, bson.D{{"chatRoomId", id}}, opts)

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
