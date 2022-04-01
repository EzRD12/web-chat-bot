package models

type ChatRoom struct {
	Id      string `json:"id" bson:"id,omitempty"`
	Name    string `max:"65" required:"true"`
	Created string
}

type ChatRoomUsers struct {
	RoomId string `json:"roomId" bson:"roomId"`
	UserId string
	Joined string
}
