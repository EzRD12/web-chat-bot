package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ezrod12/chat/messager"
	"github.com/ezrod12/chat/models"
	"github.com/ezrod12/chat/pkg"
	"github.com/ezrod12/chat/services"
	"github.com/ezrod12/chat/startup"
)

func RegisterController(sv *pkg.Server) {
	uc := newUserController()
	rm := newRoomController()
	ac := newAuthController()
	st := startup.InitStartup()

	st.SaveSeedData()

	msgs := messager.ReceiveMessageDeliveryChannel()

	go func() {
		for d := range msgs {
			var response models.StockResponse
			json.Unmarshal(d.Body, &response)

			var message string = fmt.Sprintf("%s quote is: $%f \n per share", response.Code, response.Close)
			var messageRequest models.Message = models.Message{Value: message, ChatRoomId: response.RoomId, SenderUserId: "bot"}

			services.AddMessage(messageRequest, rm.collection, rm.context)
		}
	}()

	http.Handle("/", loginForm)
	http.Handle("/home", homeForm)
	http.Handle("/chat", chatForm)
	http.Handle("/users", *uc)
	http.Handle("/users/", *uc)
	http.HandleFunc("/ws", sv.ServeWs)
	http.Handle("/rooms", *rm)
	http.Handle("/auth", *ac)
	http.Handle("/auth/", *ac)
	http.Handle("/rooms/", *rm)
}

func encodeResponseAsJson(data interface{}, w io.Writer) {
	enc := json.NewEncoder(w)
	enc.Encode(data)
}

var loginForm = http.HandlerFunc(
	func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		http.ServeFile(w, r, "./templates/login.html")
	})

var homeForm = http.HandlerFunc(
	func(w http.ResponseWriter, r *http.Request) {

		http.ServeFile(w, r, "./templates/home.html")
	})

var chatForm = http.HandlerFunc(
	func(w http.ResponseWriter, r *http.Request) {

		http.ServeFile(w, r, "./templates/chat.html")
	})
