package models

import (
	"time"
)

type Message struct {
	Id         string    `json:"id" bson:"id,omitempty" required:"true"`
	Value      string    `required:"true" max:"160" json:"value" bson:"value"`
	ChatRoomId string    `required:"true" json:"chatRoomId" bson:"chatRoomId"`
	Username   string    `required:"true" json:"username"`
	Created    time.Time `required:"true" json:"created" bson:"created"`
}

type StockMessage struct {
	HubName             string `json:"hubName"`
	ClientRemoteAddress string `json:"clientRemoteAddress"`
	Message             string `json:"message"`
	RoomId              string `json:"roomId"`
}
