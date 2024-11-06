package dto

import (
    "time"
)

type UserCreate struct {
    Username    string
	Password	string
    Email       string
}

type UserInfo struct {
    Username    string
    Email       string
    Description *string
    Status      string
    Role        string
    Tag         []string
    CreatedAt   time.Time
}

type UserUpdate struct {
    Username    string
	Password	string
    Email       string
    Description *string
}

type UserInfoByAdmin struct {
    ID          uint
    Username    string
    Email       string
    Description *string
    Status      string
    Role        string
    Tag         []string
    CreatedAt   time.Time
}

type UserUpdateByAdmin struct {
    ID          uint
    Username    string
    Password	string
    Email       string
    Description *string
    Status      string
    Role        string
    Tag         []string
}

