package services_interface

import (
	"NilCTF/dto"
)

type HomeServiceInterface interface {
	Info(userID uint) (*dto.UserInfo, error)
}