package service

import (
	"asset-management/app/model"
	"errors"
)

type entityService struct{}

var EntityService *entityService

func newEntityService() *entityService {
	return &entityService{}
}

func init() {
	EntityService = newEntityService()
}

func (entity *entityService) CreateEntity(name string) error {
	return errors.New("")
}

func (entity *entityService) GetAllEntity() ([]*model.Entity, error) {
	return []*model.Entity{}, nil
}

func (entity *entityService) ExistsEntityByID(id uint) (bool, error) {
	return true, errors.New("")
}

func (entity *entityService) ExistsEntityByName(name string) (bool, error) {
	return true, errors.New("")
}

func (entity *entityService) GetEntityInfoByID(id uint) (*model.Entity, error) {
	return &model.Entity{}, errors.New("")
}
