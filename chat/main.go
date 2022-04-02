package main

import (
	"log"
	"net/http"

	"github.com/ezrod12/chat/controllers"
	"github.com/ezrod12/chat/settings"
)

func main() {
	settings.GetConfig()
	controllers.RegisterController()
	http.ListenAndServe(":3000", nil)

	log.Printf(" [*] Service started. To exit press CTRL+C")
}
