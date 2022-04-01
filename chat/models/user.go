package models

type User struct {
	Id             string `json:"id" bson:"id,omitempty" required:"true"`
	FirstName      string
	LastName       string
	Username       string
	Password       string
	LastConnection string
}

var (
	users  []*User
	nextId = 1
)
