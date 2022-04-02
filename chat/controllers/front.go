package controllers

import (
	"encoding/json"
	"io"
	"net/http"
)

func RegisterController() {
	uc := newUserController()
	rm := newRoomController()
	ac := newAuthController()

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
