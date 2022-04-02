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

	"github.com/ezrod12/chat/models"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type userController struct {
	userIdPattern *regexp.Regexp
	context       context.Context
	collection    *mongo.Collection
}

func (uc userController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
}

func newUserController() *userController {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	collection := client.Database("chat-bot").Collection("users")

	return &userController{
		userIdPattern: regexp.MustCompile(`/users/([A-Za-z0-9\-]+)/?`),
		collection:    collection,
	}
}

func (uc *userController) post(w http.ResponseWriter, r *http.Request) {
	u, err := uc.parseRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	err = validateUserEntity(u, true)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	u, err = models.AddUser(u, uc.collection, uc.context)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	encodeResponseAsJson(u, w)
}

func (uc *userController) parseRequest(r *http.Request) (models.User, error) {
	dec := json.NewDecoder(r.Body)
	var u models.User
	err := dec.Decode(&u)

	if err != nil {
		return models.User{}, err
	}

	return u, nil
}

func validateUserEntity(user models.User, createMode bool) error {
	if strings.Trim(user.FirstName, " ") == "" {
		return errors.New("property firstname must contain a valid string value")
	}

	if strings.Trim(user.LastName, " ") == "" {
		return errors.New("property lastname must contain a valid string value")
	}

	return nil
}
