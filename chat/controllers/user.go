package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/ezrod12/chat/helpers"
	"github.com/ezrod12/chat/models"
	"github.com/ezrod12/chat/services"
	"gopkg.in/validator.v2"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type userController struct {
	userIdPattern *regexp.Regexp
	context       context.Context
	collection    *mongo.Collection
}

func (uc userController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.URL.Path == "/users" {
		switch r.Method {
		case http.MethodPost:
			uc.post(w, r)
		default:
			w.WriteHeader(http.StatusNotImplemented)
		}
	} else {
		if !helpers.IsAuthorized(r, w) {
			return
		}

		matches := uc.userIdPattern.FindStringSubmatch(r.URL.Path)
		if len(matches) == 0 {
			w.WriteHeader(http.StatusNotFound)
		}
		id := matches[1]

		fmt.Println(id)

		switch r.Method {
		case http.MethodGet:
			fmt.Println(id)
		default:
			w.WriteHeader(http.StatusNotImplemented)
		}
	}
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
		w.Write([]byte("{\"error\":\"" + err.Error() + "\"}"))
		return
	}

	err = uc.validateUserEntity(u, true)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{\"error\":\"" + err.Error() + "\"}"))
		return
	}

	u, err = services.AddUser(u, uc.collection, uc.context)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("{\"error\":\"" + err.Error() + "\"}"))
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

func (uc userController) validateUserEntity(user models.User, createMode bool) error {
	if errs := validator.Validate(user); errs != nil {
		return errs
	}

	_, err := services.GetUserByUsername(user.Username, uc.collection, uc.context)

	if err == nil {
		return errors.New("user with this username already exists")
	}

	return nil
}
