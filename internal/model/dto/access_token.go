package dto

import "time"

type AccessToken struct {
	UserUuid string        `json:"user_uuid"`
	Token    string        `json:"token"`
	Lt       time.Duration `json:"lt"`
}
