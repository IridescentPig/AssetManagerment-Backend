package service

import (
	"asset-management/app/dao"
	"asset-management/app/define"
	"asset-management/app/model"
	"asset-management/myerror"
	"asset-management/utils"
	"errors"
)

type departmentService struct{}

var DepartmentService *departmentService

func newDepartmentService() *departmentService {
	return &departmentService{}
}

func init() {
	DepartmentService = newDepartmentService()
}

func (department *departmentService) GetHeaderDepartmentID(ctx *utils.Context) uint {
	userInfo, exists := ctx.Get("user")
	if exists {
		if userInfo, ok := userInfo.(define.UserBasicInfo); ok {
			return userInfo.DepartmentID
		}
	}
	return 0
}

func (department *departmentService) CheckIsInDepartment(ctx *utils.Context, departmentID uint) bool {
	userInfo, exists := ctx.Get("user")
	if exists {
		if userInfo, ok := userInfo.(define.UserBasicInfo); ok {
			if userInfo.DepartmentID == departmentID {
				return true
			}
		}
	}
	return false
}

func (department *departmentService) CheckIsAncestor(userDepartmentID uint, operateDepartmentID uint) (bool, error) {
	operateDepartment, err := department.GetDepartmentInfoByID(operateDepartmentID)
	if err != nil {
		return false, err
	}
	if operateDepartment == nil {
		return false, errors.New(myerror.DEPARTMENT_NOT_FOUND_INFO)
	}

	flag := false
	for {
		if operateDepartment == nil {
			break
		}
		if operateDepartment.ID == userDepartmentID {
			flag = true
			break
		}
		operateDepartment = operateDepartment.Parent
	}
	return flag, nil
}

/*
拥有部门权限，也即可以看到本部门及下属部门的信息
*/
func (department *departmentService) CheckDepartmentIdentity(ctx *utils.Context, entityID uint, departmentID uint) (bool, error) {
	// 不在实体中必定没有权限
	if !EntityService.CheckIsInEntity(ctx, entityID) {
		return false, nil
	}
	// 实体的系统管理员必然有权限
	if UserService.EntitySuper(ctx) {
		return true, nil
	}
	// 是资产管理员，有对下属部门的权限
	if UserService.DepartmentSuper(ctx) {
		return department.CheckIsAncestor(department.GetHeaderDepartmentID(ctx), departmentID)
	} else {
		// 普通员工只有对自己所在部门的权限
		return department.CheckIsInDepartment(ctx, departmentID), nil
	}
}

func (department *departmentService) CreateDepartment(entityID uint, departmentID uint) error {
	var err error
	if departmentID != 0 {
		err = dao.DepartmentDao.Create(model.Department{
			EntityID: entityID,
			ParentID: departmentID,
		})
	} else {
		err = dao.DepartmentDao.Create(model.Department{
			EntityID: entityID,
		})
	}
	return err
}

func (department *departmentService) GetDepartmentInfoByID(departmentID uint) (*model.Department, error) {
	thisDepartment, err := dao.DepartmentDao.GetDepartmentByID(departmentID)
	if err != nil {
		return nil, err
	}
	return thisDepartment, err
}
