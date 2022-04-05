package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/ezrod12/chat/auth"
	"github.com/ezrod12/chat/helpers"
	"github.com/ezrod12/chat/models"
	"github.com/ezrod12/chat/services"
	"gopkg.in/validator.v2"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type roomController struct {
	roomIdPattern           *regexp.Regexp
	roomMessagesPattern     *regexp.Regexp
	stockMessagePattern     *regexp.Regexp
	context                 context.Context
	collection              *mongo.Collection
	userCollection          *mongo.Collection
	chatRoomCollection      *mongo.Collection
	chatRoomUsersCollection *mongo.Collection
}

func (uc roomController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if !helpers.IsAuthorized(r, w) {
		return
	}

	if r.URL.Path == "/rooms" {
		switch r.Method {
		case http.MethodGet:
			uc.getAll(w, r)
		case http.MethodPost:
			uc.createRoom(w, r)
		default:
			w.WriteHeader(http.StatusNotImplemented)
		}
	} else {
		matches := uc.roomIdPattern.FindStringSubmatch(r.URL.Path)
		roomMessages := uc.roomMessagesPattern.FindStringSubmatch(r.URL.Path)
		if len(matches) == 0 {
			w.WriteHeader(http.StatusNotFound)
		}
		id := matches[1]

		if len(roomMessages) > 0 {
			switch r.Method {
			case http.MethodPost:
				uc.createMessage(w, r)
			case http.MethodGet:
				uc.getRoomMessages(id, w)
			default:
				w.WriteHeader(http.StatusNotImplemented)
			}
		} else {
			if http.MethodGet == r.Method {
				uc.get(id, w)
			}
		}
		w.WriteHeader(http.StatusNotImplemented)
	}
}

func (uc *roomController) getAll(w http.ResponseWriter, r *http.Request) {
	claims, _ := auth.ExtractClaims(r.Header.Get("Authorization"))
	value := claims["userId"]
	roomsId := services.GetRoomsIdAssignedToUser(value.(string), uc.chatRoomUsersCollection, uc.context)
	encodeResponseAsJson(services.GetRoomsByIds(roomsId, uc.chatRoomCollection, uc.context), w)
}

func (uc *roomController) get(id string, w http.ResponseWriter) {
	u, err := services.GetChatRoomDetailById(id, uc.chatRoomCollection, uc.context)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	encodeResponseAsJson(u, w)
}

func (uc *roomController) getRoomMessages(id string, w http.ResponseWriter) {
	u, err := services.GetRoomMessages(id, uc.collection, uc.context)
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
	chatRoomUsersCollection := client.Database("chat-bot").Collection("chat-users")

	return &roomController{
		roomIdPattern:           regexp.MustCompile(`/rooms/([A-Za-z0-9\-]+)/?`),
		roomMessagesPattern:     regexp.MustCompile(`/rooms/([A-Za-z0-9\-]+)/messages?`),
		stockMessagePattern:     regexp.MustCompile(`/stock=([A-Za-z0-9\-]+)`),
		collection:              messageCollection,
		userCollection:          userCollection,
		chatRoomCollection:      chatRoomCollection,
		chatRoomUsersCollection: chatRoomUsersCollection,
	}
}

func (uc *roomController) createRoom(w http.ResponseWriter, r *http.Request) {
	u, err := uc.parseRoomRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("{\"error\":\"" + err.Error() + "\"}"))
		return
	}

	err = uc.validateRoomEntity(u, true)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{\"error\":\"" + err.Error() + "\"}"))
		return
	}

	var existRoom models.ChatRoom
	existRoom, err = services.GetRoomsByName(u.Name, uc.chatRoomCollection, uc.context)

	if err != nil {
		existRoom, err = services.AddRoom(u, uc.chatRoomCollection, uc.context)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{\"error\":\"Room with name" + existRoom.Name + " already exists\"}"))
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("{\"error\":\"" + err.Error() + "\"}"))
	}

	encodeResponseAsJson(existRoom, w)
}

func (uc *roomController) createMessage(w http.ResponseWriter, r *http.Request) {
	msg, err := uc.parseRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("{\"error\":\"" + err.Error() + "\"}"))
		return
	}

	claims, _ := auth.ExtractClaims(r.Header.Get("Authorization"))
	value := claims["user"]
	msg.Username = value.(string)
	msg.Created = time.Now()

	err = uc.validateMessageEntity(msg, true)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{\"error\":\"" + err.Error() + "\"}"))
		return
	}

	msg, err = services.AddMessage(msg, uc.collection, uc.context)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("{\"error\":\"" + err.Error() + "\"}"))
	}

	encodeResponseAsJson(msg, w)
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

func (uc *roomController) parseRoomRequest(r *http.Request) (models.ChatRoom, error) {
	dec := json.NewDecoder(r.Body)
	var m models.ChatRoom
	err := dec.Decode(&m)

	if err != nil {
		return models.ChatRoom{}, err
	}

	m.Created = time.Now()

	return m, nil
}

func (mc *roomController) validateMessageEntity(message models.Message, createMode bool) error {
	if strings.Trim(message.Value, " ") == "" {
		return errors.New("message value must contain a valid string value")
	}

	_, err := services.GetChatRoomDetailById(message.ChatRoomId, mc.chatRoomCollection, mc.context)

	if err != nil {
		return err
	}

	return nil
}

func (mc *roomController) validateRoomEntity(room models.ChatRoom, createMode bool) error {
	if errs := validator.Validate(room); errs != nil {
		return errs
	}

	return nil
}
