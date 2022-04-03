package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/ezrod12/chat/auth"
	"github.com/ezrod12/chat/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/validator.v2"
)

type authController struct {
	context        context.Context
	userCollection *mongo.Collection
}

func (ac authController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.URL.Path == "/auth" {
		switch r.Method {
		case http.MethodPost:
			ac.post(w, r)
		default:
			w.WriteHeader(http.StatusNotImplemented)
		}
	} else {
		w.WriteHeader(http.StatusNotImplemented)
	}
}

func newAuthController() *authController {
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

	return &authController{
		userCollection: collection,
	}
}

func (ac *authController) post(w http.ResponseWriter, r *http.Request) {
	authRequest, err := ac.parseRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("{\"error\":\"" + err.Error() + "\"}"))
		return
	}

	err = validateAuthRequest(authRequest, true)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{\"error\":\"" + err.Error() + "\"}"))
		return
	}
	auth.AuthUser(w, authRequest, ac.context, ac.userCollection)
}

func (uc *authController) parseRequest(r *http.Request) (models.AuthenthicationRequest, error) {
	dec := json.NewDecoder(r.Body)
	var u models.AuthenthicationRequest
	err := dec.Decode(&u)

	if err != nil {
		return models.AuthenthicationRequest{}, err
	}

	return u, nil
}

func validateAuthRequest(request models.AuthenthicationRequest, createMode bool) error {
	if errs := validator.Validate(request); errs != nil {
		return errs
	}

	return nil
}
