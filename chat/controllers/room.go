package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/ezrod12/chat/models"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type roomController struct {
	roomIdPattern           *regexp.Regexp
	context                 context.Context
	collection              *mongo.Collection
	userCollection          *mongo.Collection
	chatRoomCollection      *mongo.Collection
	chatRoomUsersCollection *mongo.Collection
}

func (uc roomController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/rooms" {
		switch r.Method {
		case http.MethodGet:
			uc.getAll(w, r)
		case http.MethodPost:
			uc.post(w, r)
		default:
			w.WriteHeader(http.StatusNotImplemented)
		}
	} else {
		matches := uc.roomIdPattern.FindStringSubmatch(r.URL.Path)
		if len(matches) == 0 {
			w.WriteHeader(http.StatusNotFound)
		}
		id := matches[1]

		fmt.Println(id)

		switch r.Method {
		case http.MethodGet:
			uc.get(id, w)
		default:
			w.WriteHeader(http.StatusNotImplemented)
		}
	}
}

func (uc *roomController) getAll(w http.ResponseWriter, r *http.Request) {
	encodeResponseAsJson(models.GetRooms(uc.collection, uc.context), w)
}

func (uc *roomController) get(id string, w http.ResponseWriter) {
	u, err := models.GetRoomMessages(id, uc.collection, uc.context)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	encodeResponseAsJson(u, w)
}

func newRoomController() *roomController {
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
	messageCollection := client.Database("chat-bot").Collection("messages")
	chatRoomCollection := client.Database("chat-bot").Collection("rooms")

	return &roomController{
		roomIdPattern:      regexp.MustCompile(`/rooms/([A-Za-z0-9\-]+)/?`),
		collection:         messageCollection,
		userCollection:     userCollection,
		chatRoomCollection: chatRoomCollection,
	}
}

func (uc *roomController) post(w http.ResponseWriter, r *http.Request) {
	u, err := uc.parseRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	err = uc.validateMessageEntity(u, true)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	u, err = models.AddMessage(u, uc.collection, uc.context)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	encodeResponseAsJson(u, w)
}

func (uc *roomController) parseRequest(r *http.Request) (models.Message, error) {
	dec := json.NewDecoder(r.Body)
	var m models.Message
	err := dec.Decode(&m)

	if err != nil {
		return models.Message{}, err
	}

	m.Created = time.Now()

	return m, nil
}

func (mc *roomController) validateMessageEntity(message models.Message, createMode bool) error {
	if strings.Trim(message.Value, " ") == "" {
		return errors.New("property firstname must contain a valid string value")
	}

	_, err := models.GetUserById(message.SenderUserId, mc.userCollection, mc.context)

	if err != nil {
		return err
	}

	_, err = models.GetChatRoomById(message.ChatRoomId, mc.chatRoomCollection, mc.context)

	if err != nil {
		return err
	}

	return nil
}
