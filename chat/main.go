package main

import (
	"net/http"

	"github.com/ezrod12/chat/controllers"
)

func main() {
	controllers.RegisterController()
	http.ListenAndServe(":3000", nil)
}
