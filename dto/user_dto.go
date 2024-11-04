package dto

import (
    "time"
)

type UserCreate struct {
    Username    string
	Password	string
    Email       string
}

type UserUpdate struct {
    Username    string
	Password	string
    Description *string
    Email       string
}

type UserInfo struct {
    CreatedAt   time.Time
    Username    string
    Description *string
    Email       string
    Status      string
    Role        string
    Tag         []string
}