package service

import (
	"asset-management/app/dao"
	"asset-management/app/model"
	"asset-management/myerror"
	"asset-management/utils"
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
	err := dao.EntityDao.Create(model.Entity{
		Name: name,
	})
	return err
}

func (entity *entityService) DeleteEntity(id uint) error {
	err := dao.EntityDao.Delete([]uint{id})
	return err
}

func (entity *entityService) GetAllEntity() ([]model.Entity, error) {
	entityList, err := dao.EntityDao.AllEntity()
	if err != nil {
		return []model.Entity{}, err
	}
	return entityList, nil
}

func (entity *entityService) ExistsEntityByName(name string) (bool, error) {
	thisEntity, err := dao.EntityDao.GetEntityByName(name)
	if err != nil {
		return false, err
	} else if thisEntity == nil {
		return false, nil
	} else {
		return true, nil
	}
}

func (entity *entityService) ExistsEntityByID(id uint) (bool, error) {
	thisEntity, err := dao.EntityDao.GetEntityByID(id)
	if err != nil {
		return false, err
	} else if thisEntity == nil {
		return false, nil
	} else {
		return true, nil
	}
}

func (entity *entityService) GetEntityInfoByID(id uint) (*model.Entity, error) {
	thisEntity, err := dao.EntityDao.GetEntityByID(id)
	if err != nil {
		return nil, err
	} else if thisEntity == nil {
		return nil, nil
	} else {
		return thisEntity, nil
	}
}

func (entity *entityService) GetUsersUnderEntity(id uint) ([]*model.User, error) {
	userList, err := dao.EntityDao.GetEntityAllUser(id)
	return userList, err
}

func (entity *entityService) GetAllDepartmentsUnderEntity(id uint) ([]*model.Department, error) {
	departmentList, err := dao.EntityDao.GetEntityAllDepartment(id)
	return departmentList, err
}

func (entity *entityService) CreateManager(name string, password string, entityID uint) error {
	password = utils.CreateMD5(password)
	err := dao.UserDao.Create(model.User{
		UserName:    name,
		Password:    password,
		EntityID:    entityID,
		EntitySuper: true,
	})
	return err
}

func (entity *entityService) SetManager(name string) error {
	err := dao.UserDao.UpdateByName(name, map[string]interface{}{
		"entity_super": true,
	})
	return err
}
