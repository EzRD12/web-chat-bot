package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

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

	err = collection.FindOne(ctx, bson.D{{"_id", objectId}}).Decode(&chat)

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

func GetRoomsByIds(roomsId []string, collection *mongo.Collection, ctx context.Context) []*models.ChatRoom {
	oids := make([]primitive.ObjectID, len(roomsId))
	for i := range roomsId {
		oids[i], _ = primitive.ObjectIDFromHex(roomsId[i])
	}

	query := bson.M{"_id": bson.M{"$in": oids}}
	cur, currErr := collection.Find(ctx, query)

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

func GetRoomsIdAssignedToUser(userId string, collection *mongo.Collection, ctx context.Context) []string {
	fmt.Println(userId)
	cur, currErr := collection.Find(ctx, bson.M{"userid": userId})

	if currErr != nil {
		log.Fatal(currErr)
		fmt.Println(currErr)
	}
	defer cur.Close(ctx)

	var roomsUserCollection []*models.ChatRoomUsers
	if err := cur.All(ctx, &roomsUserCollection); err != nil {
		panic(err)
	}

	roomsId := make([]string, len(roomsUserCollection))
	for i, roomUser := range roomsUserCollection {
		roomsId[i] = roomUser.RoomId
	}

	return roomsId
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
	u.Created = time.Now()
	res, insertErr := collection.InsertOne(ctx, u)
	if insertErr != nil {
		log.Fatal(insertErr)
	}

	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		u.Id = oid.Hex()
	}

	return u, nil
}

func AddUserChatRoom(u models.ChatRoomUsers, collection *mongo.Collection, ctx context.Context) (models.ChatRoomUsers, error) {
	u.Joined = time.Now()
	_, insertErr := collection.InsertOne(ctx, u)
	if insertErr != nil {
		log.Fatal(insertErr)
	}

	return u, nil
}
