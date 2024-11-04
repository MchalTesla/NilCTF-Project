package repositories_interface

import (
	"NilCTF/models"
)

type ConfigRepositoryInterface interface {
	Upsert(key, value string) error
	Get(key string) (string, error)
	Delete(key string) error
	List(condition, value string) ([]models.Config, error)
}