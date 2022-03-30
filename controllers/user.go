package controllers

import (
	"context"
	"log"
	"net/http"
	"regexp"
	"time"

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
