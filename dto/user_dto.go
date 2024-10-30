package dto

import (
    "time"
)

type UserUpdate struct {
    Username    *string
	Password	*string
    Description *string
    Email       *string
}

type UserInfo struct {
    CreatedAt   *time.Time
    Username    *string
    Description *string
    Email       *string
    Status      *string
    Role        *string
    Tag         *string
}