package startup

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/ezrod12/chat/models"
	"github.com/ezrod12/chat/services"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type startup struct {
	context             context.Context
	userCollection      *mongo.Collection
	roomCollection      *mongo.Collection
	usersRoomCollection *mongo.Collection
}

func InitStartup() *startup {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	userCollection := client.Database("chat-bot").Collection("users")
	roomCollection := client.Database("chat-bot").Collection("rooms")
	usersRoomCollection := client.Database("chat-bot").Collection("chat-users")

	return &startup{
		userCollection:      userCollection,
		roomCollection:      roomCollection,
		usersRoomCollection: usersRoomCollection,
	}
}

func (s *startup) InsertRoomSeedData() (models.ChatRoom, error) {
	var room models.ChatRoom
	err := readRoomSeedDataFile(&room)

	if err != nil {
		return models.ChatRoom{}, nil
	}

	var existRoom models.ChatRoom
	existRoom, err = services.GetRoomsByName(room.Name, s.roomCollection, s.context)

	if err != nil {
		existRoom, _ = services.AddRoom(room, s.roomCollection, s.context)
	}

	return existRoom, nil
}

func (s *startup) InsertUserSeedData() (models.User, error) {
	var user models.User
	err := readUserSeedDataFile(&user)

	if err != nil {
		return models.User{}, nil
	}

	var existUser models.User
	existUser, err = services.GetUserByUsername(user.Username, s.userCollection, s.context)

	if err != nil {
		existUser, _ = services.AddUser(user, s.userCollection, s.context)
	}

	return existUser, nil
}

func (s *startup) InsertUsersRoomData(userId, roomId string) {
	var userAlreadyAssigned, _ = services.DoesExistsUserInRoom(roomId, userId, s.usersRoomCollection, s.context)

	if userAlreadyAssigned {
		return
	}

	var userChatRoom models.ChatRoomUsers = models.ChatRoomUsers{UserId: userId, RoomId: roomId, Joined: time.Now()}
	services.AddUserChatRoom(userChatRoom, s.usersRoomCollection, s.context)
}

func readRoomSeedDataFile(data *models.ChatRoom) error {
	f, err := os.Open("seed-data/room.json")

	if err != nil {
		log.Fatal(err)
		return err
	}

	defer f.Close()

	decoder := json.NewDecoder(f)
	err = decoder.Decode(data)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func (s *startup) SaveSeedData() {
	var room models.ChatRoom
	var user models.User

	room, roomError := s.InsertRoomSeedData()
	user, userError := s.InsertUserSeedData()

	if roomError != nil || userError != nil {
		return
	}

	s.InsertUsersRoomData(user.Id, room.Id)
}

func readUserSeedDataFile(data *models.User) error {
	f, err := os.Open("seed-data/users.json")
	defer f.Close()

	if err != nil {
		log.Fatal(err)
		return err
	}

	decoder := json.NewDecoder(f)
	err = decoder.Decode(data)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}
