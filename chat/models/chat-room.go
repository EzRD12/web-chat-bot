package models

import (
	"time"
)

type ChatRoom struct {
	Id      string    `json:"id" bson:"_id,omitempty" required:"true"`
	Name    string    `validate:"nonzero, max=32, min=5" json:"name" bson:"name"`
	Created time.Time `json:"created" bson:"created"`
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
