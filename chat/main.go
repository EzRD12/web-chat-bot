package main

import (
	"log"
	"net/http"

	"github.com/ezrod12/chat/controllers"
	"github.com/ezrod12/chat/messager"
	"github.com/ezrod12/chat/pkg"
	"github.com/ezrod12/chat/settings"
)

func main() {
	pkg.RoomsMessages = make(map[string][]string)
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

	s := pkg.NewServer()
	go s.Run()

	go controllers.RegisterController(s)

	http.ListenAndServe(":8080", nil)

	log.Printf(" [*] Service started. To exit press CTRL+C")
}
