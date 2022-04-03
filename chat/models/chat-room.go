package models

import (
	"time"
)

type ChatRoom struct {
	Id      string `json:"id" bson:"_id,omitempty"`
	Name    string `max:"65" required:"true"`
	Created time.Time
}

type ChatRoomUsers struct {
	RoomId string `json:"roomId" bson:"roomId"`
	UserId string
	Joined time.Time
}

type ChatRoomDetails struct {
	Id         string
	Name       string
	Created    time.Time
	UsersCount int64
}
