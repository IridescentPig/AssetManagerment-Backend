package api

import (
	"asset-management/app/service"
	"asset-management/myerror"
	"asset-management/utils"
)

type departmentApi struct {
}

var DepartmentApi *departmentApi

func newdepartmentApi() *departmentApi {
	return &departmentApi{}
}

func init() {
	DepartmentApi = newdepartmentApi()
}

/*
目前按照 /entity/:entity_id/department/create 和 /entity/:entity_id/department/:id/department/create 来处理
*/
func (department *departmentApi) CreateDepartment(ctx utils.Context) {
	isEntitySuper := service.UserService.EntitySuper(&ctx)
	if !isEntitySuper {
		ctx.Forbidden(myerror.PERMISSION_DENIED, myerror.PERMISSION_DENIED_INFO)
		return
	}

}

func (department *departmentApi) GetDepartmentByID(ctx utils.Context) {

}
