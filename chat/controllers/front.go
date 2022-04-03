package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ezrod12/chat/messager"
	"github.com/ezrod12/chat/models"
	"github.com/ezrod12/chat/services"
	"github.com/ezrod12/chat/startup"
)

func RegisterController() {
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

			var message string = fmt.Sprintf("code: %s - date: %s \n open: %f, high: %f, low: %f, close: %f \n volume: %d", response.Code, response.Date,
				response.Open, response.High, response.Low, response.Close, response.Volume)
			var messageRequest models.Message = models.Message{Value: message, ChatRoomId: response.RoomId, SenderUserId: "system"}

			services.AddMessage(messageRequest, rm.collection, rm.context)
		}
	}()

	http.Handle("/users", *uc)
	http.Handle("/users/", *uc)
	http.Handle("/rooms", *rm)
	http.Handle("/auth", *ac)
	http.Handle("/auth/", *ac)
	http.Handle("/rooms/", *rm)
}

func encodeResponseAsJson(data interface{}, w io.Writer) {
	enc := json.NewEncoder(w)
	enc.Encode(data)
}
