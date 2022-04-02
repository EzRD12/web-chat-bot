package controllers

import (
	"encoding/json"
	"io"
	"net/http"
)

func RegisterController() {
	uc := newUserController()
	rm := newRoomController()

	http.Handle("/users", *uc)
	http.Handle("/users/", *uc)
	http.Handle("/rooms", *rm)
	http.Handle("/rooms/", *rm)
}

func encodeResponseAsJson(data interface{}, w io.Writer) {
	enc := json.NewEncoder(w)
	enc.Encode(data)
}
