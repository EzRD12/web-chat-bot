package main

import (
	"log"
	"net/http"

	"github.com/ezrod12/chat/controllers"
	"github.com/ezrod12/chat/messager"
	"github.com/ezrod12/chat/settings"
)

func main() {
	config := settings.GetConfig()
	amqp, err := messager.Connect(config)

	if err != nil {
		return
	}
	defer amqp.Close()
	ch, err := messager.OpenChannel()
	if err != nil {
		return
	}
	defer ch.Close()

	controllers.RegisterController()
	http.ListenAndServe(":3000", nil)

	log.Printf(" [*] Service started. To exit press CTRL+C")
}
