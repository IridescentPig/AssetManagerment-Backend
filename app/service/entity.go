package service

import (
	"asset-management/app/dao"
	"asset-management/app/define"
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

func (entity *entityService) CheckIsInEntity(ctx *utils.Context, entityID uint) bool {
	userInfo, exists := ctx.Get("user")
	if exists {
		if userInfo, ok := userInfo.(define.UserBasicInfo); ok {
			if userInfo.EntityID == entityID {
				return true
			}
		}
	}
	return false
}

func (entity *entityService) GetParamID(ctx *utils.Context, key string) (uint, error) {
	param := ctx.Param(key)
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
	thisEntity, err := dao.EntityDao.GetEntityByID(entityID)
	if err != nil {
		return err
	}
	err = dao.UserDao.Create(model.User{
		UserName:    name,
		Password:    password,
		EntityID:    entityID,
		Entity:      thisEntity,
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

func (entity *entityService) GetEntityManagerList(id uint) ([]*model.User, error) {
	managerList, err := dao.EntityDao.GetEntityManager(id)
	return managerList, err
}

func (entity *entityService) DeleteManager(userID uint) error {
	err := dao.UserDao.Update(userID, map[string]interface{}{
		"entity_super": false,
	})
	return err
}

func (entity *entityService) ModifyEntity(entityID uint, modifyInfo define.ModifyEntityInfoReq) error {
	if modifyInfo.EntityName != nil && modifyInfo.Description != nil {
		err := dao.EntityDao.Update(entityID, map[string]interface{}{
			"name":        *modifyInfo.EntityName,
			"description": *modifyInfo.Description,
		})
		return err
	} else if modifyInfo.EntityName != nil {
		err := dao.EntityDao.Update(entityID, map[string]interface{}{
			"name": *modifyInfo.EntityName,
		})
		return err
	} else if modifyInfo.Description != nil {
		err := dao.EntityDao.Update(entityID, map[string]interface{}{
			"description": *modifyInfo.Description,
		})
		return err
	}
	return nil
}
