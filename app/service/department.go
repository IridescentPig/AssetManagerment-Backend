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

func (department *departmentService) ExistsDepartmentByID(departmentID uint) (bool, error) {
	thisDepartment, err := dao.DepartmentDao.GetDepartmentByID(departmentID)
	if err != nil {
		return false, err
	} else if thisDepartment == nil {
		return false, nil
	}
	return true, nil
}

func (department *departmentService) ExistsDepartmentSub(departmentName string, entityID uint, departmentID uint) (bool, error) {
	thisDepartment, err := dao.DepartmentDao.GetDepartmentSub(departmentName, entityID, departmentID)
	if err != nil {
		return false, err
	} else if thisDepartment == nil {
		return false, nil
	}
	return true, nil
}

func (department *departmentService) CheckDepartmentInEntity(entityID uint, departmentID uint) (bool, error) {
	thisDepartment, err := dao.DepartmentDao.GetDepartmentByID(departmentID)
	if err != nil {
		return false, err
	}
	return thisDepartment.EntityID == entityID, nil
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

func (department *departmentService) CreateDepartment(name string, entityID uint, departmentID uint) error {
	var err error
	if departmentID != 0 {
		err = dao.DepartmentDao.Create(model.Department{
			Name:     name,
			EntityID: entityID,
			ParentID: departmentID,
		})
	} else {
		err = dao.DepartmentDao.Create(model.Department{
			Name:     name,
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

func (department *departmentService) GetSubDepartments(departmentID uint) ([]*model.Department, error) {
	departmentList, err := dao.DepartmentDao.GetSubDepartmentByID(departmentID)
	return departmentList, err
}

func (department *departmentService) GetAllUsers(departmentID uint) ([]*model.User, error) {
	userList, err := dao.DepartmentDao.GetDepartmentAllUserByID(departmentID)
	return userList, err
}

func (department *departmentService) CreateDepartmentUser(req define.CreateDepartmentUserReq, entityID uint, departmentID uint) error {
	req.Password = utils.CreateMD5(req.Password)
	err := dao.UserDao.Create(model.User{
		UserName:        req.UserName,
		Password:        req.Password,
		DepartmentSuper: req.DepartmentSuper,
		EntityID:        entityID,
		DepartmentID:    departmentID,
	})
	return err
}

func (department *departmentService) SetDepartmentManager(username string, departmentID uint) error {
	err := dao.UserDao.UpdateByName(username, map[string]interface{}{
		"department_id":    departmentID,
		"department_super": true,
	})
	return err
}

func (department *departmentService) DeleteDepartmentManager(userID uint) error {
	err := dao.UserDao.Update(userID, map[string]interface{}{
		"department_super": false,
	})
	return err
}

func (department *departmentService) GetDepartmentManagerList(id uint) ([]*model.User, error) {
	managerList, err := dao.DepartmentDao.GetDepartmentManager(id)
	return managerList, err
}

func (department *departmentService) DeleteDepartment(id uint) error {
	return dao.DepartmentDao.Delete([]uint{id})
}

func (department *departmentService) DepartmentHasUsers(departmentID uint) (bool, error) {
	hasUsers, err := dao.DepartmentDao.GetDepartmentAllUserByID(departmentID)
	if err != nil {
		return true, err
	}
	return len(hasUsers) != 0, nil
}
