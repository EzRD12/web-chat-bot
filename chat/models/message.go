package models

import (
	"time"
)

type Message struct {
	Id           string    `json:"id" bson:"id,omitempty" required:"true"`
	Value        string    `required:"true" max:"160"`
	ChatRoomId   string    `required:"true"`
	SenderUserId string    `required:"true"`
	Created      time.Time `required:"true"`
}

type StockRequest struct {
	Code   string `json:"code"`
	RoomId string `json:"roomId"`
}

type StockResponse struct {
	Code   string  `json:"code"`
	RoomId string  `json:"roomId"`
	Open   float64 `json:"open"`
	Date   string  `json:"open"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Close  float64 `json:"close"`
	Volume int64   `json:"volume"`
}
