package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/ezrod12/chat/messager"
	"github.com/ezrod12/chat/models"
	"github.com/ezrod12/chat/pkg"
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
			var response models.StockMessage
			json.Unmarshal(d.Body, &response)

			hubs := *sv.GetHubs()
			hub := hubs[response.HubName]

			message := &models.Message{Value: response.Message, Username: "system", Created: time.Now(), ChatRoomId: response.RoomId}
			result, _ := json.Marshal(message)

			hub.SendTo(string(result), response.ClientRemoteAddress)
		}
	}()

	http.Handle("/", loginForm)
	http.Handle("/home", homeForm)
	http.Handle("/signup", signUpForm)
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

var signUpForm = http.HandlerFunc(
	func(w http.ResponseWriter, r *http.Request) {

		http.ServeFile(w, r, "./templates/signup.html")
	})

var homeForm = http.HandlerFunc(
	func(w http.ResponseWriter, r *http.Request) {

		http.ServeFile(w, r, "./templates/home.html")
	})

var chatForm = http.HandlerFunc(
	func(w http.ResponseWriter, r *http.Request) {

		http.ServeFile(w, r, "./templates/chat.html")
	})
