package service

import (
	"asset-management/app/model"
	"asset-management/myerror"
	"asset-management/utils"
	"errors"
	"strconv"
)

type entityService struct{}

var EntityService *entityService

func newEntityService() *entityService {
	return &entityService{}
}

func init() {
	EntityService = newEntityService()
}

func (entity *entityService) GetParamID(ctx *utils.Context) (uint, error) {
	param := ctx.Param("id")
	tempID, err := strconv.ParseUint(param, 10, 0)
	if err != nil {
		ctx.BadRequest(myerror.INVALID_PARAM, myerror.INVALID_PARAM_INFO)
		return 0, err
	}
	entityID := uint(tempID)
	return entityID, nil
}

func (entity *entityService) CreateEntity(name string) error {
	return errors.New("")
}

func (entity *entityService) GetAllEntity() ([]*model.Entity, error) {
	entityTest := model.Entity{
		ID:   1,
		Name: "Test",
	}
	entityTest2 := model.Entity{
		ID:   2,
		Name: "Test2",
	}
	return []*model.Entity{&entityTest, &entityTest2}, nil
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

func (entity *entityService) GetUsersUnderEntity(id uint) ([]*model.User, error) {
	return []*model.User{}, nil
}
