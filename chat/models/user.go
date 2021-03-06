package models

import (
	"time"
)

type User struct {
	Id              string `json:"id" bson:"_id,omitempty" required:"true"`
	FirstName       string `required:"true" validate:"max=32"`
	LastName        string `required:"true" validate:"max=32"`
	Username        string `required:"true" validate:"max=20, nonzero"`
	Password        string `validate:"min=8"`
	AdminPermission bool   `bson:"admin_permission"`
	LastConnection  time.Time
}

type AuthenthicationRequest struct {
	Username string
	Password string
}
